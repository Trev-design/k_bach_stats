package sessionstore

import (
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/utils/connection"
	"crypto/tls"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	client    *redis.Client
	waitgroup *sync.WaitGroup
}
type RedisClient struct {
	conn              *connection.Handler[RedisConnection]
	expiry            time.Duration
	builder           *RedisClientBuilder
	credentialChannel chan connection.Credentials
}

type RedisClientBuilder struct {
	password  string
	host      string
	port      string
	tlsConfig *tlsconf.TLSBuilder
	expiry    time.Duration
	pipe      chan connection.Credentials
}

func NewRedisClientBuilder() *RedisClientBuilder {
	return &RedisClientBuilder{}
}

func (builder *RedisClientBuilder) Password(password string) *RedisClientBuilder {
	builder.password = password
	return builder
}

func (builder *RedisClientBuilder) Host(host string) *RedisClientBuilder {
	builder.host = host
	return builder
}

func (builder *RedisClientBuilder) Port(port string) *RedisClientBuilder {
	builder.port = port
	return builder
}

func (builder *RedisClientBuilder) WithTLS(tlsConf *tlsconf.TLSBuilder) *RedisClientBuilder {
	builder.tlsConfig = tlsConf
	return builder
}

func (builder *RedisClientBuilder) WithDuration(duration time.Duration) *RedisClientBuilder {
	builder.expiry = duration
	return builder
}

func (builder *RedisClientBuilder) WithCredentialChannel(pipe chan connection.Credentials) *RedisClientBuilder {
	builder.pipe = pipe
	return builder
}

func (builder *RedisClientBuilder) BuildConnection() (connection.Connection[RedisConnection], error) {
	var tlsconfig *tls.Config = nil

	if builder.tlsConfig != nil {
		newTLSConfig, err := builder.tlsConfig.Build()
		if err != nil {
			return nil, err
		}

		tlsconfig = newTLSConfig
	}

	client, err := setupClient(
		builder.password,
		builder.host,
		builder.port,
		tlsconfig,
	)
	if err != nil {
		return nil, err
	}

	return &RedisConnection{
		client:    client,
		waitgroup: &sync.WaitGroup{},
	}, nil
}

func (builder *RedisClientBuilder) Build() (*RedisClient, error) {
	conn, err := connection.NewBuilder(builder).Build()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		builder:           builder,
		credentialChannel: builder.pipe,
		expiry:            builder.expiry,
		conn:              conn,
	}, nil
}

func (client *RedisClient) CloseRedisStore() error {
	conn := client.conn.Get()
	conn.waitgroup.Wait()

	return conn.Close()
}
