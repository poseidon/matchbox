
# bootcfg Release Guide

This guide covers releasing new versions of `bootcfg`.

## Prepare

First, build the binary and ACI. Check that their version is correct and clean.

## Release Notes

Create a pre-release and draft the release notes based on the [changelog](../CHANGES.md).

## Tag

Tag and sign the release version.

    git tag -s vX.Y.Z -m 'vX.Y.Z'

Travis CI will build the Docker image and push it to Quay.io when the tag is pushed to master.

## Binaries and Images

Build the binary and ACI. Prepare the binary tarball and ACI.

    export VERSION=v0.3.0
    mkdir bootcfg-$VERSION
    cp bin/bootcfg bootcfg-$VERSION
    cp bootcfg.aci bootcfg-$VERSION-linux-amd64.aci
    tar -zcvf bootcfg-$VERSION-linux-amd64.tar.gz bootcfg-VERSION

## Signing

Sign the binary tarball and ACI.

    gpg2 -a --default-key FC8A365E --detach-sign bootcfg-$VERSION-linux-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign bootcfg-$VERSION-linux-amd64.aci

Verify the signatures.

    gpg2 --verify bootcfg-$VERSION-linux-amd64.tar.gz.asc bootcfg-$VERSION-linux-amd64.tar.gz
    gpg2 --verify bootcfg-$VERSION-linux-amd64.aci.asc bootcfg-$VERSION-linux-amd64.aci

## Publish

Publish the signed binary tarball(s) and the signed ACI with the Github release. The Docker image is published to Quay.io when the tag is pushed to master.
