#!/bin/bash

curl -i -k -L POST 'https://gigachat.devices.sberbank.ru/api/v1/chat/completions' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: Bearer eyJjdHkiOiJqd3QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwiYWxnIjoiUlNBLU9BRVAtMjU2In0.qXJ2H5UyD4RE4toLqkJDybkZ1423XiG_1Jm8Uk3IA8-VSl1eKfUKehk1nYM_ZeK4mD52RH1ZkMp5kk1UMCoDfgbZY12uM99t1XI-wBwqDi_mVIhNqvUguA3UROCzc6HeKfOchH_4wQvrTUEviI-vRfsGfgITtFrg4S2D9acIfZNOEOTgIVOEsBaFca4jnXkI65Q1IGaipb9ijSL5r-ttnSIxJAD9Ti9Tlfneeop9DYCgefhT9bLqEBPZHHsxZpclL6bdfWTHc0iaGqWGytegajTIpv9WqhF8EYGjaervl5eq1xb8xDsnAAK61A1se79Wh4dj-yE6YBAV2VxiKZIp2w.zZG7fLExFN9TFRMNRyGVYQ.4pn85AaLJ51ZV-GH2amk6JrApXBQtHXypYc0QDzb-nbX_I-99k4prYwwPTExp5zo5lC-vmH_zzc9ljtep1ltZn90TcZuL1iu7h4ZHgkkIBYwvX5wxOPI3M5zbLiaciFAk5ZnJP6KcJV4FT4qFnFxT9wOLMj8COj0RTl09DQuic8f0Vwo2sexrZrTTmsZL_GI6vo4syONEG6Cj7Y2JzLxONfxZhvs1GGOTu2j43Z5pmLIG3xLFV0KGacfZJliRIQpZg4vWk2DR0qlfmqhgCnmVW9jllWZ1GIOGb0206pCuEljB6haPOdvjVfv7V945Hzi5D6vLDYQsDuSFi8LVIKk0Wfx5MYwlCqSIUjAPX8nbMnoRv4iyzLuR6kEBTzwQ0lXqP0C4WLW8IrCM2iKLy8ZHmBR-EZVol4WW7ByJq9MnN31ELI5fEt6aFIQjoqpuIS12GWxfn9sustCe6EhpYyHbbOvxVroSZrIA8mWTb_M2GRVijZNy7rU0YW79sAGCrMKDF9XdHkyQjEUnSc_uWt9skMf9WccBeBjZNsvH5b6vcbQlG3ySgSlgDQDIqwcSjhRgD7Rtsc-uxwmVYckN8P-RwJVcWVSMYwi5P00WeUJIfjZ3ResMjSuRDAAiDl3hK_IXtPF6dVMcWtT_us6Dvl9_dtNdhlROmO4hlC6IcySDo7KUoJo-_xhCXhgM9WvCrFIkb0aiFuiVIKS3X6z6KZjBoZf6zySHhdnXETbTQlrR40.C3nnb3CiUyD__Lo3epKIGx-20moS2eNJUpiXD1MaPDc' \
  -d '{
    "model": "GigaChat-2",
    "messages": [
      {
        "role": "user",
        "content": "How much is bitcoin now?"
      }
    ],
    "function_call": "auto",
    "functions": [
      {
        "name": "get_crypto_price",
        "description": "Getter for crypto price",
        "parameters": {
          //inject InputScheme content here
        }
      }
    ],
    "temperature": 0,
    "repetition_penalty": 1
  }'