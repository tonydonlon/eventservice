#!/usr/bin/env bash

# k6 test to run; defaults to full suite run

export EVENT_URL=http://localhost:8080
export SESSION_ID=

SCRIPT=
if [ -z "$1" ]
    then
        echo "Running full suite of k6 tests: suite.js"
        SCRIPT="./suite.js"
else
    SCRIPT="$1"
    echo "Running k6 test: $SCRIPT"
fi

k6 run $SCRIPT
