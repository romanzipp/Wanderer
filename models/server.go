package models

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"gorm.io/gorm"
	"time"
)

type ServerStatus uint

const (
	StatusUp ServerStatus = 0
)

type Server struct {
	gorm.Model
	Name                 string
	Address              string
	Port                 string
	CfAccessClientId     string
	CfAccessClientSecret string
	LastStatusCheck      time.Time
	Status               ServerStatus
}

func (s Server) Check() {
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
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Println(client)

	res, _, err := client.Jobs().List(&api.QueryOptions{})
	if err != nil {
		fmt.Println("-----------err")
		fmt.Println(err)
	}
	fmt.Println(res)
}
