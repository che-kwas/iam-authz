# Makefile helper functions for tools

TOOLS := golangci-lint mockgen

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.verify.%
tools.verify.%:
	@type $* >/dev/null 2>&1 || make tools.install.$*

.PHONY: tools.install.%
tools.install.%:
	make install.$*

.PHONY: install.golangci-lint
install.golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.mockgen
install.mockgen:
	go install github.com/golang/mock/mockgen@latest
