#!/usr/bin/env bash
curl --request POST 'http://localhost:8080/route' \
--header 'Content-Type: application/json' \
--data-raw '{
        "from": "zzz",
        "to": "ttt",
        "cost": 18.1007
}'
