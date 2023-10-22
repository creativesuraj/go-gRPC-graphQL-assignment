package main

import (
	"net"

	"github.com/aeon/grpc-server/handlers"
	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBookHandlersServer(grpcServer, &handlers.BookHandlersServer{})
	logrus.Info("grpc server started")
	grpcServer.Serve(lis)
}
