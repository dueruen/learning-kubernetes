https://opensource.com/article/20/6/kubernetes-raspberry-pi
https://www.learnlinux.tv/building-a-10-node-raspberry-pi-kubernetes-cluster/
https://www.youtube.com/watch?v=MO8N79lQSWU&ab_channel=LearnLinuxTV
cilium as network driver https://www.talos.dev/v0.14/guides/deploying-cilium/

ssh ubuntu@ip
psw: ubuntu

# Change password
newPass

# Create user
sudo adduser newUser
psw: newPass

sudo usermod -aG sudo newUser
ctl + d

ssh newUser@ip

# Change hostname
echo "k8s-master-01" | sudo tee /etc/hostname
echo "k8s-worker-02" | sudo tee /etc/hostname

sudo nano /etc/hosts
127.0.1.1 k8s-worker-02

sudo apt update && sudo apt dist-upgrade

sudo sed -i '$ s/$/ cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1 swapaccount=1/' /boot/firmware/cmdline.txt

# Install docker
curl -sSL get.docker.com | sudo sh
sudo usermod -aG docker newUser

# Enable routing
sudo vim /etc/sysctl.conf
sudo reboot

# Check things are working
systemctl status docker
docker run hello-world

# Add kubernetes
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
sudo apt update
sudo apt install -y kubelet kubeadm kubectl

# Run only on CONTROLLER/MASTER
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
## Result from cmd above 
## This is important DO NOT UPLOAD OR SHARE
sudo kubeadm join <ip> --token <token> --discovery-token-ca-cert-hash <hash>

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
 
use cilium OR flannel
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

kubectl get pods --all-namespaces

kubectl get nodes

# Run containers
kubectl get pods
kubectl get services
kubectl apply -f <filename>

kubectl get pods -o wide

kubectl delete pod nginx-example
kubectl delete service nginx-example

# Configure ssh
https://www.digitalocean.com/community/tutorials/how-to-configure-ssh-key-based-authentication-on-a-linux-server

sudo vi /etc/hosts
## Add ip and names
k8s-master-01 ip
k8s-worker-01 ip

## copy public key
cat ~/.ssh/id_rsa.pub | ssh newUser@k8s-worker-01 "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"

## Add to windows terminal
{
    "commandline": "ssh newUser@k8s-worker-01",
    "guid": "{07b52e3e-de2c-5db4-bd2d-ba144ed6c800}",
    "hidden": false,
    "name": "k8s-worker-01"
}


# Reset cluser
sudo kubeadm reset