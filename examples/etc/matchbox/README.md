
## gRPC API Credentials

Create FAKE TLS credentials for running the `matchbox` gRPC API examples.

**DO NOT** use these certificates for anything other than running `matchbox` examples. Use your organization's production PKI for production deployments.

Navigate to the example directory which will be mounted as `/etc/matchbox` in examples:

    cd matchbox/examples/etc/matchbox

Set certificate subject alt names which should be used by exporting `SAN`. Use the DNS name or IP at which `matchbox` is hosted.

    # for examples on metal0 or docker0 bridges
    export SAN=IP.1:127.0.0.1,IP.2:172.18.0.2

    # production example
    export SAN=DNS.1:matchbox.example.com

Create a fake `ca.crt`, `server.crt`, `server.key`, `client.crt`, and `client.key`. Type 'Y' when prompted.

    $ ./cert-gen
    Creating FAKE CA, server cert/key, and client cert/key...
    ...
    ...
    ...
    ******************************************************************
    WARNING: Generated TLS credentials are ONLY SUITABLE FOR EXAMPLES!
    Use your organization's production PKI for production deployments!

## Inpsect

Inspect the generated FAKE certificates if desired.

    openssl x509 -noout -text -in ca.crt
    openssl x509 -noout -text -in server.crt
    openssl x509 -noout -text -in client.crt

## Verify

Verify that the FAKE server and client certificates were signed by the fake CA.

    openssl verify -CAfile ca.crt server.crt
    openssl verify -CAfile ca.crt client.crt