
runtest:
	@echo "Starting Test Suite..."
	docker-compose -f test/docker-compose.yml up --abort-on-container-exit --build

run:
	@echo "Starting docker compose"
	docker-compose up --build

clean:
	@echo "Cleaning files..."
	docker-compose down -v --remove-orphans
