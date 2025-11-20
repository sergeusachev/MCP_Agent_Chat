#!/bin/bash

curl -i -k -L POST 'https://gigachat.devices.sberbank.ru/api/v1/chat/completions' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <REQUEST_TOKEN>' \
  -d @request_functions.json