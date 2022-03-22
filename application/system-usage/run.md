./usage.sh load-one-pod-cluster-master-01.txt
./usage.sh load-one-pod-cluster-worker-01.txt
./usage.sh load-one-pod-cluster-worker-02.txt

scp k8s-master-01:/home/dueruen/code/learning-kubernetes/application/system-usage/one-pod-cluster-master-01.txt .
scp k8s-worker-01:/home/dueruen/code/usage/one-pod-cluster-worker-01.txt .
scp k8s-worker-02:/home/dueruen/code/usage/one-pod-cluster-worker-02.txt .

ssh k8s-master-01 'bash -s' < /home/dueruen/code/learning-kubernetes/application/system-usage/usage.sh

scp k8s-master-01:/home/dueruen/code/learning-kubernetes/application/system-usage/load-one-pod-cluster-master-01.txt .
scp k8s-worker-01:/home/dueruen/code/usage/load-one-pod-cluster-worker-01.txt .
scp k8s-worker-02:/home/dueruen/code/usage/load-one-pod-cluster-worker-02.txt .

siege -c5 -t10s http://127.0.0.1:38333/api