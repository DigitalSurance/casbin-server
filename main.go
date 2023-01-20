// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate protoc -I proto --go_out=plugins=grpc:proto proto/casbin.proto

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin-server/server"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	options := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := options.NewJSONHandler(os.Stdout)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	var port int
	var jsonPort int
	flag.IntVar(&port, "port", 50051, "grpc listening port")
	flag.IntVar(&jsonPort, "json-port", 50052, "json http api listening port")
	flag.Parse()

	// if port not in range or we can't listen on it, panic and exit
	if port < 1 || port > 65535 {
		panic(fmt.Sprintf("invalid port number: %d", port))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()

	casbinServer := server.NewServer(logger)
	response, err := casbinServer.NewAdapter(context.TODO(), &pb.NewAdapterRequest{})
	if err != nil {
		panic(err)
	}
	casbinServer.NewEnforcer(context.TODO(), &pb.NewEnforcerRequest{ModelText: "", AdapterHandle: response.Handler})

	pb.RegisterCasbinServer(s, casbinServer)

	// spin up json/rest server to handle relation-tuple/check requests by oathkeeper
	jsonServer := server.NewJsonServer(casbinServer)
	go func() {
		logger.Info(fmt.Sprintf("json server listening on port: %d", jsonPort))
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", jsonPort), jsonServer)
		if err != nil {
			panic(err)
		}
	}()

	// Register reflection service on gRPC server.
	reflection.Register(s)
	logger.Info(fmt.Sprintf("grpc server listening on port: %d", port))
	// if we can't serve, panic and exit
	if err := s.Serve(lis); err != nil {
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}
