package cache_redis

import (
	"time"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

const (
	emailConfirmationTokenKeyPrefix = "eConf:"
)

type (
	Options struct {
		EmailConfirmTokenLifetime int
	}
	Redis struct {
		log     *logrus.Entry
		options Options
		pool    *redis.Pool
	}
)

func New(cfg config.Redis, log *logrus.Entry, options Options) *Redis {
	r := &Redis{
		log:     log,
		options: options,
	}

	r.pool = &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout.Duration(),
		Dial: func() (redis.Conn, error) {
			return redis.Dial(cfg.Proto, cfg.Address.String())
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")

			return err
		},
	}

	return r
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) getConnect() (redis.Conn, error) {
	c := r.pool.Get()
	if err := c.Err(); err != nil {
		return nil, err
	}

	return c, nil
}

func (r *Redis) PutEmailConfirmToken(clientID uint64, token string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if err = conn.Send("SETEX", emailConfirmationTokenKeyPrefix+token,
		r.options.EmailConfirmTokenLifetime, clientID,
	); err != nil {
		return err
	}

	return nil
}

//func (r *Redis) GetSession(refreshToken string) (*models.Session, error) {
//	conn, err := r.getConnect()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err := conn.Close(); err != nil {
//			r.log.Error(err)
//		}
//	}()
//
//	resp, err := redis.String(conn.Do("GET", sessionKeyPrefix+refreshToken))
//	if err != nil {
//		return nil, err
//	}
//
//	session := &models.Session{}
//	err = json.Unmarshal([]byte(resp), &session)
//	if err != nil {
//		return nil, err
//	}
//
//	return session, nil
//}

//func (r *Redis) DeleteSession(refreshToken string) error {
//	conn, err := r.getConnect()
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err := conn.Close(); err != nil {
//			r.log.Error(err)
//		}
//	}()
//
//	if _, err = conn.Do("DEL", sessionKeyPrefix+refreshToken); err != nil {
//		return err
//	}
//
//	return nil
//}
