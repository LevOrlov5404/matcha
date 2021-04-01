package service

import (
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	randomTokenLength     = 24
	randomTokenDigitsNum  = 8
	randomTokenSymbolsNum = 0

	emailConfirmationTokenPrefix       = "ec"
	passwordResetConfirmTokenKeyPrefix = "rpc"
)

type (
	VerificationService struct {
		log       *logrus.Entry
		repo      repository.VerificationCache
		generator RandomTokenGenerator
	}
)

func NewVerificationService(
	log *logrus.Entry, repo repository.VerificationCache, generator RandomTokenGenerator,
) *VerificationService {
	return &VerificationService{
		log:       log,
		repo:      repo,
		generator: generator,
	}
}

func (s *VerificationService) CreateEmailConfirmToken(userID uint64) (string, error) {
	token, err := s.generateRandomToken()
	if err != nil {
		return "", err
	}

	confirmToken := emailConfirmationTokenPrefix + token

	err = s.repo.PutEmailConfirmToken(userID, confirmToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to put email confirmation token to cache")
	}

	return confirmToken, nil
}

func (s *VerificationService) VerifyEmailConfirmToken(confirmToken string) (userID uint64, err error) {
	userID, err = s.repo.GetEmailConfirmTokenData(confirmToken)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get email confirmation token data from cache")
	}

	if err := s.repo.DeleteEmailConfirmToken(confirmToken); err != nil {
		s.log.Error(errors.Wrap(err, "failed to delete email confirmation token from cache"))
	}

	return userID, nil
}

func (s *VerificationService) CreatePasswordResetConfirmToken(userID uint64) (string, error) {
	token, err := s.generateRandomToken()
	if err != nil {
		return "", err
	}

	confirmToken := passwordResetConfirmTokenKeyPrefix + token

	err = s.repo.PutPasswordResetConfirmToken(userID, confirmToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to put reset password confirmation token to cache")
	}

	return confirmToken, nil
}

func (s *VerificationService) VerifyPasswordResetConfirmToken(confirmToken string) (userID uint64, err error) {
	userID, err = s.repo.GetPasswordResetConfirmTokenData(confirmToken)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get reset password confirmation token data from cache")
	}

	if err := s.repo.DeletePasswordResetConfirmToken(confirmToken); err != nil {
		s.log.Error(errors.Wrap(err, "failed to delete reset password confirmation token from cache"))
	}

	return userID, nil
}

func (s *VerificationService) generateRandomToken() (string, error) {
	randomToken, err := s.generator.Generate(
		randomTokenLength, randomTokenDigitsNum, randomTokenSymbolsNum, false, false,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random token")
	}

	return randomToken, nil
}
