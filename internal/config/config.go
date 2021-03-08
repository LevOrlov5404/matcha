package config

import (
	"io/ioutil"

	"github.com/joeshaw/envdecode"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Address          AddressConfig `yaml:"address"`
		Logger           Logger        `yaml:"logger"`
		DB               PostgresDB    `yaml:"db"`
		JWT              JWT           `yaml:"jwt"`
		UserPasswordSalt string        `yaml:"userPasswordSalt" env:"USER_PASSWORD_SALT,default=test"`
	}
	Logger struct {
		Level  string `yaml:"level" env:"LOGGER_LEVEL,default=info"`
		Format string `yaml:"format" env:"LOGGER_FORMAT,default=default"`
	}
	PostgresDB struct {
		Address         AddressConfig  `yaml:"host" env:"PG_ADDRESS,default=0.0.0.0:5432"`
		User            string         `yaml:"user" env:"PG_USER,default=postgres"`
		Password        string         `yaml:"password" env:"PG_PASSWORD,default=123"`
		Database        string         `yaml:"name" env:"PG_DATABASE,default=postgres"`
		ConnMaxLifetime DurationConfig `yaml:"connMaxLifetime"`
		MaxOpenConns    int            `yaml:"maxOpenConns"`
		MaxIdleConns    int            `yaml:"maxIdleConns"`
		Timeout         DurationConfig `yaml:"timeout"`
	}
	JWT struct {
		TokenLifetime DurationConfig `yaml:"tokenLifetime"`
		SigningKey    string         `yaml:"signingKey" env:"JWT_SIGNING_KEY,default=test"`
	}
)

func DecodeYamlFile(path string, v interface{}) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf, v)
}

func ReadFromFileAndSetEnv(path string, v interface{}) error {
	if err := DecodeYamlFile(path, v); err != nil {
		return err
	}

	return envdecode.Decode(v)
}
