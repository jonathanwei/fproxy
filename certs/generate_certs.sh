#!/bin/bash

set -e

NFRONTENDS=${1:-1}
NBACKENDS=${2:-1}

CONFIGDIR=`pwd`
WORKDIR=`mktemp -d --tmpdir="$CONFIGDIR"`

echo "Generating certs in $WORKDIR for $NFRONTENDS frontend(s) and $NBACKENDS backend(s)."
cd "$WORKDIR"

cfssl gencert -initca "$CONFIGDIR/frontend_ca.json" | cfssljson -bare frontend-ca
cfssl gencert -initca "$CONFIGDIR/backend_ca.json" | cfssljson -bare backend-ca

for i in `seq 1 $NFRONTENDS`; do
  cfssl gencert -ca frontend-ca.pem -ca-key frontend-ca-key.pem "$CONFIGDIR/frontend.json" \
    | cfssljson -bare "frontend-$i"
done


for i in `seq 1 $NBACKENDS`; do
  cfssl gencert -ca backend-ca.pem -ca-key backend-ca-key.pem "$CONFIGDIR/backend.json" \
    | cfssljson -bare "backend-$i"
done
