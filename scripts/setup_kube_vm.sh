#!/bin/bash

GOPATH=${HOME}/goproj
KPATH=${GOPATH}/src/k8s.io
ETCD_VER=v3.0.15

function install_docker {
  sudo docker version
  if [ $? -ne 0 ]; then
    sudo apt-get install docker-engine
    sudo groupadd docker
    sudo gpasswd -a ${USER} docker
    sudo service docker restart
    echo added ${USER} to group docker, you might need to log out for this to take effect
  fi
}

function install_kubectl {
  kubectl
  if [ $? -ne 0 ]; then
     curl -Lo kubectl http://storage.googleapis.com/kubernetes-release/release/v1.5.1/bin/linux/amd64/kubectl 
     sudo chmod +x kubectl
     sudo mv kubectl /usr/local/bin/
  fi
}

function install_etcd {
  etcd --version
  if [ $? -ne 0 ]; then
    DOWNLOAD_URL=https://github.com/coreos/etcd/releases/download
    curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
    mkdir -p /tmp/test-etcd && tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/test-etcd --strip-components=1
    cd /tmp/test-etcd
    sudo mv etcd /usr/local/bin
    sudo mv etcdctl /usr/local/bin
  fi
}

function install_cfssl {
  cfssl version
  if [ $? -ne 0 ]; then
    go get -u github.com/cloudflare/cfssl/cmd/cfssl
    go get -u github.com/cloudflare/cfssl/cmd/..
    echo installed cfssl, this assumes ${GOPATH}/bin is in your path
  fi
}

function clone_kube {
  ls ${KPATH}/kubernetes
  if [ $? -ne 0 ]; then
    cd ${HOME}
    mkdir -p ${KPATH}
    git clone https://github.com/kubernetes/kubernetes.git
  fi
}

echo checking docker
install_docker
echo checking kubectl
install_kubectl
echo checking etcd
install_etcd
echo checking cfssl
install_cfssl
echo checking kubernetes
clone_kube
echo starting kube
$KPATH/kubernetes/hack/local-up-cluster.sh

# Install kubectl
# Install etcd
# Install cfssl
# Clone kube repo


