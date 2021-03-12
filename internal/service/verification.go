package service

import (
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/pkg/errors"
)

const (
	randomTokenLength     = 24
	randomTokenNumDigits  = 8
	randomTokenNumSymbols = 0

	emailConfirmationTokenPrefix = "ec"
)

type (
	VerificationService struct {
		repo      repository.Cache
		generator RandomTokenGenerator
	}
)

func NewVerificationService(repo repository.Cache, generator RandomTokenGenerator) *VerificationService {
	return &VerificationService{
		repo:      repo,
		generator: generator,
	}
}

func (s *VerificationService) CreateEmailConfirmToken(userID uint64) (string, error) {
	// generate random confirmation token
	randomToken, err := s.generator.Generate(
		randomTokenLength, randomTokenNumDigits, randomTokenNumSymbols, false, false,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random token")
	}

	emailConfirmToken := emailConfirmationTokenPrefix + randomToken

	err = s.repo.PutEmailConfirmToken(userID, emailConfirmToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to put email confirmation token to cache")
	}

	return emailConfirmToken, nil
}

func (s *VerificationService) VerifyEmailConfirmToken(emailConfirmToken string) (userID uint64, err error) {
	userID, err = s.repo.GetEmailConfirmTokenData(emailConfirmToken)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get email confirmation token data from cache")
	}

	return userID, nil
}
