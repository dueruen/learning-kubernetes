helm install observ . --create-namespace --namespace observability

kubectl port-forward svc/observ-grafana 3000:80 --address='0.0.0.0' -n observability