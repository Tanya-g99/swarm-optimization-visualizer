.PHONY: all frontend backend start stop

# Переменные
FRONTEND_DIR=frontend
BACKEND_DIR=backend-go
TEST_DIR=tests
PERFORMANCE_TEST_DIR=tests/test
VENV_DIR=.venv

all: start

# Запуск фронтенда
frontend:
	cd $(FRONTEND_DIR) && npm run dev

# Запуск бэкенда
backend:
	cd $(BACKEND_DIR) && go run .
# cd $(BACKEND_DIR) && go run $(BACKEND_FILES)

# Запуск тестов на Go
test:
	cd $(BACKEND_DIR)/$(TEST_DIR) && pytest test.py

test-performance:
	cd $(BACKEND_DIR)/$(PERFORMANCE_TEST_DIR) && go test

install-deps:
	if [ ! -d "$(VENV_DIR)" ]; then \
		python3 -m venv $(VENV_DIR); \
	fi

	source $(VENV_DIR)/bin/activate && pip install -r $(TEST_DIR)/requirements.txt