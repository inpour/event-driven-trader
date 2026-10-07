COMPOSE := docker compose -f deployments/docker-compose.yml
KAFKA_SERVICE := kafka
KAFKA_BOOTSTRAP_SERVER := kafka:9092
KAFKA_BIN := /opt/kafka/bin

TOPIC_CANDLE_CLOSED := market.candle.closed
TOPIC_INPUT_COMPLETED := backtest.input.completed
TOPIC_MARKER_CREATED := backtest.marker.created
TOPIC_RUN_COMPLETED := backtest.run.completed

MARKET_DATA_SERVICE := market-data-service
ORDER_EXECUTION_SERVICE := order-execution-service
BACKTEST_SERVICE := backtest-service
CHART_SERVICE := chart-service

BACKTEST_GROUP_ID := backtest-service

.PHONY: \
	rebuild \
	rebuild-no-cache \
	reset \
	up \
	down \
	stop \
	start \
	restart \
	ps \
	logs \
	logs-kafka \
	logs-topic-init \
	build-market-data \
	logs-market-data \
    run-market-data \
	build-order-execution \
    logs-order-execution \
    run-order-execution \
    build-backtest \
    logs-backtest \
    run-backtest \
	kafka-topics \
	kafka-delete-topics \
    kafka-create-topics \
    kafka-empty-topics \
	kafka-describe-candle-closed-topic \
	kafka-describe-input-completed-topic \
	kafka-describe-marker-created-topic \
	kafka-describe-run-completed-topic \
	kafka-candle-closed-events \
	kafka-input-completed-events\
	kafka-marker-created-events \
	kafka-run-completed-events\
	kafka-consumer-groups \
	kafka-backtest-group

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

build-market-data:
	$(COMPOSE) build $(MARKET_DATA_SERVICE)

logs-market-data:
	$(COMPOSE) logs $(MARKET_DATA_SERVICE)

run-market-data:
	$(COMPOSE) rm -f $(MARKET_DATA_SERVICE)
	$(COMPOSE) up $(MARKET_DATA_SERVICE)

build-order-execution:
	$(COMPOSE) build $(ORDER_EXECUTION_SERVICE)

logs-order-execution:
	$(COMPOSE) logs -f $(ORDER_EXECUTION_SERVICE)

run-order-execution:
	$(COMPOSE) up -d $(ORDER_EXECUTION_SERVICE)

build-backtest:
	$(COMPOSE) build $(BACKTEST_SERVICE)

logs-backtest:
	$(COMPOSE) logs -f $(BACKTEST_SERVICE)

run-backtest:
	$(COMPOSE) up -d $(BACKTEST_SERVICE)

build-chart:
	$(COMPOSE) build $(CHART_SERVICE)

logs-chart:
	$(COMPOSE) logs -f $(CHART_SERVICE)

run-chart:
	$(COMPOSE) up -d $(CHART_SERVICE)

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

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--create \
		--if-not-exists \
		--topic $(TOPIC_MARKER_CREATED) \
		--partitions 1 \
		--replication-factor 1

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--create \
		--if-not-exists \
		--topic $(TOPIC_RUN_COMPLETED) \
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

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--delete \
		--if-exists \
		--topic $(TOPIC_MARKER_CREATED)

	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--delete \
		--if-exists \
		--topic $(TOPIC_RUN_COMPLETED)

kafka-empty-topics:
	$(MAKE) kafka-delete-topics
	@echo "Waiting for Kafka to delete project topics..."
	# because Kafka topic deletion is often asynchronous
	@until ! $(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--list | grep -Ex \
		-e '$(TOPIC_CANDLE_CLOSED)' \
		-e '$(TOPIC_INPUT_COMPLETED)' \
		-e '$(TOPIC_MARKER_CREATED)' \
		-e '$(TOPIC_RUN_COMPLETED)' >/dev/null; do \
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

kafka-describe-marker-created-topic:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--topic $(TOPIC_MARKER_CREATED)

kafka-describe-run-completed-topic:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--topic $(TOPIC_RUN_COMPLETED)

kafka-candle-closed-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_CANDLE_CLOSED) \
		--from-beginning \
		--max-messages 100 \
		--property print.key=true \
		--property key.separator=' | ' \
		--property print.partition=true \
		--property print.offset=true \
		--property print.timestamp=true

kafka-input-completed-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_INPUT_COMPLETED) \
		--from-beginning \
		--max-messages 100 \
		--property print.key=true \
		--property key.separator=' | ' \
		--property print.partition=true \
		--property print.offset=true \
		--property print.timestamp=true

kafka-marker-created-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_MARKER_CREATED) \
		--from-beginning \
		--max-messages 100 \
		--property print.key=true \
		--property key.separator=' | ' \
		--property print.partition=true \
		--property print.offset=true \
		--property print.timestamp=true

kafka-run-completed-events:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-console-consumer.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--topic $(TOPIC_RUN_COMPLETED) \
		--from-beginning \
		--max-messages 100 \
		--property print.key=true \
		--property key.separator=' | ' \
		--property print.partition=true \
		--property print.offset=true \
		--property print.timestamp=true

kafka-consumer-groups:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-consumer-groups.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--list

kafka-backtest-group:
	$(COMPOSE) exec -T $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-consumer-groups.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--group $(BACKTEST_GROUP_ID)