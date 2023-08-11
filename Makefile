backtest:
	@echo "Starting Backend Test Suite..."
	docker-compose -f backend/test/docker-compose.yml up --abort-on-container-exit --build

authtest:
	@echo "Starting Auth Test Suite..."
	docker-compose -f auth/test/docker-compose.yml up --abort-on-container-exit --build

run:
	@echo "Starting docker compose"
	docker-compose -f docker-compose.yml up --build

lint:
	@echo "Starting golangci-lint..."
	cd backend && golangci-lint run 

clean:
	@echo "Cleaning files..."
	docker-compose down -v --remove-orphans
