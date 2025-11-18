#!/bin/bash

curl --request GET \
  --url 'https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=bitcoin' \
  --header 'x-cg-demo-api-key: CG-hsFyp7dn5iXLD1moJEsgPbPy'

# Response for this request: {"bitcoin":{"usd":92980}}