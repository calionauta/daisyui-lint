.PHONY: all test lint vocabulary build clean

all: test lint build

test:
	@go test ./...

lint:
	@gofumpt -l . | tee /tmp/gofumpt.$$$$ | (! grep .) || (cat /tmp/gofumpt.$$$$; exit 1)
	@go vet ./...

# Regenerates the snapshot embedded in the binary from the vendored skill. The
# snapshot is committed; TestBundledVocabularyMatchesSkill fails when it is
# stale, which is what makes a committed generated artifact safe.
vocabulary:
	@echo "→ regenerating daisyui/vocabulary.json..."
	@tmp=$$(mktemp) && \
		if go run ./cmd/daisyui-lint --emit-vocabulary > "$$tmp" && [ -s "$$tmp" ]; then \
			mv "$$tmp" daisyui/vocabulary.json && \
			echo "✅ $$(grep -c '"base"' daisyui/vocabulary.json) components"; \
		else \
			rm -f "$$tmp"; echo "❌ generation failed; the committed snapshot is untouched"; exit 1; \
		fi

build:
	@go build -o /dev/null ./...

clean:
	@rm -f daisyui-lint
