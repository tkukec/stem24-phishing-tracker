package authentication

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/runtimebag"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// CustomClaims mapping for claims in jwt
type CustomClaims struct {
	Audience          []string       `json:"aud,omitempty"`
	ExpiresAt         int64          `json:"exp,omitempty"`
	ID                string         `json:"jti,omitempty"`
	IssuedAt          int64          `json:"iat,omitempty"`
	Issuer            string         `json:"iss,omitempty"`
	NotBefore         int64          `json:"nbf,omitempty"`
	Subject           string         `json:"sub,omitempty"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	UserID            string         `json:"user_id"`
	Name              string         `json:"name"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	PreferredUsername string         `json:"preferred_username"`
	ClientId          string         `json:"clientId"`
	Groups            []string       `json:"groups"`
}

func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.NotBefore, 0)), nil
}

func (c CustomClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c CustomClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	return c.Audience, nil
}

// ResourceAccess all role groupings
type ResourceAccess struct {
	RealmManagement RoleGroup `json:"realm-management"`
	Templating      RoleGroup `json:"templating"`
	Account         RoleGroup `json:"account"`
}

// RoleGroup a group of roles
type RoleGroup struct {
	Roles []string `json:"roles"`
}

// Valid check if jwt is valid, for now only check if aud contains app_name
func (c CustomClaims) Valid() error {
	var audFound bool
	appName := runtimebag.GetEnvString(constants.AppName, "")
	if appName == "" {
		return nil
	}
	for _, claim := range c.Audience {
		if claim == appName {
			audFound = true
		}
	}
	if !audFound {
		return fmt.Errorf("token missing aud: %s", appName)
	}
	return nil
}
