package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//Write a log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// everything went well
	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}

func (app *Config) gRPCListen() {
	lstnr, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
	srv := grpc.NewServer()

	//register

	logs.RegisterLogServiceServer(srv, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %v", gRpcPort)

	if err := srv.Serve(lstnr); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

}
