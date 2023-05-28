#!/bin/bash
cd "$(dirname "$0")"
cfssl gencert \
  -ca=ca.pem \
  -ca-key=ca-key.pem \
  -config=ca-config.json \
  -profile=client \
  client-csr.json | cfssljson -bare client
