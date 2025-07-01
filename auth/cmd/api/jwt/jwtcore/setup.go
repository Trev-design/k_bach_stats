package jwtcore

import (
	"auth_server/cmd/api/jwt/jwtmanager"
	"time"

	"github.com/google/uuid"
)

type JWTService struct {
	seeds *jwtmanager.SeedManager
	guid  uuid.UUID
}

type JWTServiceBuilder struct {
	duration time.Duration
	guid     uuid.UUID
}

func NewJWTServiceBuilder() *JWTServiceBuilder {
	return &JWTServiceBuilder{}
}

func (builder *JWTServiceBuilder) Interval(duration time.Duration) *JWTServiceBuilder {
	builder.duration = duration
	return builder
}

func (builder *JWTServiceBuilder) Identifier(guid uuid.UUID) *JWTServiceBuilder {
	builder.guid = guid
	return builder
}

func (builder *JWTServiceBuilder) Build() *JWTService {
	return &JWTService{
		seeds: jwtmanager.NewSeedManager(builder.duration),
		guid:  builder.guid,
	}
}

func (service *JWTService) GetGUID() string { return service.guid.String() }

func (service *JWTService) CloseJWTService() error {
	return service.seeds.CloseSeedManager()
}

func (service *JWTService) ComputeBackgroundService() {
	go service.seeds.ComputeInterval()
}
