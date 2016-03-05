#!/usr/bin/env bash
# Generate Go protobuf code.

set -e

COREOS_ROOT="$GOPATH/src/"

# protobuf subpackages end in "pb"
PBUFS=$(go list ./... | grep -v /vendor | grep 'pb$')

# change into each protobuf directory
for pkg in $PBUFS ; do
  abs_path=${GOPATH}/src/${pkg}
  echo Generating $abs_path
  pushd ${abs_path} > /dev/null
  # generate protocol buffers, make other .proto files available to import
  protoc --go_out=plugins=grpc:. -I=.:"${COREOS_ROOT}" *.proto
  popd > /dev/null
done