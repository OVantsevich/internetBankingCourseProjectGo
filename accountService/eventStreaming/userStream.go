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
	SubjectNameUserDeleted = "USERS.userDeleted"
	SubjectNameUserUpdated = "USERS.userUpdated"
)

var jetStream nats.JetStreamContext = nil

func JetStreamInit() {
	if jetStream == nil {
		nc, err := nats.Connect(domain.Config.NatsUrl)
		if err != nil {
			log.Fatalf("jetstream init: %v", err)
		}

		jetStream, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
		if err != nil {
			log.Fatalf("jetstream init: %v", err)
		}
	}
}

func HandleUser() {
	_, err := jetStream.Subscribe(SubjectNameUserCreated, CreatingUserHandler)
	if err != nil {
		log.Fatalf("Subscribe for " + SubjectNameUserCreated + " failed")
	}
	_, err = jetStream.Subscribe(SubjectNameUserUpdated, UpdatingUserHandler)
	if err != nil {
		log.Fatalf("Subscribe for " + SubjectNameUserUpdated + " failed")
	}
	_, err = jetStream.Subscribe(SubjectNameUserDeleted, DeletingUserHandler)
	if err != nil {
		log.Fatalf("Subscribe for " + SubjectNameUserDeleted + " failed")
	}
}

func CreatingUserHandler(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		log.Errorf("Unable to Ack", err)
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
		return
	}

	log.Printf("Consumer  =>  Subject: %s  -  user login:%s", msg.Subject, user.UserLogin)
}

func UpdatingUserHandler(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		log.Errorf("Unable to Ack", err)
		return
	}

	var user domain.User
	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		log.Fatal(err)
	}

	currentUser, _, err := repository.GetUserByLogin(context.Background(), user.UserLogin)
	if err != nil {
		log.Errorf("error updating user "+user.UserLogin+", user doesn't exist: %v", err)
		return
	}
	currentUser.UpdateUser(&user)

	err = repository.UpdateUser(context.Background(), currentUser)
	if err != nil {
		log.Errorf("error with updating user: %v", err)
		return
	}

	log.Printf("Consumer  =>  Subject: %s  -  user login:%s  -  \n", msg.Subject, user.UserLogin)
}

func DeletingUserHandler(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		log.Errorf("Unable to Ack", err)
		return
	}

	var user domain.User
	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		log.Fatal(err)
	}

	err = repository.DeleteUser(context.Background(), user.UserLogin)
	if err != nil {
		log.Errorf("error with deleting user: %v", err)
	}

	log.Printf("Consumer  =>  Subject: %s  -  user login:%s  -  \n", msg.Subject, user.UserLogin)
}
