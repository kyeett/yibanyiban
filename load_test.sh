#!/bin/bash -e

USERS=2000

if command -v fasthttploader > /dev/null ; then
    # OK!
    echo "Running load test with $USERS of concurrent clients"
else
    echo """Command 'fasthttploader' missing! Please install:

    go get github.com/hagen1778/fasthttploader
"""
    exit 1
fi

for i in 1000 5000 10000 15000 20000 25000 30000
do
    echo -e "\nTrying to reach $i request per second"
    fasthttploader -c $USERS -q $i http://localhost:8080/validate"?iban=GB82WEST12345698765432" | grep -E "QPS:|Req done:|QPS:|Errors:" | cut -f1 -d ";"
done