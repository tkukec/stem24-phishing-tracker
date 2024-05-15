# Agent management

Service exposes an esl, rest and event bus option to control data flow from and to Freeswitch

## Install

``` 
  cp .env.example .env
  go run cmd/stem24-backend/main.go
```
## Skill group discovery
To pull skill groups and skill group members from other services
Use skill-group discovery command:


```go run .\cmd\tools\main.go skill-groups:discover -l email@email:7072,telephone@telephone:8080 -t ... -z ...```

Parameters:
``` 
-l -> service locations eg. telephone@telephone:8080,email@email:7072
-t -> tentant id to use (if not defined script will use SINGLE TENANT MODE)
Optional params: (If not defined script will use defaults from ENV) 
-i -> iam location eg. 10.135.11.103:7072
-r -> iam realm eg. evil
-z -> client secret
-c -> client name
--db-port -> database port eg. 5432
--db-user -> database user eg. live
--db-pass -> database password eg. live
--db-host -> database host address eg. 10.20.30.40
--amqp-host -> broker host address eg. 10.20.30.40
--amqp-port -> broker port eg. 6161
--amqp-user -> broker user eg. live
--amqp-pass -> broker password eg. live
```

## Dependency
Note: All configuration for non-production images should be done through an .env file (view .env.example)

### Running database
Currently, only supports mysql, mariadb, postgres

### ENV configuration
```
HTTP_PORT=8080

DB_DRIVER=mysql
DB_HOST=stem24-backend-db
DB_PORT=3306
DB_USER=root
DB_NAME=stem24-backend
DB_PASS=secret
DB_DEBUG=false
DB_SEED=true

#Possible values: DEBUG,WARN,INFO,ERROR (values are case insensetive). Defaults to ERROR
LOG_LEVEL=DEBUG

SINGLE_TENANT_MODE=false

# Set max number of connections in the idle connection pool. If negative, no idle connections are retained. 0 = DEFAULT.
MAX_IDLE_CONNS=0
# Set max number of open connections to the database. If negative then there is no limit.
MAX_OPEN_CONNS=0
# Set max amount of time a connection may be reused (in hours !!). If negative  connections are not closed due to it's age
CONN_MAX_LIFETIME=0

# Graylog service data
GRAYLOG_PORT=12201
GRAYLOG_HOSTNAME=

# should we log api endpoints
REQUEST_LOG=false

# should we log db queries
QUERY_LOG=false

# should response body log be added to a log message
RESPONSE_BODY_LOG=false

# log drivers which should be used
LOG_DRIVERS=graylog,stdout,file
```

####Logging
```
#Log for rest api request and others SEPERATE THIS LATER
logs/stem24-backend.log
```

####Open api

Open api .yml and .json files can be found under api/openapi/ directory

NOTE: when updating openapi docs, run following ```swag init -g cmd/stem24-backend/main.go -o api/openapi --parseVendor --pd --parseDepth 1 -ot yaml```. The command will read all router comments and make a swagger file

#####Tests

```go test ./... -v | go-junit-report > report.xml```

NOTE: when creating mocks run the following command  ```mockery --all --keeptree```