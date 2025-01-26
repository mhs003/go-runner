#!/bin/bash
APP="run"
VERSION="1.0.0"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    # "windows/amd64"
    # "windows/386"
    "darwin/amd64"
    "darwin/arm64"
)

mkdir -p bin

for platform in "${PLATFORMS[@]}"; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output="bin/${GOOS}-${GOARCH}/${APP}"
    
    # if [ "$GOOS" = "windows" ]; then
    #     output+=".exe"
    # fi

    chmod +x "$output"

    echo "Building ${output}"
    GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags "-s -w" -o "$output"
done