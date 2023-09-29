#!/bin/sh
set -e
num=${NUM_PHYLAXS:-1} # default value for NUM_PHYLAXS = 1
for ((i=0; i<num; i++)); do
    while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' phylax-$i.phylax:6060/readyz)" != "200" ]]; do sleep 5; done
done
CI=true npm run test-accountant
