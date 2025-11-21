#!/bin/bash

#RqUID can be any value, it doesn matter

curl -i -k -L 'https://ngw.devices.sberbank.ru:9443/api/v2/oauth' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer MDE5YTQ2NzEtNjkwOC03ZDQ3LTljZDktYWU3ZTYzYzcyYzA3OmNkNzE5NjA4LWJmMDQtNDZlOC1hNTFjLTM5OTlmYWZhY2Q4OQ==' \
  -H 'RqUID: 270fee8f-3594-4cb7-b9cb-d0690691f735' \
  -d 'scope=GIGACHAT_API_PERS'