## Singular

Singular design for deploy and operation ContainerOps platform, mostly focus on Kubernetes, Prometheus, others in the Cloud Native technology stack. We are trying to build all stack cross cloud, OpenStack, bare metals.

Singular don't use any other deploy tools like _kubeadm_, deploy everything in a hard way instead. Singular providers templates for service so could deploy any version. Singular deploys development versions of **Kubernetes**, **CoreDNS**, others in **CNCF** CI demo.

### Deployment Template

Singular uses a **YAML** file as deploy template, it describes the architecture of the cluster. It don't only deploy ContainerOps modules, also use to deploy Kubernetes, others in Cloud Native stack and some common software.

#### SSH Key

When deploy infrastructures, **Singular** need to _SSH_ to virtual machine or bare metal.

1. At least provide SSH private key, **Singular** will create public key from it. Generate ssh key follow Github document - [How to generate SSH Key](https://help.github.com/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).
2. If no SSH key files in deployment template, **Singular** will create _SSH_ pair key files in default folder(**_$HOME/.containerops/ssh_**) and name(**_id_rsa.pub_** and **_id_rsa_**).

#### Template Samples

##### Deploy Command

```
singular deploy template /tmp/deploy.yml  --verbose --timestamp
```

##### Create Nodes In DigitalOcean And Deploy Kuberentes Cluster

```YAML
uri: containerops/demo-for-cncf-ci/deploy-cncf-stack
title: Demo For Deploy Cloud Native Computing Foundation CI Working Group
version: 4
tag: latest
service:
  provider: digitalocean
  token: b516a521b14d86e59c5bb8893
  region: sfo2
  size: 4gb
  image: ubuntu-17-04-x64
tools:
  ssh:
    private: $HOME/.containerops/ssh/id_rsa
    public: $HOME/.containerops/ssh/id_rsa.pub
infras:
  -
    name: etcd
    version: etcd-3.2.2
    master: 3
    minion: 0
    components:
      -
        binary: etcd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/3.2.2/etcd
        package: false
        systemd: etcd-3.2.2
        ca: etcd-3.2.2
      -
        binary: etcdctl
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/3.2.2/etcdctl
        package: false
  -
    name: flannel
    version: flannel-0.7.1
    master: 3
    minion: 0
    dependencies:
      - etcd
    components:
      -
        binary: flanneld
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/0.7.1/flanneld
        package: false
        systemd: flannel-0.7.1
        ca: flannel-0.7.1
        before: "etcdctl --endpoints={{.Nodes}} --ca-file=/etc/kubernetes/ssl/ca.pem --cert-file=/etc/flanneld/ssl/flanneld.pem --key-file=/etc/flanneld/ssl/flanneld-key.pem set /kubernetes/network/config '{\"Network\":\"'172.30.0.0/16'\", \"SubnetLen\": 24, \"Backend\": {\"Type\": \"vxlan\"}}'"
      -
        binary: mk-docker-opts.sh
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/0.7.1/mk-docker-opts.sh
        package: false
  -
    name: docker
    version: docker-17.04.0-ce
    master: 3
    minion: 0
    dependencies:
      - flannel
    components:
      -
        binary: docker
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker
        package: false
        systemd: docker-17.04.0-ce
        before: "iptables -F && iptables -X && iptables -F -t nat && iptables -X -t nat"
        after: "iptables -P FORWARD ACCEPT"
      -
        binary: dockerd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/dockerd
        package: false
      -
        binary: docker-init
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-init
        package: false
      -
        binary: docker-proxy
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-proxy
        package: false
      -
        binary: docker-runc
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-runc
        package: false
      -
        binary: docker-containerd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd
        package: false
      -
        binary: docker-containerd-ctr
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd-ctr
        package: false
      -
        binary: docker-containerd-shim
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd-shim
        package: false
  -
    name: kubernetes
    version: kubernetes-1.6.7
    master: 1
    minion: 3
    dependencies:
      - etcd
      - flannel
      - docker
    components:
      -
        binary: kube-apiserver
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-apiserver
        package: false
        systemd: kube-apiserver-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-controller-manager
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-controller-manager
        package: false
        systemd: kube-controller-manager-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-scheduler
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-scheduler
        package: false
        systemd: kube-scheduler-1.6.7
      -
        binary: kubectl
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kubectl
        package: false
      -
        binary: kubelet
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kubelet
        package: false
        systemd: kubelet-1.6.7
      -
        binary: kube-proxy
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-proxy
        package: false
        systemd: kube-proxy-1.6.7
        ca: kube-proxy-1.6.7
description: WIP
```

##### Using Bare Metals To Deploy Kubernetes Cluster

```
uri: containerops/demo-for-cncf-ci/deploy-cncf-stack
title: Demo For Deploy Cloud Native Computing Foundation CI Working Group
version: 4
tag: latest
nodes:
  -
    ip: 138.68.226.204
    user: root
  -
    ip: 138.68.247.128
    user: root
  -
    ip: 165.227.4.126
    user: root
tools:
  ssh:
    private: /home/meaglith/.ssh/id_rsa
    public: /home/meaglith/.ssh/id_rsa.pub
infras:
  -
    name: etcd
    version: etcd-3.2.2
    master: 3
    minion: 0
    components:
      -
        binary: etcd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/3.2.2/etcd
        package: false
        systemd: etcd-3.2.2
        ca: etcd-3.2.2
      -
        binary: etcdctl
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/3.2.2/etcdctl
        package: false
  -
    name: flannel
    version: flannel-0.7.1
    master: 3
    minion: 0
    dependencies:
      - etcd
    components:
      -
        binary: flanneld
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/0.7.1/flanneld
        package: false
        systemd: flannel-0.7.1
        ca: flannel-0.7.1
        before: "etcdctl --endpoints={{.Nodes}} --ca-file=/etc/kubernetes/ssl/ca.pem --cert-file=/etc/flanneld/ssl/flanneld.pem --key-file=/etc/flanneld/ssl/flanneld-key.pem set /kubernetes/network/config '{\"Network\":\"'172.30.0.0/16'\", \"SubnetLen\": 24, \"Backend\": {\"Type\": \"vxlan\"}}'"
      -
        binary: mk-docker-opts.sh
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/0.7.1/mk-docker-opts.sh
        package: false
  -
    name: docker
    version: docker-17.04.0-ce
    master: 3
    minion: 0
    dependencies:
      - flannel
    components:
      -
        binary: docker
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker
        package: false
        systemd: docker-17.04.0-ce
        before: "iptables -F && iptables -X && iptables -F -t nat && iptables -X -t nat"
        after: "iptables -P FORWARD ACCEPT"
      -
        binary: dockerd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/dockerd
        package: false
      -
        binary: docker-init
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-init
        package: false
      -
        binary: docker-proxy
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-proxy
        package: false
      -
        binary: docker-runc
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-runc
        package: false
      -
        binary: docker-containerd
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd
        package: false
      -
        binary: docker-containerd-ctr
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd-ctr
        package: false
      -
        binary: docker-containerd-shim
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/17.04.0-ce/docker-containerd-shim
        package: false
  -
    name: kubernetes
    version: kubernetes-1.6.7
    master: 1
    minion: 3
    dependencies:
      - etcd
      - flannel
      - docker
    components:
      -
        binary: kube-apiserver
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-apiserver
        package: false
        systemd: kube-apiserver-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-controller-manager
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-controller-manager
        package: false
        systemd: kube-controller-manager-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-scheduler
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-scheduler
        package: false
        systemd: kube-scheduler-1.6.7
      -
        binary: kubectl
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kubectl
        package: false
      -
        binary: kubelet
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kubelet
        package: false
        systemd: kubelet-1.6.7
      -
        binary: kube-proxy
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/1.6.7/kube-proxy
        package: false
        systemd: kube-proxy-1.6.7
        ca: kube-proxy-1.6.7
description: WIP
```

##### Using Local Binary Files To Deploy Kubernetes Clusters

```
uri: containerops/demo-for-cncf-ci/deploy-cncf-stack
title: Demo For Deploy Cloud Native Computing Foundation CI Working Group
version: 4
tag: latest
nodes: 3
service:
  provider: digitalocean
  token: beaa3d89aac777d4d1d
  region: sfo2
  size: 4gb
  image: ubuntu-16-10-x64
infras:
  -
    name: etcd
    version: etcd-3.2.2
    nodes:
      master: 3
      node: 0
    components:
      -
        binary: etcd
        url: /root/binary/etcd/3.2.2/etcd
        package: false
        systemd: etcd-3.2.2
        ca: etcd-3.2.2
      -
        binary: etcdctl
        url: /root/binary/etcd/3.2.2/etcdctl
        package: false
  -
    name: flannel
    version: flannel-0.7.1
    nodes:
      master: 3
      node: 0
    dependencies:
      - etcd
    components:
      -
        binary: flanneld
        url: /root/binary/flannel/0.7.1/flanneld
        package: false
        systemd: flannel-0.7.1
        ca: flannel-0.7.1
        before: "etcdctl --endpoints={{.Nodes}} --ca-file=/etc/kubernetes/ssl/ca.pem --cert-file=/etc/flanneld/ssl/flanneld.pem --key-file=/etc/flanneld/ssl/flanneld-key.pem set /kubernetes/network/config '{\"Network\":\"'172.30.0.0/16'\", \"SubnetLen\": 24, \"Backend\": {\"Type\": \"vxlan\"}}'"
      -
        binary: mk-docker-opts.sh
        url: /root/binary/flannel/0.7.1/mk-docker-opts.sh
        package: false
  -
    name: docker
    version: docker-17.04.0-ce
    nodes:
      master: 3
      node: 0
    dependencies:
      - flannel
    components:
      -
        binary: docker
        url: /root/binary/docker/17.04.0-ce/docker
        package: false
        systemd: docker-17.04.0-ce
        before: "iptables -F && iptables -X && iptables -F -t nat && iptables -X -t nat"
        after: "iptables -P FORWARD ACCEPT"
      - 
        binary: dockerd
        url: /root/binary/docker/17.04.0-ce/dockerd
        package: false
      -
        binary: docker-init
        url: /root/binary/docker/17.04.0-ce/docker-init
        package: false
      -
        binary: docker-proxy
        url: /root/binary/docker/17.04.0-ce/docker-proxy
        package: false
      -
        binary: docker-runc
        url: /root/binary/docker/17.04.0-ce/docker-runc
        package: false
      -
        binary: docker-containerd
        url: /root/binary/docker/17.04.0-ce/docker-containerd
        package: false
      -
        binary: docker-containerd-ctr
        url: /root/binary/docker/17.04.0-ce/docker-containerd-ctr
        package: false
      -
        binary: docker-containerd-shim
        url: /root/binary/docker/17.04.0-ce/docker-containerd-shim
        package: false
  -
    name: kubernetes
    version: kubernetes-1.6.7
    nodes:
      master: 1
      node: 3
    dependencies:
      - etcd
      - flannel
      - docker
    components:
      -
        binary: kube-apiserver
        url: /root/binary/kubernetes/1.6.7/kube-apiserver
        package: false
        systemd: kube-apiserver-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-controller-manager
        url: /root/binary/kubernetes/1.6.7/kube-controller-manager
        package: false
        systemd: kube-controller-manager-1.6.7
        ca: kubernetes-1.6.7
      -
        binary: kube-scheduler
        url: /root/binary/kubernetes/1.6.7/kube-scheduler
        package: false
        systemd: kube-scheduler-1.6.7
      -
        binary: kubectl
        url: /root/binary/kubernetes/1.6.7/kubectl
        package: false
      -
        binary: kubelet
        url: /root/binary/kubernetes/1.6.7/kubelet
        package: false
        systemd: kubelet-1.6.7
      -
        binary: kube-proxy
        url: /root/binary/kubernetes/1.6.7/kube-proxy
        package: false
        systemd: kube-proxy-1.6.7
        ca: kube-proxy-1.6.7
description: WIP
```