kubectl apply -f alpha.yml
kubectl get pods -o wide
kubectl get services -o wide

kubectl logs simple-alpha

kubectl delete -f alpha.yml