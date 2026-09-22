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
	ps \
	logs \
	logs-kafka \
	logs-topic-init \
	logs-market-data \
	build-market-data \
	rebuild-market-data \
	reset \
	kafka-topics \
	kafka-describe-candle-closed \
	kafka-describe-input-completed \
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

kafka-describe-candle-closed:
	$(COMPOSE) exec $(KAFKA_SERVICE) \
		$(KAFKA_BIN)/kafka-topics.sh \
		--bootstrap-server $(KAFKA_BOOTSTRAP_SERVER) \
		--describe \
		--topic $(TOPIC_CANDLE_CLOSED)

kafka-describe-input-completed:
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