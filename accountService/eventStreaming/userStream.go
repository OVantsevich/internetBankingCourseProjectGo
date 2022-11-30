package eventStreaming

import (
	"context"
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/domain"
	"github.com/OVantsevich/internetBankingCourseProjectGo/accountService/repository"
	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
)

const (
	SubjectNameUserCreated = "USERS.userCreated"
)

var JetStream nats.JetStreamContext = nil

func JetStreamInit() error {
	if JetStream == nil {
		nc, err := nats.Connect("nats://host.docker.internal:4222")
		if err != nil {
			return err
		}

		JetStream, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatingUserHandler(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		log.Printf("Unable to Ack", err)
		return
	}

	var user domain.User
	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		log.Fatal(err)
	}

	err = repository.CreateUser(context.Background(), &user)
	if err != nil {
		log.Errorf("error with creating user: %v", err)
	}

	log.Printf("Consumer  =>  Subject: %s  -  user login:%s  -  \n", msg.Subject, user.UserLogin)
}
