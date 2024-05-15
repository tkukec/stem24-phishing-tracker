package tenant

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/rs/zerolog"
)

type Tenant struct {
	tenantRepo repositories.TenantRepository
	logger     zerolog.Logger
}

func NewNewTenant(
	tenantRepo repositories.TenantRepository,
	logger zerolog.Logger,
) *Tenant {
	return &Tenant{
		tenantRepo: tenantRepo,
		logger:     logger,
	}
}

type NewTenantRequest struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}

func (t *Tenant) Execute(ctx *context.RequestContext, request *NewTenantRequest) (*models.Tenant, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.ChangeGlobalStatus")

	log.Debug().Msgf("Creating new tenet %s", request.Name)
	if tenant, err := t.tenantRepo.GetByName(request.Name); err == nil && tenant != nil {
		log.Debug().Msgf("tenant %s (%s) exists, skipping creation....", request.Name, tenant.ID)
		return tenant, nil
	}

	tenant, err := t.tenantRepo.Persist(&models.Tenant{
		ID:   request.ID,
		Name: request.Name,
	})
	if err != nil {
		log.Debug().Msgf("failed creating new tenant %s with error %s", request.Name, err.Error())
		return nil, fmt.Errorf("failed creating tenant %s with error %s", request.Name, err.Error())
	}
	log.Debug().Msgf("tenet %s created. Starting seed procedure for new tenant....", tenant.ID)

	return tenant, nil
}
