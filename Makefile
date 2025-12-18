BACKEND_DIR := backend
STORAGE_DIR := storage
WEB_DIR := web

backend:
	cd $(BACKEND_DIR) && go run .

storage:
	cd $(STORAGE_DIR) && go run .

web:
	cd $(WEB_DIR) && npm run dev

dev:
	@echo "Starting backend..."
	cd $(BACKEND_DIR) && go run . &
	@echo "Starting storage..."
	cd $(STORAGE_DIR) && go run . &
	@echo "Starting web..."
	cd $(WEB_DIR) && npm run dev &
	@echo "All services started."
	wait

