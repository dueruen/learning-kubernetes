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

kubectl -n linkerd-viz patch deployment web -p '{"spec":{"spec":{"containers":[{"args":["-enforced-host=^.*"]}]}}}'