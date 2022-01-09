# Check chart without installing
helm install geared-marsupi ./helmchart --dry-run --debug

helm upgrade --install ./helmchart \
  --set name="alpha" \
  --set http_port="8080"



helm install alpha ./helmchart \
  --set name="alpha" \
  --set http_port="8080"

helm status alpha

helm uninstall alpha

helm install beta ./helmchart \
  --set name="beta" \
  --set http_port=8081 \
  --set environment.GET_URI="http://ip:30080" \
  --set environment.POST_URI=http://ip:30080 \
  --set environment.POST_SLEEP=3  \
  --set node_port=30081

helm uninstall beta

kubectl logs -f beta-deployment-5595f7c8bd-wrpvv