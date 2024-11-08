.PHONY: chatbot clean

define build_chatbot
	CGO_ENABLED=0 GOOS=linux GOARCH=$(1) \
		go build -o dist/chatbot-api-linux-$(1)/chatbot-api cmd/chatbot-api/main.go
endef

all: chatbot

chatbot:
	$(call build_chatbot,amd64)
	$(call build_chatbot,arm64)

build_image: clean chatbot
	docker buildx build --platform amd64,arm64 --tag \
		pubrepo.jiagouyun.com/chatbot/chatbot-api:$(shell git describe --always --tags)_dev . --push

build_image_prod: clean chatbot
	docker buildx build --platform amd64,arm64 --tag \
		pubrepo.jiagouyun.com/chatbot/chatbot-api:$(shell git describe --always --tags) . --push

clean:
	rm -rf dist/
