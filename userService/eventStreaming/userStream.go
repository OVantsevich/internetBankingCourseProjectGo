package eventStreaming

import (
	"encoding/json"
	"github.com/OVantsevich/internetBankingCourseProjectGo/userService/domain"
	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
)

const (
	StreamName     = "USERS"
	StreamSubjects = "USERS.*"
)

const (
	SubjectNameUserCreated = "USERS.userCreated"
	SubjectNameUserDeleted = "USERS.userDeleted"
	SubjectNameUserUpdated = "USERS.userUpdated"
)

var jetStream nats.JetStreamContext = nil

type UserForAccounts struct {
	UserLogin string `json:"user_login" sql:"type:varchar(50);not null"`
	UserName  string `json:"user_name" sql:"type:varchar(50);not null"`
	Surname   string `json:"surname" sql:"type:varchar(50);not null"`
}

func JetStreamInit() error {
	if jetStream == nil {
		nc, err := nats.Connect("nats://host.docker.internal:4222")
		if err != nil {
			log.Errorf("jetstream init: %v", err)
			return err
		}

		jetStream, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
		if err != nil {
			log.Errorf("jetstream init: %v", err)
			return err
		}
	}

	err := CreateStream()
	if err != nil {
		return err
	}

	return nil
}

func CreateStream() error {
	stream, _ := jetStream.StreamInfo(StreamName)

	if stream == nil {
		log.Printf("creating stream: %v", StreamName)

		_, err := jetStream.AddStream(&nats.StreamConfig{
			Name:      StreamName,
			Subjects:  []string{StreamSubjects},
			Retention: nats.InterestPolicy,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreatingUser(user *domain.User) error {
	createdUser, err := json.Marshal(UserForAccounts{
		UserLogin: user.UserLogin,
		UserName:  user.UserName,
		Surname:   user.Surname,
	})
	if err != nil {
		log.Errorf("creating users: %v", err)
		return err
	}

	_, err = jetStream.Publish(SubjectNameUserCreated, createdUser)
	if err != nil {
		log.Errorf("creating users: %v", err)
		return err
	}
	return nil
}

func UpdatingUser(user *domain.User) error {
	updatedUser, err := json.Marshal(UserForAccounts{
		UserLogin: user.UserLogin,
		UserName:  user.UserName,
		Surname:   user.Surname,
	})
	if err != nil {
		log.Errorf("updating users: %v", err)
		return err
	}

	_, err = jetStream.Publish(SubjectNameUserUpdated, updatedUser)
	if err != nil {
		log.Errorf("updating users: %v", err)
		return err
	}
	return nil
}

func DeletingUser(login string) error {
	deletedUser, err := json.Marshal(UserForAccounts{
		UserLogin: login,
	})
	if err != nil {
		log.Errorf("deleting users: %v", err)
		return err
	}

	_, err = jetStream.Publish(SubjectNameUserDeleted, deletedUser)
	if err != nil {
		log.Errorf("deleting users: %v", err)
		return err
	}
	return nil
}
