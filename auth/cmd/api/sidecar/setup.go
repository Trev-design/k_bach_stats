package sidecar

import (
	"sync"
	"sync/atomic"
)

type backgroundTopic struct {
	messageChannel chan Payload
	handler        BackgroundSignal
	isFinished     atomic.Bool
	waitgroup      *sync.WaitGroup
	once           sync.Once
}

type responsiveTopic struct {
	messageChannel chan responsiveRequest
	handle         ResponsiveSignal
	isFinished     atomic.Bool
	waitgroup      *sync.WaitGroup
	once           sync.Once
}

type responsiveRequest struct {
	message      Payload
	errorChannel chan error
}

type SideCar interface {
	MakeBackgroundProcess(topic string, message Payload) error
	MakeResponsiveProcess(topic string, message Payload) error
	RegisterBackgroundTopic(signal BackgroundSignal, name string)
	RegisterResponsiveTopic(signal ResponsiveSignal, name string)
	StartProcessing()
	StopSidecar()
}

type processor struct {
	backgroundTopics map[string]*backgroundTopic
	responsiveTopics map[string]*responsiveTopic
	waitgroup        *sync.WaitGroup
}

var sidecar *processor
var once sync.Once

func GetSideCar() SideCar {
	return sidecar
}

func NewSideCar() *processor {
	return &processor{
		backgroundTopics: make(map[string]*backgroundTopic),
		responsiveTopics: make(map[string]*responsiveTopic),
		waitgroup:        &sync.WaitGroup{},
	}
}

func (processor *processor) RegisterSingleton() {
	once.Do(func() {
		sidecar = processor
	})
}
