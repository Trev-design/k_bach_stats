package impl

import (
	"auth_server/cmd/api/broker/producer"
	"auth_server/cmd/api/db/dbcore"
	"auth_server/cmd/api/domain/adapters/brokerimpl"
	"auth_server/cmd/api/domain/adapters/dbimpl"
	"auth_server/cmd/api/domain/adapters/jwtimpl"
	"auth_server/cmd/api/domain/adapters/sessionimpl"
	"auth_server/cmd/api/jwt/jwtcore"
	"auth_server/cmd/api/temporary/session/sessioncore"
	"errors"
)

type Impl struct {
	db      dbimpl.Adapter
	broker  brokerimpl.Adapter
	jwt     jwtimpl.Adapter
	session sessionimpl.Adapter
}

type ImplBuilder struct {
	brokerBuilder   *producer.RMQProducerBuilder
	sessionBuilder  *sessioncore.SessionBuilder
	databaseBuilder *dbcore.DatabaseBuilder
	jwtBuilder      *jwtcore.JWTServiceBuilder
}

func NewImplBuilder() *ImplBuilder {
	return &ImplBuilder{}
}

func (builder *ImplBuilder) DatabaseSetup(dbBuilder *dbcore.DatabaseBuilder) *ImplBuilder {
	builder.databaseBuilder = dbBuilder
	return builder
}

func (builder *ImplBuilder) SessionSetup(sessionBuilder *sessioncore.SessionBuilder) *ImplBuilder {
	builder.sessionBuilder = sessionBuilder
	return builder
}

func (builder *ImplBuilder) BrokerSetup(brokerBuilder *producer.RMQProducerBuilder) *ImplBuilder {
	builder.brokerBuilder = brokerBuilder
	return builder
}

func (builder *ImplBuilder) Build() (*Impl, error) {
	if builder.brokerBuilder == nil ||
		builder.databaseBuilder == nil ||
		builder.jwtBuilder == nil ||
		builder.sessionBuilder == nil {
		return nil, errors.New("uncomplete setup")
	}

	database, err := builder.databaseBuilder.Build()
	if err != nil {
		return nil, err
	}

	session, err := builder.sessionBuilder.Build()
	if err != nil {
		return nil, err
	}

	broker, err := builder.brokerBuilder.Build()
	if err != nil {
		return nil, err
	}

	return &Impl{
		db:      database,
		session: session,
		broker:  broker,
		jwt:     builder.jwtBuilder.Build(),
	}, nil
}

func (impl *Impl) StartBackgroundServices() {
	impl.broker.ComputeBackgroundServices()
	impl.session.HandleBackground()
	impl.jwt.ComputeBackgroundService()
}
