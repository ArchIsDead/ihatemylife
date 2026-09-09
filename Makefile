.PHONY: run build update clean install fix start stop restart deps

deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || pkg install golang -y
	@command -v git >/dev/null 2>&1 || pkg install git -y
	@command -v mpv >/dev/null 2>&1 || pkg install mpv -y
	@echo "Dependencies ready."

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
