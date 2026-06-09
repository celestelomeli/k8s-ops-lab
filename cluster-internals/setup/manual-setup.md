# Manual Cluster Setup

These are the steps run by hand on each EC2 instance before switching to Ansible. Both nodes require everything up through the join command section. Only the control plane runs kubeadm init.

## Disable swap

Kubernetes requires swap to be off. If swap is on, the kubelet refuses to start.

    swapoff -a
    sed -i '/ swap / s/^/#/' /etc/fstab

## Load kernel modules

These modules enable container networking. overlay is used by containerd for the container filesystem. br_netfilter lets iptables see traffic crossing the network bridge.

    modprobe overlay
    modprobe br_netfilter
    echo -e "overlay\nbr_netfilter" > /etc/modules-load.d/k8s.conf

## Set sysctl params

Tells the kernel to route bridge traffic through iptables, which is required for Kubernetes networking to work.

    cat <<EOF > /etc/sysctl.d/k8s.conf
    net.bridge.bridge-nf-call-iptables  = 1
    net.bridge.bridge-nf-call-ip6tables = 1
    net.ipv4.ip_forward                 = 1
    EOF
    sysctl --system

## Install containerd

containerd is the container runtime. It's what actually runs the containers.

    apt-get update && apt-get install -y containerd
    mkdir -p /etc/containerd
    containerd config default > /etc/containerd/config.toml

Then edit `/etc/containerd/config.toml` and set `SystemdCgroup = true` under the runc options. Without this, the kubelet and containerd use different cgroup drivers and the node breaks.

    systemctl restart containerd

## Install kubeadm, kubelet, kubectl

    apt-get install -y apt-transport-https ca-certificates curl
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' > /etc/apt/sources.list.d/kubernetes.list
    apt-get update
    apt-get install -y kubelet kubeadm kubectl
    apt-mark hold kubelet kubeadm kubectl

`apt-mark hold` prevents them from being upgraded automatically, which can break the cluster if versions drift between nodes.

## Initialize the control plane (control plane only)

    kubeadm init --pod-network-cidr=192.168.0.0/16

Copy the join command from the output. Then set up kubeconfig for the ubuntu user:

    mkdir -p $HOME/.kube
    cp /etc/kubernetes/admin.conf $HOME/.kube/config
    chown $(id -u):$(id -g) $HOME/.kube/config

## Install Calico CNI (control plane only)

Without a CNI plugin nodes stay in NotReady state. Calico handles pod networking.

    kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.27.0/manifests/calico.yaml

## Join the worker (worker only)

Run the join command copied from the kubeadm init output:

    kubeadm join <control-plane-ip>:6443 --token <token> --discovery-token-ca-cert-hash sha256:<hash>

## Verify

Run this on the control plane. Both nodes should show Ready after a minute.

    kubectl get nodes
