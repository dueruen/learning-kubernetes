#!/bin/bash

helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

kubectl apply -k github.com/cilium/kustomize-bases/cert-manager   

(
  set -e
  kubectl wait deployment --namespace="cert-manager" --for="condition=Available" cert-manager-webhook cert-manager-cainjector cert-manager --timeout=3m
  kubectl wait pods --namespace="cert-manager" --for="condition=Ready" --all --timeout=3m
  kubectl wait apiservice --for="condition=Available" v1.cert-manager.io v1.acme.cert-manager.io --timeout=3m
  until kubectl get secret --namespace="cert-manager" cert-manager-webhook-ca 2> /dev/null ; do sleep 0.5 ; done
)

kubectl apply -k github.com/cilium/kustomize-bases/jaeger

cat > jaeger.yaml << EOF
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger-default
  namespace: jaeger
spec:
  strategy: allInOne
  storage:
    type: memory
    options:
      memory:
        max-traces: 100000
  ingress:
    enabled: false
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
EOF
kubectl apply -f jaeger.yaml

kubectl apply -k github.com/cilium/kustomize-bases/opentelemetry
cat > otelcol.yaml << EOF
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: otelcol-hubble
  namespace: kube-system
spec:
  mode: daemonset
  image: ghcr.io/cilium/hubble-otel/otelcol:v0.1.1
  env:
    # set NODE_IP environment variable using downwards API
    - name: NODE_IP
      valueFrom:
        fieldRef:
          fieldPath: status.hostIP
  volumes:
    # this example connect to Hubble socket of Cilium agent
    # using host port and TLS
    - name: hubble-tls
      projected:
        defaultMode: 256
        sources:
          - secret:
              name: hubble-relay-client-certs
              items:
                - key: tls.crt
                  path: client.crt
                - key: tls.key
                  path: client.key
                - key: ca.crt
                  path: ca.crt
    # it's possible to use the UNIX socket also, for which
    # the following volume will be needed
    # - name: cilium-run
    #   hostPath:
    #     path: /var/run/cilium
    #     type: Directory
  volumeMounts:
    # - name: cilium-run
    #   mountPath: /var/run/cilium
    - name: hubble-tls
      mountPath: /var/run/hubble-tls
      readOnly: true
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:55690
      hubble:
        # NODE_IP is substituted by the collector at runtime
        # the '\' prefix is required only in order for this config to be
        # inlined in the guide and make it easy to paste, i.e. to avoid
        # shell subtituting it
        endpoint: \${NODE_IP}:4244 # unix:///var/run/cilium/hubble.sock
        buffer_size: 100
        include_flow_types:
          # this sets an L7 flow filter, removing this section will
          # disable filtering and result all types of flows being turned
          # into spans;
          # other type filters can be set, the names are same as what's
          # used in 'hubble observe -t <type>'
          traces: ["l7"]
        tls:
          insecure_skip_verify: true
          ca_file: /var/run/hubble-tls/ca.crt
          cert_file: /var/run/hubble-tls/client.crt
          key_file: /var/run/hubble-tls/client.key
    processors:
      batch:
        timeout: 30s
        send_batch_size: 100

    exporters:
      jaeger:
        endpoint: jaeger-default-collector.jaeger.svc.cluster.local:14250
        tls:
          insecure: true

    service:
      telemetry:
        logs:
          level: info
      pipelines:
        traces:
          receivers: [hubble, otlp]
          processors: [batch]
          exporters: [jaeger]
EOF
kubectl apply -f otelcol.yaml