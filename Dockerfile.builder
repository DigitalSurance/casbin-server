FROM golang:1.23 as BACK

RUN apt-get update && \
    apt-get -y install unzip build-essential autoconf libtool

WORKDIR /go/src
COPY . .

# Install protobuf from source
RUN curl -LjO https://github.com/protocolbuffers/protobuf/archive/refs/tags/v3.17.3.zip && \
    unzip v3.17.3.zip && \
    cd protobuf-3.17.3 && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf-3.17.3 && \
    rm v3.17.3.zip

# Go environment variable to enable Go modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Download dependencies
RUN go mod download

# Install protoc-gen-go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
