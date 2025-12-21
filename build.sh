#!/bin/bash

# 创建 bin 目录（如果不存在）
mkdir -p bin

platforms=(
  "windows/amd64"
  "windows/386"
  "windows/arm64"
  "linux/amd64"
  "linux/386"
  "linux/arm64"
  "linux/arm"
  "darwin/amd64"
  "darwin/arm64"
)

for platform in "${platforms[@]}"
do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name="shellgo_${GOOS}_${GOARCH}"
  
  if [ $GOOS = "windows" ]; then
    output_name+='.exe'
  fi
  
  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o bin/$output_name .
  if [ $? -ne 0 ]; then
    echo "An error occurred. Aborting."
    exit 1
  fi
done

echo "Build complete! All binaries are in the bin/ directory."