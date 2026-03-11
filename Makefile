include .env
export

IMAGE_NAME = transya

.PHONY: build run

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run --rm \
    -e "FOLDER_ID=$(FOLDER_ID)" \
    -e "API_KEY=$(API_KEY)" \
    -e "NATS_URL=$(NATS_URL)" \
    -e "KAFKA_ADDR=$(KAFKA_ADDR)" \
    -e "KAFKA_TOPIC=$(KAFKA_TOPIC)" \
    $(IMAGE_NAME)
