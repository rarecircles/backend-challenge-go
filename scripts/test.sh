#!/bin/sh

set -e

./build/challenge -rpc-url=https://eth-mainnet-public.unifra.io &

sleep 20

curl localhost:8080/tokens?q=ThereIsNoToken > output0

diff output0 ./test/output/ThereIsNoToken

pkill challenge

rm -rf output0

echo "Test integrity successfully"


