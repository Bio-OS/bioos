#!/bin/bash
cd "$(dirname "$0")"
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
