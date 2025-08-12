# build a tiny docker image 
FROM alpine:latest

RUN mkdir /app

#COPY --from=builder /app/brokerApp /app
COPY frontEndApp /app
#COPY templates /templates

CMD [ "/app/frontEndApp" ]