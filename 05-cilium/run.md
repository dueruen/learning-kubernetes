# Install
xport HUBBLE_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)

curl -L --remote-name-all https://github.com/cilium/hubble/releases/download/$HUBBLE_VERSION/hubble-linux-arm64.tar.gz{,.sha256sum}

sha256sum --check hubble-linux-arm64.tar.gz.sha256sum

sudo tar xzvfC hubble-linux-arm64.tar.gz /usr/local/bin

rm hubble-linux-arm64.tar.gz{,.sha256sum}



helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.relay.enabled=true \
   --set hubble.ui.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"