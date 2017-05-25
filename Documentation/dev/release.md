
# Release guide

This guide covers releasing new versions of matchbox.

## Version

Create a release commit which updates old version references.

```sh
$ export VERSION=v0.6.1
```

## Tag

Tag, sign the release version, and push it to Github.

```sh
$ git tag -s vX.Y.Z -m 'vX.Y.Z'
$ git push origin --tags
$ git push origin master
```

## Images

Travis CI will build the Docker image and push it to Quay.io when the tag is pushed to master. Verify the new image and version.

```sh
$ sudo docker run quay.io/coreos/matchbox:$VERSION -version
$ sudo rkt run --no-store quay.io/coreos/matchbox:$VERSION -- -version
```

## Github release

Publish the release on Github with release notes.

## Tarballs

Build the release tarballs.

```sh
$ make release
```

Verify the reported version.

```
./_output/matchbox-v0.6.1-linux-amd64/matchbox -version
```

## Signing

Sign the release tarballs and ACI with a [CoreOS App Signing Key](https://coreos.com/security/app-signing-key/) subkey.

```sh
cd _output
gpg2 --armor --local-user A6F71EE5BEDDBA18! --detach-sign matchbox-$VERSION-linux-amd64.tar.gz
gpg2 --armor --local-user A6F71EE5BEDDBA18! --detach-sign matchbox-$VERSION-darwin-amd64.tar.gz
gpg2 --armor --local-user A6F71EE5BEDDBA18! --detach-sign matchbox-$VERSION-linux-arm.tar.gz
gpg2 --armor --local-user A6F71EE5BEDDBA18! --detach-sign matchbox-$VERSION-linux-arm64.tar.gz
```

Verify the signatures.

```sh
gpg2 --verify matchbox-$VERSION-linux-amd64.tar.gz.asc matchbox-$VERSION-linux-amd64.tar.gz
gpg2 --verify matchbox-$VERSION-darwin-amd64.tar.gz.asc matchbox-$VERSION-darwin-amd64.tar.gz
gpg2 --verify matchbox-$VERSION-linux-arm.tar.gz.asc matchbox-$VERSION-linux-arm.tar.gz
gpg2 --verify matchbox-$VERSION-linux-arm64.tar.gz.asc matchbox-$VERSION-linux-arm64.tar.gz
```

## Publish

Upload the signed tarball(s) with the Github release. Promote the release from a `pre-release` to an official release.
