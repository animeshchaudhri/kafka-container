# Kafka Docker Setup

This repository contains the configuration and instructions to set up and use Apache Kafka with Docker.

## Prerequisites

- Docker installed on your machine.
- Docker Compose installed on your machine.

## Services

This setup includes the following services:
- **Zookeeper**: Confluent Zookeeper image to coordinate and manage Kafka brokers.
- **Kafka**: Confluent Kafka image to handle message brokering.

## Setup Instructions

1. **Clone the Repository**

   ```sh
   git clone https://github.com/animeshchaudhri/kafka-container.git
   cd kafka-docker-setup



### Start Services

To start the Kafka and Zookeeper services, run the following command in the directory containing your `docker-compose.yml` file:

```sh
docker-compose up -d
```

## Locate Kafka Console Producer

To find the location of the Kafka console producer script, you can list the directories inside the Kafka container:

```sh
docker exec -it <kafka-container-id> ls /opt
docker exec -it <kafka-container-id> ls /opt/confluent
docker exec -it <kafka-container-id> ls /opt/confluent/bin

```
Replace <kafka-container-id> with the actual container ID of your Kafka service.

### Produce a Message

Once you've located the `kafka-console-producer.sh` script (typically found in `/opt/confluent/bin`), you can use it to send messages to a Kafka topic.

#### Interactive Mode

Run the Kafka console producer:

```sh
docker exec -it <kafka-container-id> /opt/confluent/bin/kafka-console-producer.sh --broker-list localhost:29092 --topic my-topic

```