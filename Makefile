BINARY := bus
CMD := ./cmd/main.go

.PHONY: all build run clean

all: run

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
