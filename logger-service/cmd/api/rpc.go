package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

// RPC is the type for our RPC Server
type RPCServer struct {
}

// RPCpayload is the data we received from RPC
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:     payload.Name,
		Data:     payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed paylod via RPS:" + payload.Name
	return nil

}
