
# Kubernetes

## TLS Assets

Use the `generate-tls` script to generate throw-away TLS assets. The script will generate a root CA and `admin`, `apiserver`, and `worker` certificates in `assets/tls`.

    cd coreos-baremetal
    ./examples/kubernetes/scripts/generate-tls

Alternately, if you have existing Public Key Infrastructure, add your CA certificate, entity certificates, and entity private keys to `assets/tls`.

    * ca.pem
    * apiserver.pem
    * apiserver-key.pem
    * worker.pem
    * worker-key.pem
    * admin.pem
    * admin-key.pem

See the [Cluster TLS OpenSSL Generation](https://coreos.com/kubernetes/docs/latest/openssl.html) document or [Kubernetes Step by Step](https://coreos.com/kubernetes/docs/latest/getting-started.html) for more details.
