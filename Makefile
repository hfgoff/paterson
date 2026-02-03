BINARY := bus
CMD := ./cmd/main.go

.PHONY: all build run clean

all: run

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY)
	python e-paper/main.py

paper:
	python e-paper/main.py

 # if not in venv, source .venv/bin/activate
fake:
	USE_FAKE_EPD=true python3 e-paper/main.py

clean:
	rm -f $(BINARY)
