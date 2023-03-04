setup:
	@echo "Setting up the environment"
	@./scripts/setup.sh

cibuild:
	./scripts/cibuild.sh

#####################################

BINARY=govel
SRC=./main.go
DB_SRC=./velocity.db
BIN_DIR=./bin
.DEFAULT_GOAL := run

build:
	@go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY) $(SRC) >/dev/null

run: build
	$(BIN_DIR)/$(BINARY)

test:
	go test ./... -v

clean:
	go clean
	rm -rf $(BIN_DIR)
	rm  $(DB_SRC)

#####################################

LINUX_ARM_DIR=./bin/linux-arm

$(BINARY)-rasp:
	@GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o $(LINUX_ARM_DIR)/$(BINARY) $(SRC) >/dev/null
