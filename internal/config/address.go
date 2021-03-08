package config

import (
	"fmt"
	"net"
	"strconv"
)

type AddressConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (a *AddressConfig) String() string {
	return fmt.Sprintf("%v:%v", a.Host, a.Port)
}

func (a *AddressConfig) IsEmpty() bool {
	return a.Host == "" || a.Port == 0
}

func (a *AddressConfig) Decode(value string) error {
	host, portStr, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	a.Host = host
	a.Port = port

	return nil
}

func (a *AddressConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var input string
	if err := unmarshal(&input); err != nil {
		return err
	}

	return a.Decode(input)
}
