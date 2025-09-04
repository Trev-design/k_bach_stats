package userregistry

import (
	"auth_server/cmd/api/grpc/userregistry/proto"
	"auth_server/cmd/api/tlsconf"
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client                     proto.UserRegistryServiceClient
	connection                 *grpc.ClientConn
	maxNumPrimaryRequests      int64
	maxNumOverflowRequests     int64
	currentNumPrimaryRequests  atomic.Int64
	currentNumOverflowRequests atomic.Int64
	isServiceFinished          atomic.Bool
	promptRequestChannel       chan struct{}
	primaryRequestDoneChannel  chan struct{}
	overflowRequestDoneChannel chan struct{}
	primaryReadyChannel        chan struct{}
	overflowReadyChannel       chan struct{}
	messageIncomeChannel       chan Request
	messagePrimaryChannel      chan Request
	messageOverflowChannel     chan Request
	serviceFinishedChannel     chan struct{}
	serviceWaitgroup           *sync.WaitGroup
}

type ClientBuilder struct {
	maxNumPrimaryRequests  int64
	maxNumOverflowRequests int64
	grpcHost               string
	grpcPort               string
	tlsConfig              *tlsconf.TLSBuilder
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		maxNumPrimaryRequests:  5,
		maxNumOverflowRequests: 5,
		grpcHost:               "localhost",
		grpcPort:               "80",
		tlsConfig:              nil,
	}
}

func (builder *ClientBuilder) MaxNumPrimaryRequests(maxPrimaryRequests int64) *ClientBuilder {
	builder.maxNumPrimaryRequests = maxPrimaryRequests
	return builder
}

func (builder *ClientBuilder) MaxNumOverflowRequests(maxOverflowRequests int64) *ClientBuilder {
	builder.maxNumOverflowRequests = maxOverflowRequests
	return builder
}

func (builder *ClientBuilder) GRPCHost(host string) *ClientBuilder {
	builder.grpcHost = host
	return builder
}

func (builder *ClientBuilder) GRPCPort(port string) *ClientBuilder {
	builder.grpcPort = port
	return builder
}

func (builder *ClientBuilder) WithTLS(tlsConfig *tlsconf.TLSBuilder) *ClientBuilder {
	builder.tlsConfig = tlsConfig
	return builder
}

func (builder *ClientBuilder) Build() (*GRPCClient, error) {
	if builder.tlsConfig != nil {
		config, err := builder.tlsConfig.Build()
		if err != nil {
			return nil, err
		}

		return builder.buildClient(credentials.NewTLS(config))
	}

	return builder.buildClient(insecure.NewCredentials())
}

func (builder *ClientBuilder) buildClient(credentials credentials.TransportCredentials) (*GRPCClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", builder.grpcHost, builder.grpcPort), grpc.WithTransportCredentials(credentials))
	if err != nil {
		return nil, err
	}

	client := proto.NewUserRegistryServiceClient(conn)

	return &GRPCClient{
		client:                     client,
		connection:                 conn,
		maxNumPrimaryRequests:      builder.maxNumPrimaryRequests,
		maxNumOverflowRequests:     builder.maxNumOverflowRequests,
		currentNumPrimaryRequests:  atomic.Int64{},
		currentNumOverflowRequests: atomic.Int64{},
		isServiceFinished:          atomic.Bool{},
		promptRequestChannel:       make(chan struct{}, builder.maxNumPrimaryRequests+builder.maxNumOverflowRequests),
		primaryRequestDoneChannel:  make(chan struct{}, builder.maxNumPrimaryRequests),
		overflowRequestDoneChannel: make(chan struct{}, builder.maxNumOverflowRequests),
		primaryReadyChannel:        make(chan struct{}, builder.maxNumPrimaryRequests),
		overflowReadyChannel:       make(chan struct{}, builder.maxNumOverflowRequests),
		messageIncomeChannel:       make(chan Request, builder.maxNumPrimaryRequests+builder.maxNumOverflowRequests),
		messagePrimaryChannel:      make(chan Request, builder.maxNumPrimaryRequests),
		messageOverflowChannel:     make(chan Request, builder.maxNumOverflowRequests),
		serviceFinishedChannel:     make(chan struct{}),
		serviceWaitgroup:           &sync.WaitGroup{},
	}, nil
}

func (client *GRPCClient) HandleBackgroundServices() {
	go client.handleMessageLoad()
	go client.handleMessage()
	go client.computePrimaryStream()
	go client.computeOverflowStream()
}
