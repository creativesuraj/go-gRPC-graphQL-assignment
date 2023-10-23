package utils

import (
	"errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGRPCConnection() (*grpc.ClientConn, error) {
	opts := make([]grpc.DialOption, 0)
	grpcEnv := LoadEnvConfig("GRPC")
	url := grpcEnv.GetString("SERVER_URL")
	if url == "" {
		return nil, errors.New("GRPC server url is undefined")
	}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		logrus.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return conn, nil
}
