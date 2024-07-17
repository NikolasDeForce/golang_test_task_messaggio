# CURL Commands to test how to works website

or - `docker-compose up`

# Before start check, you should start server PostgreSQL, change ConnectSettings in file db.go

check notSupported method - `curl localhost:8010/`

check notSupport method  - `curl localhost:8010/doesNotExist`

check GET method /getall without options - `curl -s -X PUT -H 'Content-Type: application/json' localhost:8010/getall`

check GET method /getall with options - `curl -s -X GET -H 'Content-Type: application/json' -d '{"text":"hello world!"}' localhost:8010/getall`

check GET method /message/id - `curl -X GET -H 'Content-Type: application/json' -d '{"text":"hello world!"}' localhost:8010/message/1`

check POST method /add - `curl -X POST -H 'Content-Type: application/json' -d '[{"text":"hello world!"}, {"text": "hello everybody"} ]' localhost:8010/add`

check DELETE method /message/id - `curl -X DELETE -H 'Content-Type: application/json' -d '{"text": "hello world!"}' localhost:8010/message/3 -v`