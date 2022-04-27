helm repo add fluent https://fluent.github.io/helm-charts

helm install -f ./logging/values.yaml fluent-bit fluent/fluent-bit --namespace observability