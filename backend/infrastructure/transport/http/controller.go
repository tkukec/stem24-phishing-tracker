package http

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type Controller struct {
	tenantRepo repositories.TenantRepository
}

func NewController(tenantRepo repositories.TenantRepository) Controller {
	return Controller{tenantRepo: tenantRepo}
}

func (c *Controller) BuildRequestContext(ctx *gin.Context) (*context.RequestContext, exceptions.ApplicationException) {
	xCorrelation, ok := ctx.Get(constants.XCorrelationID)
	if !ok {
		data := map[string][]string{
			constants.XCorrelationID: {
				fmt.Sprintf("%s is required", constants.XCorrelationID),
			},
		}
		return nil, exceptions.BadRequest(data, "")
	}
	tenantName, ok := ctx.Get(constants.TenantIdentifier)
	if !ok {
		data := map[string][]string{
			constants.TenantIdentifier: {
				fmt.Sprintf("%s is required", constants.TenantIdentifier),
			},
		}
		return nil, exceptions.BadRequest(data, "")
	}
	tenant, err := c.tenantRepo.GetByName(tenantName.(string))
	if err != nil {
		data := map[string][]string{
			constants.TenantIdentifier: {
				err.Error(),
			},
		}
		return nil, exceptions.BadRequest(data, "")
	}
	token, ok := ctx.Get("user")
	if !ok {
		data := map[string][]string{
			"user": {
				"missing",
			},
		}
		return nil, exceptions.BadRequest(data, "")
	}
	return context.NewRequestContext(xCorrelation.(string), tenant.ID, context.NewUser(token.(*jwt.Token), tenant.ID), ctx.Copy()), nil
}

func (c *Controller) ValidationErrors(validationErrors validator.ValidationErrors) map[string][]string {
	errors := make(map[string][]string)

	for _, fieldErr := range validationErrors {
		errors[fieldErr.Field()] = []string{
			fmt.Sprintf("%s is %s", fieldErr.Field(), fieldErr.Tag()),
		}
	}
	return errors
}
