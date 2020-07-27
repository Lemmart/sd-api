package main

import (
	"encoding/json"
	"io/ioutil"
)

type ServiceConfig struct {
	ListenAddress string `json:"listen_address"`
	UseTLS        bool   `json:"use_tls"`
	CertFile      string `json:"cert_file"`
	KeyFile       string `json:"key_file"`
}

type SdsConfig struct {
	StoreWebsites []string `json:"stores"`

	ServiceConfig *ServiceConfig `json:"service_config"`
}

type SdsRequest struct {
	Websites []string `json:"websites"`
}

type Offer struct {
	Amount   int
	Code     string
	Category string
}

type Company struct {
	Name    string
	Website string
	Offers  []*Offer
}

type SdsResponse struct {
	Companies []*Company
}

func LoadConfig(fileName string) (*SdsConfig, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &SdsConfig{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
