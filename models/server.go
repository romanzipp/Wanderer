package models

import (
	"github.com/hashicorp/nomad/api"
	"gorm.io/gorm"
	"log"
	"time"
)

type ServerStatus uint

const (
	StatusUp      ServerStatus = 0
	StatusErr     ServerStatus = 1
	StatusUnknown ServerStatus = 99999
)

func (status ServerStatus) GetTitle() string {
	switch status {
	case StatusUp:
		return "Up"
	case StatusErr:
		return "Error"
	case StatusUnknown:
		return "Unknown"
	}
	return ""
}

type Server struct {
	gorm.Model
	Name                 string
	Address              string
	CfAccessClientId     string
	CfAccessClientSecret string
	LastStatusCheck      time.Time
	Status               ServerStatus
}

func (s Server) NewNomadClient() (*api.Client, error) {
	headers := make(map[string][]string)
	if s.CfAccessClientId != "" {
		headers["CF-Access-Client-Id"] = []string{s.CfAccessClientId}
		headers["CF-Access-Client-Secret"] = []string{s.CfAccessClientSecret}
	}

	client, err := api.NewClient(&api.Config{
		Address: s.Address,
		Headers: headers,
	})

	if err != nil {
		log.Fatalf("error getting nomad client: %s", err)
		return nil, err
	}

	return client, nil
}

func (s Server) Check(db *gorm.DB) (ServerStatus, error) {
	client, err := s.NewNomadClient()
	if err != nil {
		return StatusUnknown, err
	}

	_, _, err = client.Jobs().List(&api.QueryOptions{})
	if err != nil {
		s.Status = StatusErr
	} else {
		s.Status = StatusUp
	}

	s.LastStatusCheck = time.Now()

	db.Save(s)

	return s.Status, nil
}
