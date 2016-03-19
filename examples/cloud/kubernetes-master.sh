#!/bin/bash
set -e

# List of etcd servers (http://ip:port), comma separated
export ETCD_ENDPOINTS={{.k8s_etcd_endpoints}}

# Specify the version (vX.Y.Z) of Kubernetes assets to deploy
export K8S_VER=v1.1.8_coreos.0

# The CIDR network to use for pod IPs.
# Each pod launched in the cluster will be assigned an IP out of this range.
# Each node will be configured such that these IPs will be routable using the flannel overlay network.
export POD_NETWORK={{.k8s_pod_network}}

# The CIDR network to use for service cluster IPs.
# Each service will be assigned a cluster IP out of this range.
# This must not overlap with any IP ranges assigned to the POD_NETWORK, or other existing network infrastructure.
# Routing to these IPs is handled by a proxy service local to each node, and are not required to be routable between nodes.
export SERVICE_IP_RANGE={{.k8s_service_ip_range}}

# The IP address of the Kubernetes API Service
# If the SERVICE_IP_RANGE is changed above, this must be set to the first IP in that range.
export K8S_SERVICE_IP={{.k8s_service_ip}}

# The IP address of the cluster DNS service.
# This IP must be in the range of the SERVICE_IP_RANGE and cannot be the first IP in the range.
# This same IP must be configured on all worker nodes to enable DNS service discovery.
export DNS_SERVICE_IP={{.k8s_dns_service_ip}}

# ADVERTISE_IP is the host node's IP.
export ADVERTISE_IP={{.k8s_advertise_ip}}

# TLS Certificate assets are hosted by the Config Server
export CERT_ENDPOINT={{.k8s_cert_endpoint}}

function init_config {
  local REQUIRED=('ADVERTISE_IP' 'POD_NETWORK' 'ETCD_ENDPOINTS' 'SERVICE_IP_RANGE' 'K8S_SERVICE_IP' 'DNS_SERVICE_IP' 'K8S_VER' )
  
  for REQ in "${REQUIRED[@]}"; do
      if [ -z "$(eval echo \$$REQ)" ]; then
          echo "Missing required config value: ${REQ}"
          exit 1
      fi
  done
}

function get_certs {
  DEST=/etc/kubernetes/ssl
  mkdir -p $DEST
  curl $CERT_ENDPOINT/tls/apiserver.pem -o $DEST/apiserver.pem
  curl $CERT_ENDPOINT/tls/apiserver-key.pem -o $DEST/apiserver-key.pem
  curl $CERT_ENDPOINT/tls/ca.pem -o $DEST/ca.pem
}

function init_flannel {
    echo "Waiting for etcd..."
    while true
    do
        IFS=',' read -ra ES <<< "$ETCD_ENDPOINTS"
        for ETCD in "${ES[@]}"; do
            echo "Trying: $ETCD"
            if [ -n "$(curl --silent "$ETCD/v2/machines")" ]; then
                local ACTIVE_ETCD=$ETCD
                break
            fi
            sleep 1
        done
        if [ -n "$ACTIVE_ETCD" ]; then
            break
        fi
    done
    RES=$(curl --silent -X PUT -d "value={\"Network\":\"$POD_NETWORK\",\"Backend\":{\"Type\":\"vxlan\"}}" "$ACTIVE_ETCD/v2/keys/coreos.com/network/config?prevExist=false")
    if [ -z "$(echo $RES | grep '"action":"create"')" ] && [ -z "$(echo $RES | grep 'Key already exists')" ]; then
        echo "Unexpected error configuring flannel pod network: $RES"
    fi
}

