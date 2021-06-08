#!/bin/bash
# Build generic golang application with provided version.
#
# Example:
#	build-go cmd/server example-server pkg.version.Version 1.2.0
set -e

if [ "$#" -ne 4 ]; then
	echo "Usage: $0 <package> <target_binary> <version_var> <version>"
	exit 1
fi

readonly package="$1"
readonly target_binary="$2"
readonly version_var="$3"
readonly version="$4"

if [ -z "$5" ];then
readonly bin_dir="bin/"
else
readonly bin_dir="$CI_PROJECT_DIR/bin/"
fi

echo "Building $target_binary at version $version"

mkdir -p "$bin_dir"
CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -ldflags="-X $version_var=$version" -o "$bin_dir/$target_binary" "$package"
