
## gRPC API Credentials

Create FAKE TLS credentials for running the `bootcfg` gRPC API examples.

**DO NOT** use these certificates for anything other than running `bootcfg` examples. Use your organization's production PKI for production deployments.

Navigate to the example directory which will be mounted as `/etc/bootcfg` in examples:

    cd coreos-baremetal/examples/etc/bootcfg

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