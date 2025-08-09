package main

import (
	"billing-service/internal/api"
	grpcserver "billing-service/internal/grpc"
	"billing-service/internal/storage"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	pb "billing-service/proto"

	"github.com/gin-gonic/gin"
)

func main() {
	redisAddr := "localhost:6379"
	redisClient := storage.NewRedisClient(redisAddr)

	// Start REST API server
	go func() {
		r := gin.Default()
		handler := api.NewHandler(redisClient)
		r.POST("/charge/:userid", handler.ChargeUser)
		log.Println("REST API listening on :8080")
		if err := r.Run(":8081"); err != nil {
			log.Fatalf("failed to start REST API: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatalf("failed to listen on :4040: %v", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterBillingServiceServer(grpcSrv, grpcserver.NewServer(redisClient))

	go func() {
		log.Println("gRPC server listening on :4040")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Graceful shutdown on interrupt
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down servers...")
	grpcSrv.GracefulStop()
	// REST server has no graceful shutdown here for brevity
	time.Sleep(time.Second)
	log.Println("Shutdown complete")
}
