package application

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"
)

// NewAuth constructor for Auth
func NewAuth(
	authRepo repositories.UserRepository,
	logger zerolog.Logger,
) *Auth {
	app := &Auth{
		logger:   logger,
		authRepo: authRepo,
	}
	return app
}

// Auth ....
type Auth struct {
	logger   zerolog.Logger
	authRepo repositories.UserRepository
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" bindinig:"required"`
	Password string `json:"password" binding:"required"`
}

type OTPInput struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

type UpdateAuthRequest struct {
	UserID    string `json:"user_id"`
	Extension string `json:"extension"`
	Number    string `json:"number"`
}

func (a *Auth) CreateUser(requestContext *context.RequestContext, request *RegisterUserInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.CreateUser")
	log.Debug().Msgf("creating new user")
	newUser := models.NewUserNoOtp(request.Name, request.Email, request.Password)
	log.Debug().Msgf("created new user %s", newUser.Name)
	user, err := a.authRepo.Persist(requestContext.TenantID(), newUser)
	if err != nil {
		return nil, exceptions.FailedPersisting(models.UserModelName, err)
	}
	log.Debug().Msgf("user %s persisted", user.Name)
	return user, nil
}

func (a *Auth) LoginUser(requestContext *context.RequestContext, request *LoginUserInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.LoginUser")
	log.Debug().Msgf("login new user")
	user, err := a.authRepo.GetByEmail(requestContext.TenantID(), request.Email)
	if err != nil {
		return nil, exceptions.FailedPersisting(models.UserModelName, err)
	}
	log.Debug().Msgf("user %s fetched", user.Name)
	return user, nil
}

func (a *Auth) GenerateOTP(requestContext *context.RequestContext, request *OTPInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.GenerateOTP")
	log.Debug().Msgf("login new user")
	user, err := a.authRepo.Get(requestContext.TenantID(), request.UserId)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.UserModelName, err)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "stem24-backend.com",
		AccountName: "admin@admin.com",
		SecretSize:  15,
	})
	if err != nil {
		return nil, exceptions.FailedFetchingServiceToken(err)
	}
	user.Otp_secret = key.Secret()
	user.Otp_auth_url = key.URL()
	user, err = a.authRepo.Update(requestContext.TenantID(), user)
	if err != nil {
		return nil, exceptions.FailedUpdating(models.UserModelName, err)
	}
	log.Debug().Msgf("user %s updated", user.Name)
	return user, nil
}

func (a *Auth) VerifyOTP(requestContext *context.RequestContext, request *OTPInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.GenerateOTP")
	message := "token is invalid or user doesn't exist"
	user, err := a.authRepo.Get(requestContext.TenantID(), request.UserId)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.UserModelName, err)
	}

	valid := totp.Validate(request.Token, user.Otp_secret)
	if !valid {
		return nil, exceptions.FailedQuerying(models.UserModelName, fmt.Errorf(message))
	}
	user.Otp_enabled = true
	user.Otp_verified = true
	user, err = a.authRepo.Update(requestContext.TenantID(), user)
	if err != nil {
		return nil, exceptions.FailedUpdating(models.UserModelName, err)
	}
	log.Debug().Msgf("user %s updated", user.Name)
	return user, nil
}

func (a *Auth) ValidateOTP(requestContext *context.RequestContext, request *OTPInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.GenerateOTP")
	message := "token is invalid or user doesn't exist"
	user, err := a.authRepo.Get(requestContext.TenantID(), request.UserId)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.UserModelName, err)
	}
	valid := totp.Validate(request.Token, user.Otp_secret)
	if !valid {
		return nil, exceptions.FailedQuerying(models.UserModelName, fmt.Errorf(message))
	}
	log.Debug().Msgf("user %s validated", user.Name)
	return user, nil
}

func (a *Auth) DisableOTP(requestContext *context.RequestContext, request *OTPInput) (*models.User, exceptions.ApplicationException) {
	log := requestContext.BuildLog(a.logger, "application.Auth.GenerateOTP")
	user, err := a.authRepo.Get(requestContext.TenantID(), request.UserId)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.UserModelName, err)
	}
	user.Otp_enabled = false
	user, err = a.authRepo.Update(requestContext.TenantID(), user)
	if err != nil {
		return nil, exceptions.FailedUpdating(models.UserModelName, err)
	}
	log.Debug().Msgf("user %s updated", user.Name)
	return user, nil
}
