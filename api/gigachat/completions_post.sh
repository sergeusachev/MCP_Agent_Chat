#!/bin/bash

curl -i -k -L POST 'https://gigachat.devices.sberbank.ru/api/v1/chat/completions' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer <REQUEST_TOKEN>' \
  -d '{
    "model": "GigaChat-2",
    "messages": [
      {
        "role": "user",
        "content": "Что написано в файле readme?"
      }
    ],
    "temperature": 0,
    "repetition_penalty": 1
  }'