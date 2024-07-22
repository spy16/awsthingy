build:
	@echo "Building..."
	@go build -o bin/ ./...

install:
	@echo "Installing..."
	@go install ./...
