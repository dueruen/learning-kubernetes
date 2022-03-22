Forward service to remote host:
kubectl -n emojivoto port-forward svc/web-svc 8080:80 --address='0.0.0.0'

kubectl get namespace
kubectl get pods -n linkerd-viz
kubectl get service -n linkerd-viz