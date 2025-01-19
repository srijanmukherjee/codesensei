#!/bin/sh

set -xe

DEPENDENCIES_DIR=deps
ISOLATE_REVISION=ad39cc4d0fbb577fb545910095c9da5ef8fc9a1a

mkdir -p $DEPENDENCIES_DIR

# Download and build isolate
if [[ ! -d $DEPENDENCIES_DIR/isolate ]]; then
    git clone https://github.com/judge0/isolate.git $DEPENDENCIES_DIR/isolate
fi
(cd $DEPENDENCIES_DIR/isolate; git checkout $ISOLATE_REVISION)
(cd $DEPENDENCIES_DIR/isolate; sudo make -j$(nproc) install)
