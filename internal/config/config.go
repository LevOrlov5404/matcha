package config

import (
	"io/ioutil"

	"github.com/joeshaw/envdecode"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Address      AddressConfig `yaml:"address"`
		Logger       Logger        `yaml:"logger"`
		PostgresDB   PostgresDB    `yaml:"postgresDB"`
		Redis        Redis         `yaml:"redis"`
		JWT          JWT           `yaml:"jwt"`
		Cookie       Cookie        `yaml:"cookie"`
		UserBlocking UserBlocking  `yaml:"userBlocking"`
		Verification Verification  `yaml:"verification"`
		Mailer       Mailer        `yaml:"mailer"`
	}
	Logger struct {
		Level  string `yaml:"level" env:"LOGGER_LEVEL,default=info"`
		Format string `yaml:"format" env:"LOGGER_FORMAT,default=default"`
	}
	PostgresDB struct {
		Address         AddressConfig  `yaml:"address" env:"PG_ADDRESS,default=0.0.0.0:5432"`
		User            string         `yaml:"user" env:"PG_USER,default=postgres"`
		Password        string         `yaml:"password" env:"PG_PASSWORD,default=123"`
		Database        string         `yaml:"name" env:"PG_DATABASE,default=postgres"`
		ConnMaxLifetime DurationConfig `yaml:"connMaxLifetime"`
		MaxOpenConns    int            `yaml:"maxOpenConns"`
		MaxIdleConns    int            `yaml:"maxIdleConns"`
		Timeout         DurationConfig `yaml:"timeout"`
	}
	Redis struct {
		Address     AddressConfig  `yaml:"address" env:"REDIS_ADDRESS,default=0.0.0.0:6379"`
		Proto       string         `yaml:"proto"`
		MaxActive   int            `yaml:"maxActive"`
		MaxIdle     int            `yaml:"maxIdle"`
		IdleTimeout DurationConfig `yaml:"idleTimeout"`
	}
	JWT struct {
		AccessTokenLifetime  DurationConfig `yaml:"accessTokenLifetime"`
		RefreshTokenLifetime DurationConfig `yaml:"refreshTokenLifetime"`
		SigningKey           StdBase64      `yaml:"signingKey" env:"JWT_SIGNING_KEY,default=dGVzdA=="`
	}
	Cookie struct {
		HashKey  StdBase64 `yaml:"hashKey" env:"COOKIE_HASH_KEY,default=dGVzdA=="`
		BlockKey StdBase64 `yaml:"blockKey" env:"COOKIE_BLOCK_KEY,default=dGVzdA=="`
	}
	UserBlocking struct {
		Lifetime  DurationConfig `yaml:"lifetime"`
		MaxErrors int            `yaml:"maxErrors"`
	}
	Verification struct {
		EmailConfirmTokenLifetime         DurationConfig `yaml:"emailConfirmTokenLifetime"`
		PasswordResetConfirmTokenLifetime DurationConfig `yaml:"passwordResetConfirmTokenLifetime"`
	}
	Mailer struct {
		ServerAddress     AddressConfig  `yaml:"serverAddress" env:"EMAIL_SERVER_ADDRESS,default=smtp.gmail.com:587"`
		Username          string         `yaml:"username" env:"EMAIL_USERNAME,default=test"`
		Password          string         `yaml:"password" env:"EMAIL_PASSWORD,default=test"`
		Timeout           DurationConfig `yaml:"timeout"`
		MsgToSendChanSize int            `yaml:"msgToSendChanSize"`
		WorkersNum        int            `yaml:"workersNum"`
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
