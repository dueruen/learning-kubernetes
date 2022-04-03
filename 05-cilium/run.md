# Install
xport HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)

curl -L --remote-name-all https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-arm64.tar.gz{,.sha256sum}

sha256sum --check hubble-linux-arm64.tar.gz.sha256sum

sudo tar xzvfC hubble-linux-arm64.tar.gz /usr/local/bin

rm hubble-linux-arm64.tar.gz{,.sha256sum}

kubectl -n kube-system port-forward svc/hubble-ui 8080:80 --address='0.0.0.0'
kubectl -n cilium-monitoring port-forward svc/grafana 3000:3000 --address='0.0.0.0'
kubectl -n cilium-monitoring port-forward svc/prometheus 3001:9090 --address='0.0.0.0'
kubectl port-forward svc/nats-demo-subscriber-service 8081:8080 --address='0.0.0.0'

kubectl -n jaeger port-forward svc/jaeger-default-query 8082:16686 --address='0.0.0.0'

kubectl port-forward svc/my-kube-state-metrics 8083:8080 --address='0.0.0.0'

kubectl create -f l7-policy.yaml
kubectl apply -f producer-low-access.yaml
kubectl apply -f producer.yaml
kubectl apply -f subscriber.yaml

helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

helm install cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --set nodeinit.enabled=true \
   --set kubeProxyReplacement=partial \
   --set hostServices.enabled=false \
   --set externalIPs.enabled=true \
   --set nodePort.enabled=true \
   --set hostPort.enabled=true \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes \
   --set hubble.enabled=true \
   --set hubble.listenAddress=":4244" \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true

helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,icmp,http}" \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.relay.enabled=true \
   --set externalIPs.enabled=true \
   --set hubble.ui.enabled=true


helm install cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"   