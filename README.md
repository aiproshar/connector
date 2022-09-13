
# Connector

A standalone bridge application between MQTT and Kafka.


## Application Details
MQTT is a the preferred go to protocol for IoT devices, due to its low protocol overhead. it was not built for data integration and data processing. In these scenaros, developers prefer kafka or similar event streaming platform to achieve event driven architecture to process data in motion.

To clarify, MQTT and Kafka complement each other. So how we pass our huge amount of events from MQTT to kafka easily ? This is where **connector** comes to help.

**Connector** is a highly scalable and fault tolerant Application which continiously pulls events from MQTT brokers and publishes those events to Air-Gapped kafka clusters with minimal latency. And depending on configuration, it can achieve fault tolerence, so when the application crashes or kafka cluster is offline, the messages are stored and published as soon as the cluster comes online.

## Deployment

To deploy this project run, we need

- MOTT Instance
- ZooKeeper Instance
- Kafka Instance

Lets deploy the MQTT broker

```bash
    docker run --rm --net=host --name mosquitto -p 1883:1883  -p 9001:9001 eclipse-mosquitto mosquitto -c /mosquitto-no-auth.conf
```
ZooKeeper 

```bash
    docker run --rm --net=host --name zookeeper -p 2181:2181 zookeeper
```
Find the zookeeper ip using docker command

```bash
    docker inspect zookeeper
```
Deploy the kafka cluster, update envar KAFKA_ZOOKEEPER_CONNECT as necessary

```bash
    docker run --rm --net=host --name kafka -p 9092:9092 -e KAFKA_ZOOKEEPER_CONNECT=172.17.0.1 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 confluentinc/cp-kafka
```
Deploy the Connector Application after setting the required environment variables (described later)

```bash
docker run --rm --net=host --name connector \
-e RUN_MODE=PRODUCTION \
-e MQTT_BROKER_URL=127.0.0.1 \
-e MQTT_BROKER_PORT=1883 \
-e MQTT_CLIENT_ID=arafat \
-e MQTT_USER_NAME=user \
-e MQTT_USER_PASSWORD=12345678 \
-e LOG_LVL=info \
-e TOPIC_LIST=default:default-kafka,custom:custom-kafka,dev_arafat:dev_arafat \
-e WILDCARD_TOPIC_LIST="default/+/mqtt:custom-kafka,default/kc/#:general" \
-e INVOKE_LATENCY_SEC=10 \
-e INVOKE_MAX_GOROUTINE=1000 \
-e INVOKE_MIN_GOROUTINE=50 \
-e KAFKA_BROKER=localhost:9092 \
-e SASL_SSL=false \
arafatkhan/connector
```
## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

- `RUN_MODE` can `DEVELOP` or `PRODUCTION`. if this variable is not set then it will default to DEVELOP
and the required environment variables will be loaded from .env file.

- `MQTT_BROKER_URL` URL of the MQTT Broker 


- `MQTT_BROKER_PORT` Specify Broker Port


- `MQTT_CLIENT_ID` Specify the identitiy of the client

- `MQTT_USER_NAME` Set Username if auth-conf is added


- `MQTT_USER_PASSWORD` Set Password if auth-conf is added

- `LOG_LVL` Connector Uses Zerolog as logging package, it has 7 logging levels. Levels: - panic fatal error warn info debug trace


- `TOPIC_LIST` Comma separated topic mapping from MQTT to Kafka Cluster, for Example

`TOPIC_LIST=default:default-kafka,custom:custom-kafka`
maps IoT events from `default` topic to `default-kafka` topic in kafka and so on.

- `WILDCARD_TOPIC_LIST` MQTT supports wildcard topic, allowing thousand of individual topic to map into a 

single wildcard. Its is possible that a topic may end into multiple wildcards, in that scenario,
connector will push it only into the first mapped kafka topic and will throw waring log.
Eg.  `WILDCARD_TOPIC_LIST="default/+/cars:cars,default/general/#:general"`

- `INVOKE_LATENCY_SEC` Specify after how much time DbInvoke Should be called. This will try to publish all the messages stored in the database, which was previously failed to publish.

- `INVOKE_MAX_GOROUTINE` no of maximum concorrent go routines to publish stored messages
- `INVOKE_MIN_GOROUTINE` number of go routines upon which DbInvoke will be called to spawn new goroutines
- `KAFKA_BROKER` kafkka broker with port, Eg. `KAFKA_BROKER=localhost:9092`
- `SASL_SSL` connector supports PLAINTEXT and SAS/OAUTHBEARER protocols. If set to false it will use PLAINTEXT (default)
- `SASL_OAUTHBEARER_CLIENT_ID` if SASL_SSL=true then set Oauth CLient ID
- `SASL_OAUTHBEARER_CLIENT_SECRET` Oauth Clinet Secret
- `SASL_OAUTHBEARER_TOKEN_ENDPOINT_URI` Oauth token Endpoint Providers URI

Eg.  `SASL_OAUTHBEARER_TOKEN_ENDPOINT_URI=https://yourdomain.com/auth/realms/myrealm/protocol/openid-connect/token`
- `DB_PERSIST_PATH` If we want to persist our data on Sqlite, if this path is not found then switches to in-memory mode.
## Run Locally with Clients 

Clone the project

```bash
  git repo clone aiproshar/connector
```

Go to the project directory  

```bash
  cd connector
```

Install dependencies

```bash
  go mod tidy
```

Start the application

```bash
  go run application.go
```
Note: Connector Requires 

### Run the kafka subscriber

Go to the subscriber directory

```bash
  cd cd clients/subscriber-kafka
```
Configure the .env files as necessary (topic list)

Install dependencies

```bash
  go mod tidy
```

Start the subscriber

```bash
  go run application.go
```

### Run the mqtt publisher

Go to the publisher directory

```bash
  cd clients/publisher-mqtt
```
Configure the .env files as necessary (topic list, publish_interval, total msg etc)

Install dependencies

```bash
  go mod tidy
```

Start the subscriber

```bash
  go run application.go
```
## Architecture

Connector uses 4 channels to pass incomming events to the publish handler, 

or database insert handler after failure, database channel to restore messages failed 

previously and database delete channel to delete successfully published events which were failed previously.

Let's have a look at the whole design diagram.

![Alt text](images/architecture.png?raw=true "ARCHITECTURE")



## License

[MIT](https://choosealicense.com/licenses/mit/)

