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
)

var jetStream nats.JetStreamContext = nil

type UserForAccounts struct {
	UserLogin string `json:"user_login" sql:"type:varchar(50);not null"`
	UserName  string `json:"user_name" sql:"type:varchar(50);not null"`
	Surname   string `json:"surname" sql:"type:varchar(50);not null"`
	IsDeleted bool   `json:"is_deleted" sql:"type:boolean;default false;not null"`
}

func JetStreamInit() error {
	if jetStream == nil {
		nc, err := nats.Connect("nats://host.docker.internal:4222")
		if err != nil {
			return err
		}

		jetStream, err = nc.JetStream(nats.PublishAsyncMaxPending(256))
		if err != nil {
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
		IsDeleted: user.IsDeleted,
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
