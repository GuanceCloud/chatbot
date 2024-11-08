.PHONY: chatbot clean

define build_chatbot
	CGO_ENABLED=0 GOOS=linux GOARCH=$(1) \
		go build -o dist/chatbot/linux-$(1)/chatbot cmd/chatbot/main.go
endef

chatbot:
	$(call build_chatbot,amd64)
	$(call build_chatbot,arm64)

clean:
	rm -rf dist/
