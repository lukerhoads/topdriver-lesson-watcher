package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InstructorId        string `yaml:"instructor_id"`
	Timings             string `yaml:"timings"`
	PickupLocation      string `yaml:"pickup_location"`
	DropoffLocation     string `yaml:"dropoff_location"`
	ReceiverPhoneNumber string `yaml:"receiver_phone_number"`
	PhoneNumber         string `yaml:"phone_number"`
}

func GetConfig(path string) (*Config, error) {
	// yamlFile, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	return nil, err
	// }

	// config := &Config{}
	// err = yaml.Unmarshal(yamlFile, config)
	// if err != nil {
	// 	return nil, err
	// }

	return &Config{
		InstructorId: "-1",
		Timings: "-1",
		PickupLocation: "Home",
		DropoffLocation: "Home",
		ReceiverPhoneNumber: "+13123994384",
		PhoneNumber: "+18484209423"
	}, nil
}
