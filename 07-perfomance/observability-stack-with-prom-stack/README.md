helm dependency build
helm install observ . --create-namespace --namespace observability

kubectl get secret --namespace observability observ-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
kubectl port-forward svc/observ-grafana 3000:80 --address='0.0.0.0' -n observability