
# Release Guide

This guide covers releasing new versions of coreos-baremetal.

## Version

Create a release commit which updates old version references.

    export VERSION=v0.4.2

## Tag

Tag, sign the release version, and push it to Github.

    git tag -s vX.Y.Z -m 'vX.Y.Z'
    git push origin --tags
    git push origin master

## Images

Travis CI will build the Docker image and push it to Quay.io when the tag is pushed to master. Verify the new image and version.

    sudo docker run quay.io/coreos/bootcfg:$VERSION -version
    sudo rkt run --no-store quay.io/coreos/bootcfg:$VERSION -- -version

## Github Release

Publish the release on Github with release notes.

## Tarballs

Build the release tarballs.

    make build
    make release

Verify the reported version.

    ./_output/coreos-baremetal-v0.4.2-linux-amd64/bootcfg -version

## ACI

Build the rkt ACI on a Linux host with `acbuild`,

    ./build-aci

Check that the listed version is correct/clean.

    sudo rkt --insecure-options=image run bootcfg.aci -- -version

Add the ACI to `output` for signing.

    mv bootcfg.aci _output/bootcfg-$VERSION-linux-amd64.aci

## Signing

Sign the release tarballs and ACI with a [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) subkey.

    cd _output
    gpg2 -a --default-key FC8A365E --detach-sign bootcfg-$VERSION-linux-amd64.aci
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-linux-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-darwin-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-linux-arm.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign coreos-baremetal-$VERSION-linux-arm64.tar.gz

Verify the signatures.

    gpg2 --verify bootcfg-$VERSION-linux-amd64.aci.asc bootcfg-$VERSION-linux-amd64.aci
    gpg2 --verify coreos-baremetal-$VERSION-linux-amd64.tar.gz.asc coreos-baremetal-$VERSION-linux-amd64.tar.gz
    gpg2 --verify coreos-baremetal-$VERSION-darwin-amd64.tar.gz.asc coreos-baremetal-$VERSION-darwin-amd64.tar.gz
    gpg2 --verify coreos-baremetal-$VERSION-linux-arm.tar.gz.asc coreos-baremetal-$VERSION-linux-arm.tar.gz
    gpg2 --verify coreos-baremetal-$VERSION-linux-arm64.tar.gz.asc coreos-baremetal-$VERSION-linux-arm64.tar.gz

## Publish

Upload the signed tarball(s) and ACI with the Github release. Promote the release from a `pre-release` to an official release.
