.PHONY: install
install:
	go install

.PHONY: release
release:
	GITHUB_TOKEN=$GORELEASER_GITHUB_TOKEN goreleaser --rm-dist
