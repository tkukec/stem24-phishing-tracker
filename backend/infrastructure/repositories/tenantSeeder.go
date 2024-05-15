package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"os"
)

type TenantSeeder struct {
	tenantRepo TenantRepository
	logger     zerolog.Logger
}

func NewTenantSeeder(
	tenantRepo TenantRepository,
	logger zerolog.Logger) *TenantSeeder {
	return &TenantSeeder{
		tenantRepo: tenantRepo,
		logger:     logger}
}

func (t *TenantSeeder) Description() string {
	return fmt.Sprintf("seeds data from file, if file does not exists seed will be for %s", constants.DefaultTenant)
}

func (t *TenantSeeder) Execute(location string) error {
	if location == "" {
		_, err := t.Run(context.Background(), &NewTenantRequest{
			Name: constants.DefaultTenant,
		})
		if err != nil {
			return err
		}
		return nil
	}

	jsonFile, err := os.Open(location)
	if err != nil {
		log.Panicln(err.Error())
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tenants []*models.Tenant
	err = json.Unmarshal(byteValue, &tenants)
	if err != nil {
		log.Panicln(err.Error())
	}

	for _, tenant := range tenants {
		_, err := t.Run(context.Background(), &NewTenantRequest{
			Name: tenant.Name,
		})
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

type NewTenantRequest struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}

func (t *TenantSeeder) Run(ctx *context.RequestContext, request *NewTenantRequest) (*models.Tenant, error) {
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
