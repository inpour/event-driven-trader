COMPOSE := docker compose -f deployments/docker-compose.yml
KAFKA_SERVICE := kafka
KAFKA_BOOTSTRAP_SERVER := kafka:9092
KAFKA_BIN := /opt/kafka/bin

TOPIC_CANDLE_CLOSED := market.candle.closed
TOPIC_INPUT_COMPLETED := backtest.input.completed

.PHONY: \
	rebuild \
	rebuild-no-cache \
	up \
	down \
	stop \
	start \
	restart \
	ps \
	logs \
	logs-kafka \
	logs-topic-init \
	logs-market-data \
	build-market-data \
	rebuild-market-data \
	reset \
	kafka-topics \
	kafka-delete-topics \
    kafka-create-topics \
    kafka-empty-topics \
	kafka-describe-candle-closed-topic \
	kafka-describe-input-completed-tipic \
	kafka-candle-closed-events \
	kafka-input-completed-events

rebuild:
	$(COMPOSE) up --build --force-recreate -d

rebuild-no-cache:
	$(COMPOSE) build --no-cache
	$(COMPOSE) up --force-recreate -d

reset:
	$(COMPOSE) down
	$(COMPOSE) up -d

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

stop:
	$(COMPOSE) stop

start:
	$(COMPOSE) start

restart:
	$(COMPOSE) stop
	$(COMPOSE) start

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f

logs-kafka:
	$(COMPOSE) logs -f kafka

logs-topic-init:
	$(COMPOSE) logs topic-init

logs-market-data:
	$(COMPOSE) logs market-data-service

build-market-data:
	$(COMPOSE) build market-data-service

rebuild-market-data:
	$(COMPOSE) build market-data-service
	$(COMPOSE) rm -f market-data-service
	$(COMPOSE) up market-data-service

kafka-topics:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--list

kafka-create-topics:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--create \
		--if-not-exists \
		--topic $(TOPIC_CANDLE_CLOSED) \
		--partitions 1 \
		--replication-factor 1

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--create \
		--if-not-exists \
		--topic $(TOPIC_INPUT_COMPLETED) \
		--partitions 1 \
		--replication-factor 1

kafka-delete-topics:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--delete \
		--if-exists \
		--topic $(TOPIC_CANDLE_CLOSED)

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--delete \
		--if-exists \
		--topic $(TOPIC_INPUT_COMPLETED)

kafka-empty-topics:
	$(MAKE) kafka-delete-topics
	@echo "Waiting for Kafka to delete project topics..."
	# because Kafka topic deletion is often asynchronous
	@until ! $(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--list | grep -Eq '^($(TOPIC_CANDLE_CLOSED)|($(TOPIC_INPUT_COMPLETED))$$'; do \
		sleep 1; \
	done
	$(MAKE) kafka-create-topics
	$(MAKE) kafka-topics

kafka-describe-candle-closed-topic:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--topic $(TOPIC_CANDLE_CLOSED)

kafka-describe-input-completed-topic:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--topic $(TOPIC_INPUT_COMPLETED)

kafka-candle-closed-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_CANDLE_CLOSED) \
		--from-beginning \
  		--max-messages 100

kafka-input-completed-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_INPUT_COMPLETED) \
		--from-beginning \
  		--max-messages 100