## base go image 
#FROM golang:1.24-alpine AS builder

#RUN mkdir /app

#COPY . /app

#WORKDIR /app

#RUN GO_ENABLED=0 go build -o brokerApp ./cmd/api

#RUN chmod +x /app/brokerApp

# build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

#COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

CMD [ "/app/brokerApp" ]