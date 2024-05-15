# Agent management

Service exposes an esl, rest and event bus option to control data flow from and to Freeswitch

## Install

``` 
  cp .env.example .env
  go run cmd/agent-management/main.go
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

To pull skill groups and members on start of Agent-Management for every tenant, define locations in env:
```
#services to pull skill groups and members from seperated by "," eg. telephone@telephone:7280,email@10.135.150.90:8080
SG_DISCOVERY_LOCATIONS=telephone@telephone:8080,email@email:7072
```

## Dependency
Note: All configuration for non-production images should be done through an .env file (view .env.example)

### Running database
Currently, only supports mysql, mariadb, postgres

### Running broker
Currently, only supports amqp 1.0 event brokers

### Sip Service
In order to use the telephone|videochat capabilities a running instance of freeswitch is required: [More here](https://git.asseco-see.hr/asseco-hr-voice/evil/sipservice)

### ENV configuration
```
HTTP_PORT=8080

DB_DRIVER=mysql
DB_HOST=agent-management-db
DB_PORT=3306
DB_USER=root
DB_NAME=agent-management
DB_PASS=secret
DB_DEBUG=false
DB_SEED=true

BROKER_TYPE=amqp
BROKER_HOST=broker
BROKER_PORT=61616
BROKER_USER=admin
BROKER_PASS=admin
BROKER_TOPICS=sip_service_notifications,telephone,iam_notifications,email,sms,social-network

IAM_URI=http://live-iam:8080
IAM_REALM=evil
PUSHER_URI=http://pusher:8080
CONTACTS_URI=http://contacts:8080
SIP_SERVICE_URI=http://sip-service:8080

CLIENT_ID=agent-management
CLIENT_SECRET=
#DEFAULTS TO THIS
TELEPHONE_CHANNEL=telephone
VIDEO_CHAT_CHANNEL=video_chat

#Possible values: DEBUG,WARN,INFO,ERROR (values are case insensetive). Defaults to ERROR
LOG_LEVEL=DEBUG

#services to pull skill groups and members from seperated by "," eg. telephone@telephone:8080,email@email:8080
SG_DISCOVERY_LOCATIONS=

GRACE_LOGIN_PERIOD=60
#if set to true, tenant is not required and the base seeded tenant "*" will be used
SINGLE_TENANT_MODE=false
#true by default, persist to database agent ponder data when matching for user, "./tools queue:explain" looks at this data 
PERSIST_PONDER_SNAPSHOTS=true

#Time span to take into account when calculating agnet ponder values
IDLE_PONDER_TIME_SPAN=1440
ACTIVITIES_PONDER_TIME_SPAN=1440

# Set max number of connections in the idle connection pool. If negative, no idle connections are retained. 0 = DEFAULT.
MAX_IDLE_CONNS=0
# Set max number of open connections to the database. If negative then there is no limit.
MAX_OPEN_CONNS=0
# Set max amount of time a connection may be reused (in hours !!). If negative  connections are not closed due to it's age
CONN_MAX_LIFETIME=0

# Enable CRON jobs such as auto skill-group sync by using CRON Schedule experssions e.g. * * 10 * * or @daily or @every 10h5m30s ...
# USEFUL TOOL: https://crontab.guru/
SKILL_GROUP_CRON_SCHEDULE=

# Graylog service data
GRAYLOG_PORT=12201
GRAYLOG_HOSTNAME=

# should we log api endpoints
REQUEST_LOG=false

# should we log broker events
QUEUE_LOG=false

# should we log db queries
QUERY_LOG=false

# should response body log be added to a log message
RESPONSE_BODY_LOG=false

# log drivers which should be used
LOG_DRIVERS=graylog,stdout,file

#How many workers for each system channel to spin up (match channel, queued channel, ...)
#Defaults to 3, WARNING: high numbers could impact CPU and RAM usage
WORKERS_PER_SYSTEM_CHANNEL=3

# pprof port to be used
PPROF_ENABLED=false
PPROF_PORT=6061

# INT TIME IN SECONDS HOW LONG WILL AM WAITS FOR TELEPHONE/VIDEO/SN SERVICE RESPONSE
# DEFAULT 4000 miliseconds
QUEUE_RESOLVE_SLEEP_TIME=4000
```

####Logging
There are three main logs.
```
#Log for rest api request
logs/agent-management.log
#Log for socket connections
logs/agent-management-io.log
#Log for broker events
logs/broker_log.log
```

####Open api

Open api .yml and .json files can be found under api/openapi/ directory

NOTE: when updating openapi docs, run following ```swag init -g cmd/agent-management/main.go -o api/openapi --parseVendor --pd --parseDepth 1 -ot yaml```. The command will read all router comments and make a swagger file

#####Tests

```go test ./... -v | go-junit-report > report.xml```

NOTE: when creating mocks run the following command  ```mockery --all --keeptree```