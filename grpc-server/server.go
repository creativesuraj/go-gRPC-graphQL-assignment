package main

import (
	"net"

	"github.com/aeon/grpc-server/handlers"
	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/aeon/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	grpcEnv := utils.LoadEnvConfig("GRPC")
	port := grpcEnv.GetString("SERVER_PORT")
	if port == "" {
		logrus.Fatal("GRPC server port is undefined")
	}

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBookHandlersServer(grpcServer, &handlers.BookHandlersServer{})
	logrus.Info("grpc server started")
	grpcServer.Serve(lis)
}
