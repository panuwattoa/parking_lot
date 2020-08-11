#!/usr/bin/env bash
package_split=(${package//\// })
package_name="parkinglot"

platforms=("linux/amd64")
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    platforms=("linux/amd64")
elif [[ "$OSTYPE" == "darwin"* ]]; then
    platforms=("darwin/amd64")
fi
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    echo '[Build Success]'
done
cd bin
./run_functional_tests