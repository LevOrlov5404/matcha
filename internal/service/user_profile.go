package service

import (
	"bytes"
	"context"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/l-orlov/matcha/internal/config"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	matchaBucketName  = "matcha"
	pictureURLExpires = 3 * time.Hour
)

type (
	UserProfileService struct {
		log                *logrus.Entry
		maxUserPicturesNum int
		pathTemplates      config.FilePathTemplates
		repo               *repository.Repository
	}
)

func NewUserProfileService(
	log *logrus.Entry, maxUserPicturesNum int,
	pathTemplates config.FilePathTemplates, repo *repository.Repository,
) *UserProfileService {
	return &UserProfileService{
		log:                log,
		maxUserPicturesNum: maxUserPicturesNum,
		pathTemplates:      pathTemplates,
		repo:               repo,
	}
}

func (s *UserProfileService) GetUserProfileByID(ctx context.Context, id uint64) (*models.UserProfile, error) {
	profile, err := s.repo.User.GetUserProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if profile.AvatarPath != "" {
		profile.AvatarPath = s.getPictureURLByPath(ctx, profile.AvatarPath)
	}

	for i := 0; i < len(profile.PicturesPath); i++ {
		profile.PicturesPath[i] = s.getPictureURLByPath(ctx, profile.PicturesPath[i])
	}

	return profile, nil
}

func (s *UserProfileService) UpdateUserProfile(ctx context.Context, user models.UserProfile) error {
	return s.repo.User.UpdateUserProfile(ctx, user)
}

func (s *UserProfileService) UploadUserAvatar(ctx context.Context, userID uint64, file io.ReadSeeker) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ierrors.NewBusiness(ErrUserNotFound, "")
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(mime.String(), "image/") {
		return ierrors.NewBusiness(errors.Errorf("avatar file %s is not an image", mime.String()), "")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	path, err := prepareFilePath(s.pathTemplates.UserAvatar, map[string]interface{}{"UserID": userID})
	if err != nil {
		return err
	}

	if err := s.repo.Storage.PutFile(
		ctx, matchaBucketName, path, mime.String(), file,
	); err != nil {
		return err
	}

	if err := s.repo.User.UpdateUserAvatarPath(ctx, userID, path); err != nil {
		return err
	}

	return nil
}

func (s *UserProfileService) UploadUserPicture(ctx context.Context, userID uint64, file io.ReadSeeker) error {
	userProfile, err := s.repo.User.GetUserProfileByID(ctx, userID)
	if err != nil {
		return err
	}

	if userProfile == nil {
		return ierrors.NewBusiness(ErrUserNotFound, "")
	}

	if len(userProfile.PicturesPath) >= s.maxUserPicturesNum {
		return ierrors.NewBusiness(
			errors.Errorf("can not upload more than %d pictures", s.maxUserPicturesNum), "",
		)
	}

	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(mime.String(), "image/") {
		return ierrors.NewBusiness(
			errors.Errorf("picture file %s is not an image", mime.String()), "",
		)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	path, err := prepareFilePath(s.pathTemplates.UserPicture, map[string]interface{}{
		"UserID": userID,
		"Num":    len(userProfile.PicturesPath) + 1,
	})
	if err != nil {
		return err
	}

	if err := s.repo.Storage.PutFile(
		ctx, matchaBucketName, path, mime.String(), file,
	); err != nil {
		return err
	}

	userProfile.PicturesPath = append(userProfile.PicturesPath, path)
	if err := s.repo.User.UpdateUserPicturesPaths(ctx, userID, userProfile.PicturesPath); err != nil {
		return err
	}

	return nil
}

func prepareFilePath(pathTemplate string, pathParams map[string]interface{}) (string, error) {
	tpl, err := template.New("").Parse(pathTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, pathParams); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *UserProfileService) getPictureURLByPath(ctx context.Context, path string) string {
	url, err := s.repo.Storage.GetFileURL(ctx, matchaBucketName, path, pictureURLExpires)
	if err != nil {
		s.log.Errorf("failed to get picture URL by path %s: %v", path, err)
		return ""
	}

	return url
}
