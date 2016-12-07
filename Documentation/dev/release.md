
# Release Guide

This guide covers releasing new versions of coreos-baremetal.

## Version

Create a release commit which updates old version references.

    export VERSION=v0.4.2

## Tag

Tag, sign the release version, and push it to Github.

    git tag -s vX.Y.Z -m 'vX.Y.Z'
    git push origin --tags

Travis CI will build the Docker image and push it to Quay.io when the tag is pushed to master.

## Github Release

Publish the release on Github with release notes.

## Tarballs

Build the release tarballs.

    make release

## ACI

Build the rkt ACI on a Linux host with `acbuild`,

    make build
    ./build-aci
    mv bootcfg.aci _output/bootcfg-$VERSION-linux-amd64.aci

Check that the listed version is correct/clean.

    sudo rkt --insecure-options=image run bootcfg.aci -- -version

## Signing

Sign the release tarballs and ACI with a [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) subkey.

    cd _output
    gpg2 -a --default-key FC8A365E --detach-sign bootcfg-$VERSION-linux-amd64.aci
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-linux-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-darwin-amd64.tar.gz

Verify the signatures.

    gpg2 --verify bootcfg-$VERSION-linux-amd64.aci.asc bootcfg-$VERSION-linux-amd64.aci
    gpg2 --verify coreos-baremetal-$VERSION-linux-amd64.tar.gz.asc coreos-baremetal-$VERSION-linux-amd64.tar.gz
    gpg2 --verify coreos-baremetal-$VERSION-darwin-amd64.tar.gz.asc coreos-baremetal-$VERSION-darwin-amd64.tar.gz

## Publish

Publish the signed tarball(s) and ACI with the Github release. Verify the Docker image was published to Quay.io when the release tag was pushed.
