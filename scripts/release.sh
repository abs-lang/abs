#!/usr/bin/env bash

package="main.go"
package_split=(${package//\// })
package_name=${package_split[-1]}

# Platforms commented out seem
# to be a PITA to compile without
# too much tooling, we'll have a 
# look at those later on.
platforms=(
    "linux/386"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "linux/ppc64"
    "linux/ppc64le"
    # "linux/mips"
    #"linux/mipsle"
    # "linux/mips64"
    # "linux/mips64le"
    "linux/arm"
    "linux/arm64"
    "windows/amd64"
    "windows/386"
    "darwin/amd64"
    # "darwin/arm"
    "darwin/386"
    "freebsd/386"
    "freebsd/amd64"
    "freebsd/arm"
    "netbsd/386"
    "netbsd/amd64"
    "netbsd/arm"
    # "plan9/386"
    # "plan9/amd64"
    # "openbsd/386"
    # "openbsd/amd64"
    # "openbsd/arm"
    # "dragonfly/amd64"
    # "android/arm"
    # "solaris/amd64"
)

echo "Deleting previous builds..."
rm -rf builds/abs-*

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='abs-preview-1-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    echo "Building $GOOS-$GOARCH"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o builds/$output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done