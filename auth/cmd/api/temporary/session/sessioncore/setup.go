package sessioncore

import (
	"auth_server/cmd/api/temporary/session/sessioncrypto"
	"auth_server/cmd/api/temporary/session/sessionstore"
	"auth_server/cmd/api/tlsconf"
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
	host       string
	port       string
	password   string
	tlsBuilder *tlsconf.TLSBuilder
	duration   time.Duration
}

func NewSessionBuilder() *SessionBuilder {
	return &SessionBuilder{
		host:     "localhost",
		port:     "6379",
		duration: 2 * time.Hour,
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
	builder.duration = duration
	return builder
}

func (builder *SessionBuilder) Build() (*Session, error) {
	client := new(sessionstore.RedisClient)
	if builder.tlsBuilder != nil {
		tlsConfig, err := builder.tlsBuilder.Build()
		if err != nil {
			return nil, err
		}

		newClient, err := sessionstore.NewRedisClient(
			builder.duration,
			builder.password,
			builder.host,
			builder.port,
			tlsConfig)
		if err != nil {
			return nil, err
		}

		client = newClient
	} else {
		newClient, err := sessionstore.NewRedisClient(
			builder.duration,
			builder.password,
			builder.host,
			builder.port,
			nil)
		if err != nil {
			return nil, err
		}

		client = newClient
	}

	verifyCrypt, err := sessioncrypto.NewCrypt(builder.duration)
	if err != nil {
		return nil, err
	}

	cookieCrypt, err := sessioncrypto.NewCrypt(builder.duration)
	if err != nil {
		return nil, err
	}

	refreshCrypt, err := sessioncrypto.NewCrypt(builder.duration)
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
