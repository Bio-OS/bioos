#!/bin/bash
cd "$(dirname "$0")"

if ! [[ -f ca.pem ]];then
  ./ca.sh
fi

if ! [[ -f server.pem ]];then
  ./server.sh
fi

if ! [[ -f client.pem ]];then
  ./client.sh
fi