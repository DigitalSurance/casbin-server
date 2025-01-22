FROM casbin-builder:0.0.0 AS builder

# Copy the source and generate the .proto file
ADD . /go/src/github.com/casbin/casbin-server
WORKDIR $GOPATH/src/github.com/casbin/casbin-server
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false \
    --go-grpc_opt=paths=source_relative proto/casbin.proto

# Install app
RUN go install .

RUN cd /go/src && go build -o casbin-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/casbin-server /app/
COPY ./config/connection_for_docker.json ./connection_config.json
COPY ./examples/rbac_model.conf ./rbac_model.conf

EXPOSE 50051 50052
ENTRYPOINT ./casbin-server