function init_templates {
    local TEMPLATE=/etc/systemd/system/kubelet.service
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
[Service]
ExecStartPre=/usr/bin/mkdir -p /etc/kubernetes/manifests
Environment=KUBELET_VERSION=${K8S_VER}
ExecStart=/usr/lib/coreos/kubelet-wrapper \
  --api_servers=http://127.0.0.1:8080 \
  --register-node=false \
  --allow-privileged=true \
  --config=/etc/kubernetes/manifests \
  --hostname-override=${ADVERTISE_IP} \
  --cluster_dns=${DNS_SERVICE_IP} \
  --cluster_domain=cluster.local
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
    }

    local TEMPLATE=/etc/kubernetes/manifests/kube-proxy.yaml
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
apiVersion: v1
kind: Pod
metadata:
  name: kube-proxy
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-proxy
    image: quay.io/coreos/hyperkube:$K8S_VER
    command:
    - /hyperkube
    - proxy
    - --master=http://127.0.0.1:8080
    - --proxy-mode=iptables
    securityContext:
      privileged: true
    volumeMounts:
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
    }

    local TEMPLATE=/etc/kubernetes/manifests/kube-apiserver.yaml
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
apiVersion: v1
kind: Pod
metadata:
  name: kube-apiserver
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-apiserver
    image: quay.io/coreos/hyperkube:$K8S_VER
    command:
    - /hyperkube
    - apiserver
    - --bind-address=0.0.0.0
    - --etcd-servers=${ETCD_ENDPOINTS}
    - --allow-privileged=true
    - --service-cluster-ip-range=${SERVICE_IP_RANGE}
    - --secure-port=443
    - --advertise-address=${ADVERTISE_IP}
    - --admission-control=NamespaceLifecycle,NamespaceExists,LimitRanger,SecurityContextDeny,ServiceAccount,ResourceQuota
    - --tls-cert-file=/etc/kubernetes/ssl/apiserver.pem
    - --tls-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --client-ca-file=/etc/kubernetes/ssl/ca.pem
    - --service-account-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    ports:
    - containerPort: 443
      hostPort: 443
      name: https
    - containerPort: 8080
      hostPort: 8080
      name: local
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
    }

    local TEMPLATE=/etc/kubernetes/manifests/kube-podmaster.yaml
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
apiVersion: v1
kind: Pod
metadata:
  name: kube-podmaster
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: scheduler-elector
    image: gcr.io/google_containers/podmaster:1.1
    command:
    - /podmaster
    - --etcd-servers=${ETCD_ENDPOINTS}
    - --key=scheduler
    - --whoami=${ADVERTISE_IP}
    - --source-file=/src/manifests/kube-scheduler.yaml
    - --dest-file=/dst/manifests/kube-scheduler.yaml
    volumeMounts:
    - mountPath: /src/manifests
      name: manifest-src
      readOnly: true
    - mountPath: /dst/manifests
      name: manifest-dst
  - name: controller-manager-elector
    image: gcr.io/google_containers/podmaster:1.1
    command:
    - /podmaster
    - --etcd-servers=${ETCD_ENDPOINTS}
    - --key=controller
    - --whoami=${ADVERTISE_IP}
    - --source-file=/src/manifests/kube-controller-manager.yaml
    - --dest-file=/dst/manifests/kube-controller-manager.yaml
    terminationMessagePath: /dev/termination-log
    volumeMounts:
    - mountPath: /src/manifests
      name: manifest-src
      readOnly: true
    - mountPath: /dst/manifests
      name: manifest-dst
  volumes:
  - hostPath:
      path: /srv/kubernetes/manifests
    name: manifest-src
  - hostPath:
      path: /etc/kubernetes/manifests
    name: manifest-dst
