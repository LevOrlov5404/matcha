package config

import (
	cr "github.com/l-orlov/task-tracker/pkg/configreader"
)

type (
	Config struct {
		Port               string            `yaml:"port" env:"PORT,default=8080"`
		Logger             Logger            `yaml:"logger"`
		PostgresDB         PostgresDB        `yaml:"postgresDB"`
		Redis              Redis             `yaml:"redis"`
		JWT                JWT               `yaml:"jwt"`
		Cookie             Cookie            `yaml:"cookie"`
		UserBlocking       UserBlocking      `yaml:"userBlocking"`
		Verification       Verification      `yaml:"verification"`
		Mailer             Mailer            `yaml:"mailer"`
		Minio              Minio             `yaml:"minio"`
		FilePathTemplates  FilePathTemplates `yaml:"filePathTemplates"`
		MaxUserPicturesNum int               `yaml:"maxUserPicturesNum"`
	}
	Logger struct {
		Level  string `yaml:"level" env:"LOGGER_LEVEL,default=info"`
		Format string `yaml:"format" env:"LOGGER_FORMAT,default=default"`
	}
	PostgresDB struct {
		Address         cr.AddressConfig  `yaml:"address" env:"PG_ADDRESS,default=0.0.0.0:5432"`
		User            string            `yaml:"user" env:"PG_USER,default=postgres"`
		Password        string            `yaml:"password" env:"PG_PASSWORD,default=123"`
		Database        string            `yaml:"name" env:"PG_DATABASE,default=postgres"`
		ConnMaxLifetime cr.DurationConfig `yaml:"connMaxLifetime"`
		MaxOpenConns    int               `yaml:"maxOpenConns"`
		MaxIdleConns    int               `yaml:"maxIdleConns"`
		Timeout         cr.DurationConfig `yaml:"timeout"`
	}
	Redis struct {
		Address     cr.AddressConfig  `yaml:"address" env:"REDIS_ADDRESS,default=0.0.0.0:6379"`
		Proto       string            `yaml:"proto"`
		MaxActive   int               `yaml:"maxActive"`
		MaxIdle     int               `yaml:"maxIdle"`
		IdleTimeout cr.DurationConfig `yaml:"idleTimeout"`
	}
	JWT struct {
		AccessTokenLifetime  cr.DurationConfig `yaml:"accessTokenLifetime"`
		RefreshTokenLifetime cr.DurationConfig `yaml:"refreshTokenLifetime"`
		SigningKey           cr.StdBase64      `yaml:"signingKey" env:"JWT_SIGNING_KEY,default=dGVzdA=="`
	}
	Cookie struct {
		HashKey  cr.StdBase64 `yaml:"hashKey" env:"COOKIE_HASH_KEY,default=dGVzdA=="`
		BlockKey cr.StdBase64 `yaml:"blockKey" env:"COOKIE_BLOCK_KEY,default=dGVzdA=="`
		Domain   string       `yaml:"domain" env:"COOKIE_DOMAIN"`
	}
	UserBlocking struct {
		Lifetime  cr.DurationConfig `yaml:"lifetime"`
		MaxErrors int               `yaml:"maxErrors"`
	}
	Verification struct {
		EmailConfirmTokenLifetime         cr.DurationConfig `yaml:"emailConfirmTokenLifetime"`
		PasswordResetConfirmTokenLifetime cr.DurationConfig `yaml:"passwordResetConfirmTokenLifetime"`
	}
	Mailer struct {
		ServerAddress     cr.AddressConfig  `yaml:"serverAddress" env:"EMAIL_SERVER_ADDRESS,default=smtp.gmail.com:587"`
		Username          string            `yaml:"username" env:"EMAIL_USERNAME,default=test"`
		Password          string            `yaml:"password" env:"EMAIL_PASSWORD,default=test"`
		Timeout           cr.DurationConfig `yaml:"timeout"`
		MsgToSendChanSize int               `yaml:"msgToSendChanSize"`
		WorkersNum        int               `yaml:"workersNum"`
	}
	Minio struct {
		Endpoint  cr.AddressConfig  `yaml:"endpoint" env:"MINIO_ENDPOINT,default=0.0.0.0:9000"`
		AccessKey string            `yaml:"accessKey" env:"MINIO_ACCESS_KEY,default=minio"`
		SecretKey string            `yaml:"secretKey" env:"MINIO_SECRET_KEY,default=minio123"`
		UseSSL    bool              `yaml:"useSSL"`
		Timeout   cr.DurationConfig `yaml:"timeout"`
	}
	FilePathTemplates struct {
		UserAvatar  string `yaml:"userAvatar"`
		UserPicture string `yaml:"userPicture"`
	}
)

func Init(path string) (*Config, error) {
	var cfg Config
	if err := cr.ReadYamlAndSetEnv(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
