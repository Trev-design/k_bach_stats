package keyorchestrator

import (
	"auth_server/cmd/api/grpc/keyorchestrator/proto/proto"
	"auth_server/cmd/api/tlsconf"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client           proto.KeyOrchestratorServiceClient
	connection       *grpc.ClientConn
	serviceWaitgroup *sync.WaitGroup
}

type GRPCClientBuilder struct {
	grpcHost  string
	grpcPort  string
	tlsConfig *tlsconf.TLSBuilder
}

func (builder *GRPCClientBuilder) GRPCHost(host string) *GRPCClientBuilder {
	builder.grpcHost = host
	return builder
}

func (builder *GRPCClientBuilder) GRPCPort(port string) *GRPCClientBuilder {
	builder.grpcPort = port
	return builder
}

func (builder *GRPCClientBuilder) WithTLS(tlsConfig *tlsconf.TLSBuilder) *GRPCClientBuilder {
	builder.tlsConfig = tlsConfig
	return builder
}

func (builder *GRPCClientBuilder) Build() (*GRPCClient, error) {
	if builder.tlsConfig != nil {
		config, err := builder.tlsConfig.Build()
		if err != nil {
			return nil, err
		}

		return builder.buildClient(credentials.NewTLS(config))
	}

	return builder.buildClient(insecure.NewCredentials())
}

func (builder *GRPCClientBuilder) buildClient(credentials credentials.TransportCredentials) (*GRPCClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", builder.grpcHost, builder.grpcPort), grpc.WithTransportCredentials(credentials))
	if err != nil {
		return nil, err
	}

	client := proto.NewKeyOrchestratorServiceClient(conn)

	return &GRPCClient{
		client:           client,
		connection:       conn,
		serviceWaitgroup: &sync.WaitGroup{},
	}, nil
}
