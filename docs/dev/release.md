
# Release guide

This guide covers releasing new versions of matchbox.

## Version

Create a release commit which updates old version references.

```sh
$ export VERSION=v0.10.0
```

## Tag

Tag, sign the release version, and push it to Github.

```sh
$ git tag -s vX.Y.Z -m 'vX.Y.Z'
$ git push origin --tags
$ git push origin main
```

## Images

Travis CI will build the Docker image and push it to Quay.io when the tag is pushed to master. Verify the new image and version.

```sh
$ sudo docker run quay.io/poseidon/matchbox:$VERSION -version
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
./_output/matchbox-v0.10.0-linux-amd64/matchbox -version
```

## Signing

Release tarballs are signed by Dalton Hubble's GPG [Key](/docs/deployment.md#download)

```sh
make release-sign
```

Verify the signatures.

```sh
make release-verify
```

## Publish

Upload the signed tarball(s) with the Github release. Promote the release from a `pre-release` to an official release.
