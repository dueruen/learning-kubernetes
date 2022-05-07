./start.sh

kubectl port-forward svc/observ-grafana 3000:80 --address='0.0.0.0' -n observability
kubectl port-forward svc/observ-kube-prometheus-sta-prometheus 3001:9090 --address='0.0.0.0' -n observability