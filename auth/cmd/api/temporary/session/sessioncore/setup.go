package sessioncore

import (
	"auth_server/cmd/api/temporary/session/sessioncrypto"
	"auth_server/cmd/api/temporary/session/sessionstore"
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/utils/connection"
	"sync"
	"time"
)

type Session struct {
	cookieCrypt  *sessioncrypto.Crypt
	verifyCrypt  *sessioncrypto.Crypt
	refreshCrypt *sessioncrypto.Crypt
	store        *sessionstore.RedisClient
	mutex        sync.RWMutex
}

type SessionBuilder struct {
	host               string
	port               string
	password           string
	tlsBuilder         *tlsconf.TLSBuilder
	intervalDuration   time.Duration
	sessionDuration    time.Duration
	credentialsChannel chan connection.Credentials
}

func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{
		host:               "localhost",
		port:               "6379",
		sessionDuration:    2 * time.Hour,
		intervalDuration:   15 * time.Minute,
		credentialsChannel: nil,
	}
}

func (builder *SessionBuilder) Host(host string) *SessionBuilder {
	builder.host = host
	return builder
}

func (builder *SessionBuilder) Port(port string) *SessionBuilder {
	builder.port = port
	return builder
}

func (builder *SessionBuilder) Password(password string) *SessionBuilder {
	builder.password = password
	return builder
}

func (builder *SessionBuilder) WithTLS(tlsBuilder *tlsconf.TLSBuilder) *SessionBuilder {
	builder.tlsBuilder = tlsBuilder
	return builder
}

func (builder *SessionBuilder) IntevalDuration(duration time.Duration) *SessionBuilder {
	builder.intervalDuration = duration
	return builder
}

func (builder *SessionBuilder) SessionDuration(duration time.Duration) *SessionBuilder {
	builder.sessionDuration = duration
	return builder
}

func (builder *SessionBuilder) WithCredentialChannel(pipe chan connection.Credentials) *SessionBuilder {
	builder.credentialsChannel = pipe
	return builder
}

func (builder *SessionBuilder) Build() (*Session, error) {
	client := new(sessionstore.RedisClient)
	newClient, err := sessionstore.NewRedisClientBuilder().
		Host(builder.host).
		Port(builder.port).
		Password(builder.password).
		WithDuration(builder.sessionDuration).
		WithTLS(builder.tlsBuilder).
		WithCredentialChannel(builder.credentialsChannel).
		Build()
	if err != nil {
		return nil, err
	}

	client = newClient

	verifyCrypt, err := sessioncrypto.NewCrypt(builder.intervalDuration)
	if err != nil {
		return nil, err
	}

	cookieCrypt, err := sessioncrypto.NewCrypt(builder.intervalDuration)
	if err != nil {
		return nil, err
	}

	refreshCrypt, err := sessioncrypto.NewCrypt(builder.intervalDuration)
	if err != nil {
		return nil, err
	}

	return &Session{
		store:        client,
		verifyCrypt:  verifyCrypt,
		cookieCrypt:  cookieCrypt,
		refreshCrypt: refreshCrypt,
		mutex:        sync.RWMutex{},
	}, nil
}
