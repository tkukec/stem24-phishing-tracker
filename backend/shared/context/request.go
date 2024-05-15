package context

import (
	"context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/authentication"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"time"
)

type User struct {
	token    *jwt.Token
	tenantID string
}

func NewUser(token *jwt.Token, tenantID string) *User {
	return &User{token: token, tenantID: tenantID}
}

func (u *User) TenantID() string {
	return u.tenantID
}

func (u *User) ID() string {
	return u.token.Claims.(*authentication.CustomClaims).UserID
}

func (u *User) Name() string {
	return u.token.Claims.(*authentication.CustomClaims).Name
}

func (u *User) FamilyName() string {
	return u.token.Claims.(*authentication.CustomClaims).FamilyName
}

func (u *User) GivenName() string {
	return u.token.Claims.(*authentication.CustomClaims).GivenName
}

func (u *User) Token() *jwt.Token {
	return u.token
}

func (u *User) Claims() *authentication.CustomClaims {
	return u.token.Claims.(*authentication.CustomClaims)
}

type RequestContext struct {
	xCorrelationID string
	tenantID       string
	user           *User
	context        context.Context
}

func NewRequestContext(xCorrelationID string, tenantID string, user *User, context context.Context) *RequestContext {
	return &RequestContext{xCorrelationID: xCorrelationID, tenantID: tenantID, user: user, context: context}
}

func Background() *RequestContext {
	return &RequestContext{
		xCorrelationID: "",
		tenantID:       "",
		user:           nil,
		context:        context.Background(),
	}
}

func WithContext(ctx context.Context) *RequestContext {
	return &RequestContext{xCorrelationID: "", tenantID: "", user: nil, context: ctx}
}

func (r *RequestContext) XCorrelationID() string {
	return r.xCorrelationID
}

func (r *RequestContext) TenantID() string {
	return r.tenantID
}

func (r *RequestContext) User() *User {
	return r.user
}

func (r *RequestContext) Context() context.Context {
	return r.context
}

// BackgroundCopy returns a copy of context with a new context.Background() instead of the original one
func (r *RequestContext) BackgroundCopy() *RequestContext {
	return &RequestContext{
		xCorrelationID: r.XCorrelationID(),
		tenantID:       r.TenantID(),
		user:           r.User(),
		context:        context.Background(),
	}
}

// Copy returns a copy of context with a new context.Context
func (r *RequestContext) Copy(ctx context.Context) *RequestContext {
	return &RequestContext{
		xCorrelationID: r.XCorrelationID(),
		tenantID:       r.TenantID(),
		user:           r.User(),
		context:        ctx,
	}
}

func (r *RequestContext) WithTimeout(timeout time.Duration) (*RequestContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	return r.Copy(ctx), cancel
}

func (r *RequestContext) BuildLog(logger zerolog.Logger, method string, fields ...map[string]string) zerolog.Logger {
	log := logger.With().Str(constants.Method, method).
		Str(constants.XCorrelationID, r.XCorrelationID()).
		Str(constants.TenantIdentifier, r.TenantID()).Timestamp()

	if r.User() != nil {
		log = log.Str(constants.UserDisplayName, r.User().Name()).
			Str(constants.UserId, r.User().ID())
	}

	for _, fieldPair := range fields {
		for fieldName, filedValue := range fieldPair {
			log = log.Str(fieldName, filedValue)
		}
	}

	return log.Logger()
}
