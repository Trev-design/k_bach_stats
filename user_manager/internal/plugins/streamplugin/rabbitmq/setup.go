package rabbitmq

import (
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func makeEnv() (*StreamAdapter, error) {
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("IAmTheUser").
			SetPassword("ThisIsmyPassword").
			SetVHost("kbach"),
	)
	if err != nil {
		return nil, err
	}

	if err = env.DeclareStream(
		"invitation",
		stream.NewStreamOptions().
			SetMaxSegmentSizeBytes(stream.ByteCapacity{}.MB(3)).
			SetMaxLengthBytes(stream.ByteCapacity{}.MB(12)),
	); err != nil {
		return nil, err
	}

	producerOptions := stream.NewProducerOptions()
	producerOptions.SetProducerName("invitation_archive_producer")

	invitationProducer, err := env.NewProducer("invitation", producerOptions)
	if err != nil {
		return nil, err
	}

	if err = env.DeclareStream(
		"join_request",
		stream.NewStreamOptions().
			SetMaxSegmentSizeBytes(stream.ByteCapacity{}.MB(2)).
			SetMaxLengthBytes(stream.ByteCapacity{}.MB(8)),
	); err != nil {
		return nil, err
	}

	options := stream.NewProducerOptions()
	options.SetProducerName("join_request_archive_producer")

	joinRequestProducer, err := env.NewProducer("join_request", producerOptions)
	if err != nil {
		return nil, err
	}

	return &StreamAdapter{
		environment:         env,
		invitationProducer:  invitationProducer,
		joinRequestProducer: joinRequestProducer,
	}, nil
}
