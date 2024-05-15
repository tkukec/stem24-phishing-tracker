package application

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	tenantService "github.com/andrezz-b/stem24-phishing-tracker/domain/services/tenant"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/rs/zerolog"
)

type Tenant struct {
	createTenantService *tenantService.Tenant
	tenantRepo          repositories.TenantRepository
	logger              zerolog.Logger
}

func NewTenant(createTenantService *tenantService.Tenant, tenantRepo repositories.TenantRepository, logger zerolog.Logger) *Tenant {
	return &Tenant{createTenantService: createTenantService, tenantRepo: tenantRepo, logger: logger}
}

type CreateTenantRequest struct {
	ID   string `json:"-"`
	Name string `json:"name" binding:"required"`
}

func (t *Tenant) HandleIamMessage(ctx *context.RequestContext, msg *amqp.EventMessage) {
	log := ctx.BuildLog(t.logger, "application.Tenant.HandleIamMessage")
	log.Debug().Msg(shared.ToJsonString(msg))
}

func (t *CreateTenantRequest) toServiceRequest() *tenantService.NewTenantRequest {
	return &tenantService.NewTenantRequest{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (t *Tenant) Create(ctx *context.RequestContext, request *CreateTenantRequest) (*models.Tenant, exceptions.ApplicationException) {
	tenant, err := t.createTenantService.Execute(ctx, request.toServiceRequest())
	if err != nil {
		return nil, exceptions.Internal("failed creating tenant", map[string][]string{
			"tenant": {
				err.Error(),
			},
		}, "")
	}
	return tenant, nil
}

func (t *Tenant) Delete(ctx *context.RequestContext, tenantID string) exceptions.ApplicationException {
	tenant, err := t.tenantRepo.Get(tenantID)
	if err != nil {
		return exceptions.TenantNotFound(err)
	}
	err = t.tenantRepo.Delete(tenant)
	if err != nil {
		return exceptions.FailedDeleting(models.TenantModelName, err)
	}
	return nil
}
