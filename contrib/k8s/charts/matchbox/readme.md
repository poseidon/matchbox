# Matchbox services
Network boot and provision Container Linux clusters (e.g. etcd3, Kubernetes, more...)

# Default values
* matchbox.tag: latest
* matchbox.stage: dev
* matchbox.log_level: debug
* matchbox.ingress_name: DNS record of the matchbox service (matchbox.example.com) 
* matchbox.ingress_rpc_name: DNS record of the matchbox RPC service (matchbox-rpc.example.com)
* matchbox.resources.cpu_limit: Max CPU resource of the matchbox pod (50m)
* matchbox.resources.mem_limit: Max MEM resource of the matchbox pod (50Mi)
* matchbox.ssl.ca: default CA for matchbox.example.com, matchbox-rpc.example.com
* matchbox.ssl.crt: default CRT for matchbox.example.com, matchbox-rpc.example.com
* matchbox.ssl.key: default KEY for matchbox.example.com, matchbox-rpc.example.com
* matchbox.pvc.data.size: Size of the volume claim for data directory (2Gi)
* matchbox.pvc.data.annotations: Additional annotations for the volume claim data
* matchbox.pvc.assets.size: Size of the volume claim for data directory (8Gi)
* matchbox.pvc.assets.annotations: Additional annotations for the volume claim assets 

# SSL certificate
```
# Generate the cert
# and export to env variables

$ export CA=`cat ca.crt | base64 -w0`
$ export CRT=`cat server.crt | base64 -w0`
$ export KEY=`cat server.key | base64 -w0`

$ helm install --name matchbox --set matchbox.ssl.ca=$CA --set matchbox.ssl.crt=$CRT --set matchbox.ssl.key=$KEY ./matchbox

$ export CA=
$ export CRE=
$ export key=
```

# Installation and upgrade
```bash
$ git pull git@github.com:coreos/matchbox.git
$ cd matchbox/contrib/k8s/charts/
$ helm inspect ./matchbox > matchbox-dev.yaml
$ vi matchbox.yaml
$ helm upgrade --install -f matchbox-dev.yaml matchbox-dev ./matchbox
```
