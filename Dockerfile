# Build Stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go build -o main cmd/main/main.go
RUN cp `which migrate` /app/migrate

# Run stage
FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]