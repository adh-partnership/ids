SUBDIVISION ?= zan

.PHONY: dev
dev:
	@npx vite --port 5174 --host 0.0.0.0

build:
	cp configs/$(SUBDIVISION).json config.json
	yarn build
