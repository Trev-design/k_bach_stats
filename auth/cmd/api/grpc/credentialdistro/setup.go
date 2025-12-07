package credentialdistro

import (
	"auth_server/cmd/api/grpc/credentialdistro/proto"
	"auth_server/cmd/api/kms/local/asymmetric"
	"auth_server/cmd/api/utils/connection"
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

type GRPCClient struct {
	connection     *grpc.ClientConn
	waitgroup      *sync.WaitGroup
	client         proto.CredentialDistroServiceClient
	saltStreams    map[string]*grpcSaltStream
	newCredsStream *grpcNewCredsStreamHandler
	mutex          sync.RWMutex
}

type grpcSaltStream struct {
	id             string
	topic          string
	stream         grpc.BidiStreamingClient[proto.SaltRequest, proto.Response]
	messageChannel chan []byte
}

type grpcNewCredsStreamHandler struct {
	streams    map[string]*grpcNewCredsStream
	scheduled  map[string]bool
	keyManager *asymmetric.KeyManager
	mutex      sync.Mutex
	waitgroup  *sync.WaitGroup
}

type grpcNewCredsStream struct {
	id              string
	topic           string
	stream          grpc.BidiStreamingClient[proto.NewCredsRequest, proto.Response]
	messageChannels chan *connection.Credentials
}
type SaltStreamHandler interface {
	MakeSaltStream(topic, id string, messageChannel chan []byte)
}

type NewCredsStreamHandler interface {
	MakeNewCredsStream(topic string, messageChannel chan connection.Credentials)
}

func NewGRPCClient(host, port string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	client := proto.NewCredentialDistroServiceClient(conn)

	keyManager, err := asymmetric.NewKeyManager()
	if err != nil {
		return nil, err
	}

	newCredsStream := &grpcNewCredsStreamHandler{
		streams:    make(map[string]*grpcNewCredsStream),
		scheduled:  make(map[string]bool),
		keyManager: keyManager,
		mutex:      sync.Mutex{},
		waitgroup:  &sync.WaitGroup{},
	}

	return &GRPCClient{
		connection:     conn,
		waitgroup:      &sync.WaitGroup{},
		client:         client,
		saltStreams:    make(map[string]*grpcSaltStream),
		newCredsStream: newCredsStream,
		mutex:          sync.RWMutex{},
	}, nil
}
