kubectl create namespace confluent

helm repo add confluentinc https://packages.confluent.io/helm
helm repo update

helm upgrade --install confluent-operator confluentinc/confluent-for-kubernetes --namespace confluent

helm install confluent-kafka . --create-namespace --namespace confluent

kubectl delete pvc --all -n confluent