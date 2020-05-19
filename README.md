# Test Log Server

## Append Log
curl -X POST localhost:8080 -d \
'{"record": {"value": "TGV0J3MgR28gIzEK"}}'

## Read Log
curl -X GET localhost:8080 -d \
'{"offset": 0}'

## Unsupported
curl -X PUT localhost:8080