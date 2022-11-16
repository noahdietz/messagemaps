#!/bin/bash

out=all.txt # resources_only.txt
flags=out_file=$out #,resources_only=true

rm -rf $out

go install ./protoc-gen-messagemaps

if [ -z $GOOGLEAPIS ]; then
  echo "Cloning googleapis to current directory"
  git clone --depth=1 https://github.com/googleapis/googleapis.git
  echo "export GOOGLEAPIS=$(pwd)/googleapis - DO THIS IN YOUR SHELL TO RETAIN"
  export GOOGLEAPIS=$(pwd)/googleapis
fi

pushd $GOOGLEAPIS

find google/cloud -name '*.proto' -type f -exec dirname {} \; | uniq | while read line; do
    protoc -I. --messagemaps_out=$flags:. $line/*.proto
done

popd

mv $GOOGLEAPIS/$out .
