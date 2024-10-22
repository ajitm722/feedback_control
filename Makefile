# Define variables
DOCKER_COMPOSE_FILE = docker-compose.yaml
PRODUCER_EXEC = producer
CONSUMER_EXEC = consumer
SIMULATION_SCRIPT = simulation.py

# Default target: Starts the entire process
.PHONY: all
all: docker-start build-producer build-consumer

# Start Docker Compose in detached mode
.PHONY: docker-start
docker-start:
	@echo "Starting Docker Compose in detached mode..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Build the producer executable
.PHONY: build-producer
build-producer:
	@echo "Building the Kafka producer..."
	cd client && go build -o ../$(PRODUCER_EXEC) producer.go

# Build the consumer executable
.PHONY: build-consumer
build-consumer:
	@echo "Building the Kafka consumer..."
	cd server && go build -o ../$(CONSUMER_EXEC) consumer.go

# Run the simulation (Assuming simulation.py visualizes the CSV data)
.PHONY: run-simulation
run-simulation:
	@echo "Running the temperature simulation..."
	python3 $(SIMULATION_SCRIPT)

# Clean up executables and stop Docker Compose
.PHONY: clean
clean:
	@echo "Stopping Docker Compose and removing executables..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	rm -f $(PRODUCER_EXEC) $(CONSUMER_EXEC)

# Force stop all Docker containers
.PHONY: docker-stop
docker-stop:
	@echo "Stopping all running Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Help target for listing all the available commands
.PHONY: help
help:
	@echo "Makefile usage:"
	@echo "  make all             - Start Docker, build producer & consumer"
	@echo "  make docker-start    - Start Docker Compose in detached mode"
	@echo "  make build-producer  - Build the producer executable"
	@echo "  make build-consumer  - Build the consumer executable"
	@echo "  make run-simulation  - Run the simulation script"
	@echo "  make clean           - Stop Docker and remove executables"
	@echo "  make docker-stop     - Stop Docker containers only"

