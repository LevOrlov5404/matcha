package service

import (
	"context"

	"github.com/l-orlov/matcha/internal/config"
	iErrs "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	AuthenticationService struct {
		cfg  *config.Config
		log  *logrus.Entry
		repo *repository.Repository
	}
)

const (
	ErrMsgWrongPassword  = "wrong password"
	ErrMsgBlockedByLimit = "user is blocked due to exceeding the error limit"
)

func NewAuthenticationService(
	cfg *config.Config, log *logrus.Entry, repo *repository.Repository,
) *AuthenticationService {
	return &AuthenticationService{
		cfg:  cfg,
		log:  log,
		repo: repo,
	}
}

func (s *AuthenticationService) AuthenticateUserByUsername(
	ctx context.Context, username, password, fingerprint string,
) (userID uint64, err error) {
	if err := s.checkUserBlocking(fingerprint); err != nil {
		return 0, err
	}

	user, err := s.repo.User.GetUserByUsername(ctx, username)
	if err != nil {
		return 0, err
	}

	if user == nil {
		return 0, iErrs.NewBusiness(errors.New("user does not exist"), "")
	}

	if err := s.checkUserPasswordHash(fingerprint, user.Password, password); err != nil {
		return 0, err
	}

	if err := s.repo.SessionCache.DeleteUserBlocking(fingerprint); err != nil {
		s.log.Errorf("err while DeleteUserBlocking: %v", err)
	}

	return user.ID, nil
}

func (s *AuthenticationService) checkUserBlocking(fingerprint string) error {
	count, err := s.repo.SessionCache.GetUserBlocking(fingerprint)
	if err != nil {
		s.log.Errorf("err while GetUserBlocking: %v", err)
	}

	if count >= s.cfg.UserBlocking.MaxErrors {
		return errors.New(ErrMsgBlockedByLimit)
	}

	return nil
}

func (s *AuthenticationService) checkUserPasswordHash(fingerprint, hash, password string) error {
	if !models.CheckPasswordHash(hash, password) {
		if _, err := s.repo.SessionCache.AddUserBlocking(fingerprint); err != nil {
			s.log.Errorf("err while AddUserBlocking: %v", err)
		}

		return errors.New(ErrMsgWrongPassword)
	}

	return nil
}
