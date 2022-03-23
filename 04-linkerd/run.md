Forward service to remote host:
kubectl -n emojivoto port-forward svc/web-svc 8080:80 --address='0.0.0.0'

kubectl get namespace
kubectl get pods -n linkerd-viz
kubectl get service -n linkerd-viz

kubectl get deployment web -n linkerd-viz -o yaml

install step binary from github release fixes Qt error

helm install linkerd2 \
  --set-file identityTrustAnchorsPEM=ca.crt \
  --set-file identity.issuer.tls.crtPEM=issuer.crt \
  --set-file identity.issuer.tls.keyPEM=issuer.key \
  linkerd/linkerd2

helm install linkerd-viz \
  --set dashboard.enforcedHostRegexp="^.*" \
  linkerd/linkerd-viz

helm install linkerd-viz \
  --set dashboard.enforcedHostRegexp="^.*" \
  linkerd-edge/linkerd-viz  

helm install linkerd-jaeger \
  --set collector.image.version="0.43.0-arm64" \
  linkerd/linkerd-jaeger

helm install linkerd-jaeger \
  linkerd-edge/linkerd-jaeger  

Clone repo and edit manualy, not working
helm dependency update
helm install linkerd-viz . --create-namespace -n linkerd-viz

kubectl -n linkerd-viz port-forward svc/web 8084:8084 --address='0.0.0.0'

helm install linkerd-jaeger . --create-namespace -n linkerd-jaeger \
  --set webhook.image.version="2.11.1"

kubectl -n linkerd-viz patch deployment web -p '{"metadata":{"annotations":{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"linkerd.io/created-by":"linkerd/helm stable-2.11.1"},"labels":{"app.kubernetes.io/name":"web","app.kubernetes.io/part-of":"Linkerd","app.kubernetes.io/version":"stable-2.11.1","component":"web","linkerd.io/extension":"viz","namespace":"linkerd-viz"},"name":"web","namespace":"linkerd-viz"},"spec":{"replicas":1,"selector":{"matchLabels":{"component":"web","linkerd.io/extension":"viz","namespace":"linkerd-viz"}},"template":{"metadata":{"annotations":{"linkerd.io/created-by":"linkerd/helm stable-2.11.1"},"labels":{"component":"web","linkerd.io/extension":"viz","namespace":"linkerd-viz"}},"spec":{"containers":[{"args":["-linkerd-metrics-api-addr=metrics-api.linkerd-viz.svc.cluster.local:8085","-cluster-domain=cluster.local","-grafana-addr=grafana.linkerd-viz.svc.cluster.local:3000","-controller-namespace=linkerd","-viz-namespace=linkerd-viz","-log-level=info","-log-format=plain","-enforced-host=^.*"],"image":"cr.l5d.io/linkerd/web:stable-2.11.1","imagePullPolicy":"IfNotPresent","livenessProbe":{"httpGet":{"path":"/ping","port":9994},"initialDelaySeconds":10},"name":"web","ports":[{"containerPort":8084,"name":"http"},{"containerPort":9994,"name":"admin-http"}],"readinessProbe":{"failureThreshold":7,"httpGet":{"path":"/ready","port":9994}},"resources":null,"securityContext":{"runAsUser":2103}}],"nodeSelector":{"beta.kubernetes.io/os":"linux"},"serviceAccountName":"web"}}}}}'

kubectl -n linkerd-viz patch deployment web -p '{"metadata":{"annotations":{"kubectl.kubernetes.io/last-applied-configuration":{"spec":{"spec":{"containers":[{"args":["-linkerd-metrics-api-addr=metrics-api.linkerd-viz.svc.cluster.local:8085","-cluster-domain=cluster.local","-grafana-addr=grafana.linkerd-viz.svc.cluster.local:3000","-controller-namespace=linkerd","-viz-namespace=linkerd-viz","-log-level=info","-log-format=plain","-enforced-host=^.*"]}]}}}}}}'

kubectl -n linkerd-viz patch deployment web -p '{"spec":{"spec":{"containers":[{"args":["-linkerd-metrics-api-addr=metrics-api.linkerd-viz.svc.cluster.local:8085","-cluster-domain=cluster.local","-grafana-addr=grafana.linkerd-viz.svc.cluster.local:3000","-controller-namespace=linkerd","-viz-namespace=linkerd-viz","-log-level=info","-log-format=plain","-enforced-host=^.*"]}]}}}'

kubectl -n linkerd-viz patch deployment web -p '{"spec":{"spec":{"containers":[{"args":["-enforced-host=^.*"]}]}}}'