EOF
    }

    local TEMPLATE=/srv/kubernetes/manifests/kube-controller-manager.yaml
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
apiVersion: v1
kind: Pod
metadata:
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
  - name: kube-controller-manager
    image: quay.io/coreos/hyperkube:$K8S_VER
    command:
    - /hyperkube
    - controller-manager
    - --master=http://127.0.0.1:8080
    - --service-account-private-key-file=/etc/kubernetes/ssl/apiserver-key.pem
    - --root-ca-file=/etc/kubernetes/ssl/ca.pem
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10252
      initialDelaySeconds: 15
      timeoutSeconds: 1
    volumeMounts:
    - mountPath: /etc/kubernetes/ssl
      name: ssl-certs-kubernetes
      readOnly: true
    - mountPath: /etc/ssl/certs
      name: ssl-certs-host
      readOnly: true
  hostNetwork: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/ssl
    name: ssl-certs-kubernetes
  - hostPath:
      path: /usr/share/ca-certificates
    name: ssl-certs-host
EOF
    }

    local TEMPLATE=/srv/kubernetes/manifests/kube-scheduler.yaml
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
apiVersion: v1
kind: Pod
metadata:
  name: kube-scheduler
  namespace: kube-system
spec:
  hostNetwork: true
  containers:
  - name: kube-scheduler
    image: quay.io/coreos/hyperkube:$K8S_VER
    command:
    - /hyperkube
    - scheduler
    - --master=http://127.0.0.1:8080
    livenessProbe:
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10251
      initialDelaySeconds: 15
      timeoutSeconds: 1
