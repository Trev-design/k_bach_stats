package sidecar_test

import (
	"auth_server/cmd/api/mocks"
	"auth_server/cmd/api/sidecar"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func Test_BackgroundTopic(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.StartProcessing()

	err := sc.MakeBackgroundProcess("logging", mocks.LoggingPayload([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}

	sc.StopSidecar()
}

func Test_BackgroundTopicMultipleMessages(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.StartProcessing()

	for range 20 {
		err := sc.MakeBackgroundProcess("logging", mocks.LoggingPayload([]byte("hello")))
		if err != nil {
			t.Fatal(err)
		}
	}

	sc.StopSidecar()
}

func Test_BackgroundTopicMultipleTopics(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "more_logging")
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "even_more_logging")

	sc.StartProcessing()

	err := sc.MakeBackgroundProcess("logging", mocks.LoggingPayload([]byte("Hello")))
	if err != nil {
		t.Fatal(err)
	}

	err = sc.MakeBackgroundProcess("more_logging", mocks.LoggingPayload([]byte("Hello")))
	if err != nil {
		t.Fatal(err)
	}

	err = sc.MakeBackgroundProcess("even_more_logging", mocks.LoggingPayload([]byte("Hello")))
	if err != nil {
		t.Fatal(err)
	}

	sc.StopSidecar()
}

func Test_BackgroundTopicMultipleTopicsMultipleMessages(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "more_logging")
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "even_more_logging")

	sc.StartProcessing()

	for range 20 {
		err := sc.MakeBackgroundProcess("logging", mocks.LoggingPayload([]byte("Hello")))
		if err != nil {
			t.Fatal(err)
		}

		err = sc.MakeBackgroundProcess("more_logging", mocks.LoggingPayload([]byte("Hello")))
		if err != nil {
			t.Fatal(err)
		}

		err = sc.MakeBackgroundProcess("even_more_logging", mocks.LoggingPayload([]byte("Hello")))
		if err != nil {
			t.Fatal(err)
		}
	}

	sc.StopSidecar()
}

func Test_BackgroundTopicFailedFalseTopicName(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.StartProcessing()

	err := sc.MakeBackgroundProcess("even_more_logging", mocks.LoggingPayload([]byte("Hello")))
	if err == nil {
		t.Fatal("should fail")
	}

	sc.StopSidecar()
}

func Test_BackgroundTopicFailedSideCarClosed(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterBackgroundTopic(&mocks.LoggingProcessor{}, "logging")
	sc.StartProcessing()
	sc.StopSidecar()

	err := sc.MakeBackgroundProcess("even_more_logging", mocks.LoggingPayload([]byte("Hello")))
	if err == nil {
		t.Fatal("should fail")
	}
}

func Test_ResponsiveTopic(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key")
	sc.StartProcessing()

	err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}

	sc.StopSidecar()
}

func Test_ResponsiveTopicMultipleMessages(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key")
	sc.StartProcessing()

	for range 5 {
		err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
		if err != nil {
			t.Fatal(err)
		}
	}

	sc.StopSidecar()
}

func Test_ResponsiveTopicMultipleTopics(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key")
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key_2")
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key_3")
	sc.StartProcessing()

	err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}

	err = sc.MakeResponsiveProcess("send_key_2", mocks.LoggingPayload([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}

	err = sc.MakeResponsiveProcess("send_key_3", mocks.LoggingPayload([]byte("hello")))
	if err != nil {
		t.Fatal(err)
	}

	sc.StopSidecar()
}

func Test_ResponsiveTopicMultipleTopicsAndMultipleMessages(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key")
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key_2")
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key_3")
	sc.StartProcessing()

	for range 5 {
		err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
		if err != nil {
			t.Fatal(err)
		}

		err = sc.MakeResponsiveProcess("send_key_2", mocks.LoggingPayload([]byte("hello")))
		if err != nil {
			t.Fatal(err)
		}

		err = sc.MakeResponsiveProcess("send_key_3", mocks.LoggingPayload([]byte("hello")))
		if err != nil {
			t.Fatal(err)
		}
	}

	sc.StopSidecar()
}

func Test_ResponsiveTopicFailedFalseToppicName(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.StartProcessing()

	err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
	if err == nil {
		t.Fatal("should faile but succeed")
	}

	sc.StopSidecar()
}

func Test_ResponsiveTopicFailedSideCarClosed(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.StartProcessing()
	sc.StopSidecar()

	err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("hello")))
	if err == nil {
		t.Fatal("should faile but succeed")
	}
}

func Test_ResponsiveTopicFailedTopicError(t *testing.T) {
	sc := sidecar.NewSideCar()
	sc.RegisterResponsiveTopic(&mocks.KeyDistributor{}, "send_key")
	sc.StartProcessing()

	err := sc.MakeResponsiveProcess("send_key", mocks.LoggingPayload([]byte("some invalid payload")))
	if err == nil {
		t.Fatal("should faile but succeed")
	}

	sc.StopSidecar()
}

func Test_RegisterSingleton(t *testing.T) {
	sidecar.NewSideCar().RegisterSingleton()
	sidecar.GetSideCar().StartProcessing()
	sidecar.GetSideCar().StopSidecar()
}
