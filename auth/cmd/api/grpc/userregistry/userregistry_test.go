package userregistry_test

import (
	"auth_server/cmd/api/grpc/userregistry"
	"auth_server/cmd/api/mocks"
	"log"
	"sync"
	"testing"
	"time"
)

var testClient *userregistry.GRPCClient

func TestMain(m *testing.M) {
	go mocks.NewRegistryServer(10, 5)

	time.Sleep(1 * time.Second)

	client, err := userregistry.NewClientBuilder().
		GRPCHost("localhost").
		GRPCPort("5670").
		MaxNumPrimaryRequests(10).
		MaxNumOverflowRequests(5).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	testClient = client

	testClient.HandleBackgroundServices()

	defer testClient.CloseGRPCClient()

	m.Run()
}

func TestInitClient(t *testing.T) {
	newClient, err := userregistry.NewClientBuilder().
		GRPCHost("localhost").
		GRPCPort("5670").
		MaxNumPrimaryRequests(10).
		MaxNumOverflowRequests(5).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	newClient.HandleBackgroundServices()

	newClient.CloseGRPCClient()
}

func TestSend(t *testing.T) {
	if err := testClient.SendMessage(userregistry.Message{
		Name:   "halli",
		Email:  "some@email.com",
		Entity: "hallo",
	}); err.Error() != "halli" {
		t.Fatal(err)
	}
}

func TestSendMultiple(t *testing.T) {
	testSubjects := []userregistry.Message{
		{Name: "halli", Email: "halli@hallo.de", Entity: "1"},
		{Name: "hallo", Email: "hallo@halli.de", Entity: "2"},
		{Name: "halloechen", Email: "hallo@holli.at", Entity: "3"},
	}

	wg := &sync.WaitGroup{}

	for _, subject := range testSubjects {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := testClient.SendMessage(subject); err.Error() != subject.Name {
				log.Fatal(err)
			}
		}()
	}

	wg.Wait()
}

func TestSendOverflow(t *testing.T) {
	testSubjects := []userregistry.Message{
		{Name: "halli", Email: "halli@hallo.de", Entity: "1"},
		{Name: "hallo", Email: "hallo@halli.de", Entity: "2"},
		{Name: "halloechen", Email: "hallo@holli.at", Entity: "3"},
		{Name: "ciao", Email: "ciao@bella.it", Entity: "4"},
		{Name: "bon jour", Email: "bon@jour.fr", Entity: "5"},
		{Name: "goode tag", Email: "goode@tag.nl", Entity: "6"},
		{Name: "guude wie", Email: "guud@wie.de", Entity: "7"},
		{Name: "moin moin", Email: "moin@moin.de", Entity: "8"},
		{Name: "buenos diaz", Email: "buenoz@diaz.es", Entity: "9"},
		{Name: "karoffeln", Email: "kartoffeln@kartoffeln.de", Entity: "10"},
		{Name: "eier", Email: "eier@eier.de", Entity: "11"},
		{Name: "zwiebeln", Email: "zwiebeln@zwiebeln.at", Entity: "12"},
		{Name: "knoblauch", Email: "knoblauch@knoblauch.de", Entity: "13"},
		{Name: "speck", Email: "speck@speck.de", Entity: "14"},
		{Name: "salz", Email: "salz@salz.at", Entity: "15"},
		{Name: "pfeffer", Email: "pfeffer@pfeffer.de", Entity: "16"},
		{Name: "dat is aber lecker", Email: "dat.is@aber-lecker.de", Entity: "17"},
		{Name: "eric", Email: "eric@eric.at", Entity: "18"},
		{Name: "stan", Email: "stan@stan.de", Entity: "19"},
		{Name: "kyle", Email: "kyle@kyle.de", Entity: "20"},
		{Name: "kenny", Email: "kenny@kenny.at", Entity: "21"},
	}

	wg := &sync.WaitGroup{}

	for _, subject := range testSubjects {
		wg.Add(1)
		go func() {
			wg.Done()
			if err := testClient.SendMessage(subject); err.Error() != subject.Name {
				log.Fatal(err)
			}
		}()
	}

	wg.Wait()
}