EOF
    }

    local TEMPLATE=/srv/kubernetes/manifests/kube-system.json
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
{
  "apiVersion": "v1",
  "kind": "Namespace",
  "metadata": {
    "name": "kube-system"
  }
}
EOF
    }

    local TEMPLATE=/srv/kubernetes/manifests/kube-dns-rc.json
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
{
    "apiVersion": "v1",
    "kind": "ReplicationController",
    "metadata": {
        "labels": {
            "k8s-app": "kube-dns",
            "kubernetes.io/cluster-service": "true",
            "version": "v9"
        },
        "name": "kube-dns-v9",
        "namespace": "kube-system"
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "k8s-app": "kube-dns",
            "version": "v9"
        },
        "template": {
            "metadata": {
                "labels": {
                    "k8s-app": "kube-dns",
                    "kubernetes.io/cluster-service": "true",
                    "version": "v9"
                }
            },
            "spec": {
                "containers": [
                    {
                        "command": [
                            "/usr/local/bin/etcd",
                            "-data-dir",
                            "/var/etcd/data",
                            "-listen-client-urls",
                            "http://127.0.0.1:2379,http://127.0.0.1:4001",
                            "-advertise-client-urls",
                            "http://127.0.0.1:2379,http://127.0.0.1:4001",
                            "-initial-cluster-token",
                            "skydns-etcd"
                        ],
                        "image": "gcr.io/google_containers/etcd:2.0.9",
                        "name": "etcd",
                        "resources": {
                            "limits": {
                                "cpu": "100m",
                                "memory": "50Mi"
                            }
                        },
                        "volumeMounts": [
                            {
                                "mountPath": "/var/etcd/data",
                                "name": "etcd-storage"
                            }
                        ]
                    },
                    {
                        "args": [
                            "-domain=cluster.local"
                        ],
                        "image": "gcr.io/google_containers/kube2sky:1.11",
                        "name": "kube2sky",
                        "resources": {
                            "limits": {
                                "cpu": "100m",
                                "memory": "50Mi"
                            }
                        }
                    },
                    {
                        "args": [
                            "-machines=http://127.0.0.1:4001",
                            "-addr=0.0.0.0:53",
                            "-ns-rotate=false",
                            "-domain=cluster.local."
                        ],
                        "image": "gcr.io/google_containers/skydns:2015-10-13-8c72f8c",
                        "livenessProbe": {
                            "httpGet": {
                                "path": "/healthz",
                                "port": 8080,
                                "scheme": "HTTP"
                            },
                            "initialDelaySeconds": 30,
                            "timeoutSeconds": 5
                        },
                        "name": "skydns",
                        "ports": [
                            {
                                "containerPort": 53,
                                "name": "dns",
                                "protocol": "UDP"
                            },
                            {
                                "containerPort": 53,
                                "name": "dns-tcp",
                                "protocol": "TCP"
                            }
                        ],
                        "readinessProbe": {
                            "httpGet": {
                                "path": "/healthz",
                                "port": 8080,
                                "scheme": "HTTP"
                            },
                            "initialDelaySeconds": 1,
                            "timeoutSeconds": 5
                        },
                        "resources": {
                            "limits": {
                                "cpu": "100m",
                                "memory": "50Mi"
                            }
                        }
                    },
                    {
                        "args": [
                            "-cmd=nslookup kubernetes.default.svc.cluster.local localhost >/dev/null",
                            "-port=8080"
                        ],
                        "image": "gcr.io/google_containers/exechealthz:1.0",
                        "name": "healthz",
                        "ports": [
                            {
                                "containerPort": 8080,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {
                            "limits": {
                                "cpu": "10m",
                                "memory": "20Mi"
                            }
                        }
                    }
                ],
                "dnsPolicy": "Default",
                "volumes": [
                    {
                        "emptyDir": {},
                        "name": "etcd-storage"
                    }
                ]
            }
        }
    }
}
EOF
    }

    local TEMPLATE=/srv/kubernetes/manifests/kube-dns-svc.json
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "kube-dns",
    "namespace": "kube-system",
    "labels": {
      "k8s-app": "kube-dns",
      "kubernetes.io/name": "KubeDNS",
      "kubernetes.io/cluster-service": "true"
    }
  },
  "spec": {
    "clusterIP": "$DNS_SERVICE_IP",
    "ports": [
      {
        "protocol": "UDP",
        "name": "dns",
        "port": 53
      },
      {
        "protocol": "TCP",
        "name": "dns-tcp",
        "port": 53
      }
    ],
    "selector": {
      "k8s-app": "kube-dns"
    }
  }
}
EOF
    }

    local TEMPLATE=/etc/flannel/options.env
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
FLANNELD_IFACE=$ADVERTISE_IP
FLANNELD_ETCD_ENDPOINTS=$ETCD_ENDPOINTS
EOF
    }

    local TEMPLATE=/etc/systemd/system/flanneld.service.d/40-ExecStartPre-symlink.conf.conf
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
[Service]
ExecStartPre=/usr/bin/ln -sf /etc/flannel/options.env /run/flannel/options.env
EOF
    }

    local TEMPLATE=/etc/systemd/system/docker.service.d/40-flannel.conf
    [ -f $TEMPLATE ] || {
        echo "TEMPLATE: $TEMPLATE"
        mkdir -p $(dirname $TEMPLATE)
        cat << EOF > $TEMPLATE
[Unit]
Requires=flanneld.service
After=flanneld.service
EOF
    }

}

function start_addons {
  echo "Waiting for Kubernetes API..."
  until curl --silent "http://127.0.0.1:8080/version"
  do
    sleep 5
  done
  echo
    echo "K8S: kube-system namespace"
    curl --silent -XPOST -d"$(cat /srv/kubernetes/manifests/kube-system.json)" "http://127.0.0.1:8080/api/v1/namespaces" > /dev/null
    echo "K8S: DNS addon"
    curl --silent -XPOST -d"$(cat /srv/kubernetes/manifests/kube-dns-rc.json)" "http://127.0.0.1:8080/api/v1/namespaces/kube-system/replicationcontrollers" > /dev/null
    curl --silent -XPOST -d"$(cat /srv/kubernetes/manifests/kube-dns-svc.json)" "http://127.0.0.1:8080/api/v1/namespaces/kube-system/services" > /dev/null
}

init_config
get_certs
init_templates
init_flannel

systemctl stop update-engine; systemctl mask update-engine
systemctl daemon-reload
systemctl enable kubelet; systemctl start kubelet

start_addons

echo "done" > /home/core/master
