#!/usr/bin/env bash

mkdir -p dist/
rm -rf dist/*

package="github.com/scaleway/c14-cli/cmd/c14"
package_name="c14"

platforms=(
        "windows/amd64"
        "windows/386"
        "darwin/amd64"
        "darwin/386"
        "linux/amd64"
        "linux/386"
        "linux/arm"
        "linux/arm64"
        "freebsd/amd64"
        "freebsd/386"
        "freebsd/arm"
        )

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $platform... (dist/$output_name)"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o dist/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

echo "Done."
