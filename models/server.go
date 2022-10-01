package models

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"gorm.io/gorm"
	"log"
	"time"
)

type ServerStatus uint

const (
	StatusUnknown ServerStatus = 0
	StatusUp      ServerStatus = 1
	StatusErr     ServerStatus = 2
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
	BasicAuthUser        string
	BasicAuthPassword    string
	LastStatusCheck      time.Time
	Status               ServerStatus
}

func (s Server) NewNomadClient() (*api.Client, error) {
	headers := make(map[string][]string)
	if s.CfAccessClientId != "" {
		headers["CF-Access-Client-Id"] = []string{s.CfAccessClientId}
		headers["CF-Access-Client-Secret"] = []string{s.CfAccessClientSecret}
	}

	if s.BasicAuthUser != "" && s.BasicAuthPassword != "" {
		headers["Authorization"] = []string{fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.BasicAuthUser, s.BasicAuthPassword))))}
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

	return s.Status, err
}
