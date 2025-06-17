#
# Display help.
#
.PHONY: help
help:
	@echo 'Available tasks:'
	@echo '  make help     -- Display this help message'
	@echo '  make build    -- Build the `wholidisuka` binary'
	@echo '  make clean    -- Remove build artifacts'
	@echo '  make test     -- Run tests'
	@echo '  make license  -- Collect license information and save it to `./licenses/`'
	@echo '  make release  -- Release the wholidisuka binary'

#
# Build the binary.
#
.PHONY: build
build: clean
	@echo 'Building wholidisuka binary...'

	go build \
	  -v \
	  ./cmd/wholidisuka/

#
# Remove build artifacts.
#
.PHONY: clean
clean:
	@echo 'Cleaning build artifacts...'

	rm -rf \
	  ./wholidisuka \
	  ./dist/ \
	  ./licenses/

#
# Run tests.
#
.PHONY: test
test:
	@echo 'Running tests...'

	go test -v ./...

#
# Collect license information.
#
.PHONY: license
license:
	@echo 'Collecting license information...'

	go tool go-licenses save ./cmd/wholidisuka/ \
	  --force \
	  --save_path ./licenses/

#
# Release the binary.
#
.PHONY: release
release: clean license
	@echo 'Releasing wholidisuka binary...'

	go tool goreleaser release --clean
