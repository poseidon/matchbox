
# Release Guide

This guide covers releasing new versions of matchbox.

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

    sudo docker run quay.io/coreos/matchbox:$VERSION -version
    sudo rkt run --no-store quay.io/coreos/matchbox:$VERSION -- -version

## Github Release

Publish the release on Github with release notes.

## Tarballs

Build the release tarballs.

    make release

Verify the reported version.

    ./_output/matchbox-v0.4.2-linux-amd64/matchbox -version

## ACI

Build the rkt ACI on a Linux host with `acbuild`,

    make aci

Check that the listed version is correct/clean.

    sudo rkt --insecure-options=image run matchbox.aci -- -version

Add the ACI to `output` for signing.

    mv matchbox.aci _output/matchbox-$VERSION-linux-amd64.aci

## Signing

Sign the release tarballs and ACI with a [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) subkey.

    cd _output
    gpg2 -a --default-key FC8A365E --detach-sign matchbox-$VERSION-linux-amd64.aci
    gpg2 -a --default-key FC8A365E --detach-sign matchbox-$VERSION-linux-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign matchbox-$VERSION-darwin-amd64.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign matchbox-$VERSION-linux-arm.tar.gz
    gpg2 -a --default-key FC8A365E --detach-sign matchbox-$VERSION-linux-arm64.tar.gz

Verify the signatures.

    gpg2 --verify matchbox-$VERSION-linux-amd64.aci.asc matchbox-$VERSION-linux-amd64.aci
    gpg2 --verify matchbox-$VERSION-linux-amd64.tar.gz.asc matchbox-$VERSION-linux-amd64.tar.gz
    gpg2 --verify matchbox-$VERSION-darwin-amd64.tar.gz.asc matchbox-$VERSION-darwin-amd64.tar.gz
    gpg2 --verify matchbox-$VERSION-linux-arm.tar.gz.asc matchbox-$VERSION-linux-arm.tar.gz
    gpg2 --verify matchbox-$VERSION-linux-arm64.tar.gz.asc matchbox-$VERSION-linux-arm64.tar.gz

## Publish

Upload the signed tarball(s) and ACI with the Github release. Promote the release from a `pre-release` to an official release.
