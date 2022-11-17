FROM golang:1.16 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o start-casbin-server

FROM alpine:3.14
WORKDIR /casbin-server
COPY --from=builder /app/start-casbin-server ./
COPY ./config/connection_for_docker.json ./connection_config.json
COPY ./examples/rbac_model.conf ./rbac_model.conf
EXPOSE 50051 50052
RUN ["chmod", "+x", "./start-casbin-server"]
CMD ["./start-casbin-server"]