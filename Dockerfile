# stage Build
FROM golang:1.21.1-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# stage Run
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8000
CMD [ "/app/main" ]
# it works like /app/start.sh /app/main ENTRYPOINT CMD
ENTRYPOINT [ "/app/start.sh" ]
