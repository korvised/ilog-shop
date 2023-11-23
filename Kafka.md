# Kafka

## Install Kafka
```shell
docker-compose -f docker-compose.kafka.yml up -d
```

## This project we use topic to send message to kafka

1. Create topic

```shell
# run in kafka-1
docker exec -it kafka-1 bash
# go to kafka bin
cd /opt/bitnami/kafka/bin
# create inventory topic
kafka-topics.sh --create --topic inventory --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
# create payment topic
kafka-topics.sh --create --topic payment --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
# create player topic
kafka-topics.sh --create --topic player --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
# show topic list
kafka-topics.sh --list --bootstrap-server localhost:9092
# describe topic
kafka-topics.sh --describe --topic inventory --bootstrap-server localhost:9092
```