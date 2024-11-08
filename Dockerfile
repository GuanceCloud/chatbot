FROM pubrepo.guance.com/base/ubuntu:20.04 AS base
ARG TARGETARCH

COPY dist/chatbot-api-linux-${TARGETARCH}/ /app/

WORKDIR /app

CMD ["/app/chatbot-api"]
