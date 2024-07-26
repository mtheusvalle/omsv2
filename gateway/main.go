package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/mtheusvalle/commons"
	pb "github.com/mtheusvalle/commons/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr          = common.EnvString("HTTP_ADDR", ":8080")
	ordersServiceAddr = "localhost:2000"
)

func main() {
	conn, err := grpc.Dial(ordersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing orders serice at", ordersServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server")
	}
}
