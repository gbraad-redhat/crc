#!/bin/sh

# clean workspace folder
rm -rf /workspace/crc
mkdir /workspace/crc
ln -s /workspace/crc ~/Projects
git init /workspace/crc

cd ~/Projects

exit 0
