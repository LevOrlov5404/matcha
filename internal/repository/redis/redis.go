package redis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	sessionKeyPrefix                   = "session:"
	userToSessionKeyPrefix             = "uToSession:"
	accessTokenKeyPrefix               = "at:"
	userBlockingKeyPrefix              = "ub:"
	emailConfirmTokenKeyPrefix         = "eConf:"
	passwordResetConfirmTokenKeyPrefix = "rpConf:"
)

type (
	Options struct {
		AccessTokenLifetime               int
		RefreshTokenLifetime              int
		UserBlockingLifetime              int
		EmailConfirmTokenLifetime         int
		PasswordResetConfirmTokenLifetime int
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

func (r *Redis) PutSessionAndAccessToken(session models.Session, refreshToken string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return err
	}

	if err := conn.Send("MULTI"); err != nil {
		return err
	}

	if err = conn.Send("SETEX", sessionKeyPrefix+refreshToken,
		r.options.RefreshTokenLifetime, sessionBytes,
	); err != nil {
		return err
	}

	// link user to session to be able to throw user
	if err = conn.Send("SETEX", userToSessionKeyPrefix+session.UserID+":"+refreshToken,
		r.options.RefreshTokenLifetime, session.AccessTokenID,
	); err != nil {
		return err
	}

	if err = conn.Send("SETEX", accessTokenKeyPrefix+session.AccessTokenID,
		r.options.AccessTokenLifetime, refreshToken,
	); err != nil {
		return err
	}

	if _, err = conn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetSession(refreshToken string) (*models.Session, error) {
	conn, err := r.getConnect()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	resp, err := redis.String(conn.Do("GET", sessionKeyPrefix+refreshToken))
	if err != nil {
		return nil, err
	}

	session := &models.Session{}
	err = json.Unmarshal([]byte(resp), &session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Redis) DeleteSession(refreshToken string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if _, err = conn.Do("DEL", sessionKeyPrefix+refreshToken); err != nil {
		return err
	}

	return nil
}

func (r *Redis) DeleteUserToSession(userID, refreshToken string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if _, err = conn.Do("DEL", userToSessionKeyPrefix+userID+":"+refreshToken); err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetAccessTokenData(accessTokenID string) (refreshToken string, err error) {
	conn, err := r.getConnect()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	refreshToken, err = redis.String(conn.Do("GET", accessTokenKeyPrefix+accessTokenID))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (r *Redis) DeleteAccessToken(accessTokenID string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if _, err = conn.Do("DEL", accessTokenKeyPrefix+accessTokenID); err != nil {
		return err
	}

	return nil
}

func (r *Redis) AddUserBlocking(fingerprint string) (int64, error) {
	conn, err := r.getConnect()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	count, err := conn.Do("INCR", userBlockingKeyPrefix+fingerprint)
	if err != nil {
		return 0, err
	}

	_, err = conn.Do("EXPIRE", userBlockingKeyPrefix+fingerprint, r.options.UserBlockingLifetime)
	if err != nil {
		return 0, err
	}

	return count.(int64), nil
}

func (r *Redis) GetUserBlocking(fingerprint string) (int, error) {
	conn, err := r.getConnect()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	count, err := redis.Int(conn.Do("GET", userBlockingKeyPrefix+fingerprint))
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return 0, err
	}

	return count, nil
}

func (r *Redis) DeleteUserBlocking(fingerprint string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	_, err = conn.Do("DEL", userBlockingKeyPrefix+fingerprint)

	return err
}

func (r *Redis) PutEmailConfirmToken(userID uint64, token string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if err = conn.Send("SETEX", emailConfirmTokenKeyPrefix+token,
		r.options.EmailConfirmTokenLifetime, userID,
	); err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetEmailConfirmTokenData(token string) (userID uint64, err error) {
	conn, err := r.getConnect()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	userID, err = redis.Uint64(conn.Do("GET", emailConfirmTokenKeyPrefix+token))
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Redis) DeleteEmailConfirmToken(token string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if _, err = conn.Do("DEL", emailConfirmTokenKeyPrefix+token); err != nil {
		return err
	}

	return nil
}

func (r *Redis) PutPasswordResetConfirmToken(userID uint64, token string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if err = conn.Send("SETEX", passwordResetConfirmTokenKeyPrefix+token,
		r.options.PasswordResetConfirmTokenLifetime, userID,
	); err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetPasswordResetConfirmTokenData(token string) (userID uint64, err error) {
	conn, err := r.getConnect()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	userID, err = redis.Uint64(conn.Do("GET", passwordResetConfirmTokenKeyPrefix+token))
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Redis) DeletePasswordResetConfirmToken(token string) error {
	conn, err := r.getConnect()
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			r.log.Error(err)
		}
	}()

	if _, err = conn.Do("DEL", passwordResetConfirmTokenKeyPrefix+token); err != nil {
		return err
	}

	return nil
}
