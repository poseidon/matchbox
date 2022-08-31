#!/usr/bin/env bash
# USAGE: ./scripts/get-flatcar
# USAGE: ./scripts/get-flatcar channel version dest
#
# ENV VARS:
# - OEM_ID - specify OEM image id to download, alongside the default one
set -eou pipefail

GPG=${GPG:-/usr/bin/gpg}

CHANNEL=${1:-"stable"}
VERSION=${2:-"current"}
DEST_DIR=${3:-"$PWD/examples/assets"}
OEM_ID=${OEM_ID:-""}
DEST=$DEST_DIR/flatcar/$VERSION
BASE_URL=https://$CHANNEL.release.flatcar-linux.net/amd64-usr/$VERSION

# check channel/version exist based on the header response
if ! curl -s -I "${BASE_URL}/flatcar_production_pxe.vmlinuz" | grep -q -E '^HTTP/[0-9.]+ [23][0-9][0-9]'; then
  echo "Channel or Version not found"
  exit 1
fi

if [[ ! -d "$DEST" ]]; then
  echo "Creating directory ${DEST}"
  mkdir -p "${DEST}"
fi

if [[ -n "${OEM_ID}" ]]; then
  IMAGE_NAME="flatcar_production_${OEM_ID}_image.bin.bz2"

  # check if the oem version exists based on the header response
  if ! curl -s -I "${BASE_URL}/${IMAGE_NAME}" | grep -q -E '^HTTP/[0-9.]+ [23][0-9][0-9]'; then
    echo "OEM version not found"
    exit 1
  fi
fi

echo "Downloading Flatcar Linux $CHANNEL $VERSION images and sigs to $DEST"

echo "Flatcar Linux Image Signing Key"
curl -# https://www.flatcar.org/security/image-signing-key/Flatcar_Image_Signing_Key.asc -o "${DEST}/Flatcar_Image_Signing_Key.asc"
$GPG --import <"$DEST/Flatcar_Image_Signing_Key.asc" || true

# Version
echo "version.txt"
curl -# "${BASE_URL}/version.txt" -o "${DEST}/version.txt"

# PXE kernel and sig
echo "flatcar_production_pxe.vmlinuz..."
curl -# "${BASE_URL}/flatcar_production_pxe.vmlinuz" -o "${DEST}/flatcar_production_pxe.vmlinuz"
echo "flatcar_production_pxe.vmlinuz.sig"
curl -# "${BASE_URL}/flatcar_production_pxe.vmlinuz.sig" -o "${DEST}/flatcar_production_pxe.vmlinuz.sig"

# PXE initrd and sig
echo "flatcar_production_pxe_image.cpio.gz"
curl -# "${BASE_URL}/flatcar_production_pxe_image.cpio.gz" -o "${DEST}/flatcar_production_pxe_image.cpio.gz"
echo "flatcar_production_pxe_image.cpio.gz.sig"
curl -# "${BASE_URL}/flatcar_production_pxe_image.cpio.gz.sig" -o "${DEST}/flatcar_production_pxe_image.cpio.gz.sig"

# Install image
echo "flatcar_production_image.bin.bz2"
curl -# "${BASE_URL}/flatcar_production_image.bin.bz2" -o "${DEST}/flatcar_production_image.bin.bz2"
echo "flatcar_production_image.bin.bz2.sig"
curl -# "${BASE_URL}/flatcar_production_image.bin.bz2.sig" -o "${DEST}/flatcar_production_image.bin.bz2.sig"

# Install oem image
if [[ -n "${IMAGE_NAME-}" ]]; then
  echo "${IMAGE_NAME}"
  curl -# "${BASE_URL}/${IMAGE_NAME}" -o "${DEST}/${IMAGE_NAME}"
  echo "${IMAGE_NAME}.sig"
  curl -# "${BASE_URL}/${IMAGE_NAME}.sig" -o "${DEST}/${IMAGE_NAME}.sig"
fi

# verify signatures
$GPG --verify "${DEST}/flatcar_production_pxe.vmlinuz.sig"
$GPG --verify "${DEST}/flatcar_production_pxe_image.cpio.gz.sig"
$GPG --verify "${DEST}/flatcar_production_image.bin.bz2.sig"

# verify oem signature
if [[ -n "${IMAGE_NAME-}" ]]; then
  $GPG --verify "${DEST}/${IMAGE_NAME}.sig"
fi
