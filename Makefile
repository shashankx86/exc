# Define the name of the output binary
BINARY_NAME=exc.out

# Define the Go compiler command
GO_BUILD=go build

# Default to RELEASE build if not specified
BUILD_TYPE ?= DEV

# Set build flags based on the build type
ifeq ($(BUILD_TYPE), DEV)
    GO_BUILD_FLAGS=-ldflags "-X main.dev=1"
else
    GO_BUILD_FLAGS=
endif

all: $(BINARY_NAME)

$(BINARY_NAME):
	$(GO_BUILD) $(GO_BUILD_FLAGS) -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
