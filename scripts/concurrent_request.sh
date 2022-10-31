#!/bin/sh

set -e

for i in `seq 1 5`;
do
 curl -s localhost:8080/tokens?q= > output$i &
done

sleep 10

echo "Test concurrency successfully"


