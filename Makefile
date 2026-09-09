.PHONY: run build update clean install fix start stop restart deps

deps:
	@echo ""
	@echo "[*] Checking Go..."
	@if command -v go >/dev/null 2>&1; then \
		echo "    Go: INSTALLED"; \
	else \
		echo "    Go: NOT FOUND - Installing..."; \
		pkg install golang -y; \
	fi
	@echo ""
	@echo "[*] Checking Git..."
	@if command -v git >/dev/null 2>&1; then \
		echo "    Git: INSTALLED"; \
	else \
		echo "    Git: NOT FOUND - Installing..."; \
		pkg install git -y; \
	fi
	@echo ""
	@echo "[*] Checking MPV..."
	@if command -v mpv >/dev/null 2>&1; then \
		echo "    MPV: INSTALLED"; \
	else \
		echo "    MPV: NOT FOUND - Installing..."; \
		pkg install mpv -y; \
	fi
	@echo ""
	@echo "[*] All dependencies ready."
	@echo ""

run: deps
	@clear
	@go run .

build: deps
	@clear
	@go build -o s .

update: deps
	@clear
	@git pull
	@rm -f s
	@go mod tidy
	@go build -o s .
	@./s

clean:
	@clear
	@rm -f s
	@go clean

install: deps
	@clear
	@go mod tidy
	@go build -o s .

fix:
	@clear
	@rm -rf ~/suic1.de
	@rm -f s
	@go mod tidy
	@go build -o s .

start: deps
	@clear
	@./s

stop:
	@pkill -f "./s" || true

restart: stop start
