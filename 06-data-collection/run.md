helm install my-kube-state-metrics prometheus-community/kube-state-metrics --version 4.7.0

promtail add to default values:
extraArgs:
  - -client.url=http://scraper-service:80
helm install loki-promtail grafana/promtail --values /tmp/loki-promtail.yaml -n loki --create-namespace

http://scraper-service:80/

-client.url=http://loki-stack:3100/loki/api/v1/push
http://loki-stack:3100/loki/api/v1/push

helm uninstall loki-promtail -n loki

kubectl --namespace loki port-forward daemonset/loki-promtail 8084 --address='0.0.0.0'

kubectl describe pod/loki-promtail-9gs8w -n loki

kubectl port-forward svc/kafka-demo-service 8080:80 --address='0.0.0.0'
kubectl port-forward svc/my-kube-state-metrics 8080:8080 --address='0.0.0.0'