run-video-service:
	@echo "Running video service..."
	cd cmd/video && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/nghiatrann0502/trinity/cmd/video
.PHONY: run-video-service

run-ranking-service:
	@echo "Running ranking service..."
	cd cmd/ranking && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/nghiatrann0502/trinity/cmd/ranking
.PHONY: run-ranking-service


run:
	docker compose --profile default up --build  -d
.PHONY: run

stop:
	docker compose --profile default down
.PHONY: stop

run-core:
	docker compose --profile core up --build  -d
.PHONY: run-core

stop-core:
	docker compose --profile core down
.PHONY: stop-core