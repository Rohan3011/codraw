# Go compiler
GO := go

# Directory where the source files are located
SRCDIR := cmd

# Output directory for built binaries
OUTDIR := bin

# List of source files (assuming all .go files in the src directory)
SOURCES := $(wildcard $(SRCDIR)/*.go)

# Name of the final executable (you can change this to your desired name)
TARGET := myapp

# Build command
build:
	@echo "Building $(TARGET)..."
	$(GO) build -o $(OUTDIR)/$(TARGET) $(SOURCES)

# Run command
run: build
	@echo "Running $(TARGET)..."
	$(OUTDIR)/$(TARGET)

# Clean command
clean:
	@echo "Cleaning up..."
	rm -rf $(OUTDIR)

# Help command (to show available commands)
help:
	@echo "Available commands:"
	@echo "  make build     - Build the project"
	@echo "  make run       - Build and run the project"
	@echo "  make clean     - Clean up built binaries"