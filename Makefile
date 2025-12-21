BINARY := bus
CMD := ./cmd/main.go

.PHONY: all build run clean

all: run

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

paper:
	python e-paper/main.py

clean:
	rm -f $(BINARY)
