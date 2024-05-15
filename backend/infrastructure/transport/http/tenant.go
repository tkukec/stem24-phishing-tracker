package http

import (
	"github.com/asseco-voice/agent-management/application"
	"github.com/asseco-voice/agent-management/domain/models"
	"github.com/asseco-voice/agent-management/infrastructure/dto"
	"github.com/asseco-voice/agent-management/shared/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// Tenant ....
type Tenant struct {
	Controller
	tenantApp *application.Tenant
}

// NewTenant constructor for Tenant
func NewTenant(tenantApp *application.Tenant, controller Controller) *Tenant {
	return &Tenant{Controller: controller, tenantApp: tenantApp}
}

// Create
// @Summary Create new tenant
// @Description create new tenant
// @Tags Tenants
// @Accept  json
// @Produce  json
// @Param tenant body application.CreateTenantRequest true "Tenant"
// @Success 201 {object} dto.Tenant
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /statuses [post]
func (t *Tenant) Create(ctx *gin.Context) {
	requestContext, appErr := t.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.CreateTenantRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(t.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.SeedSignals = models.SignalSeed()
	request.ActivityStatus = models.ActivityStatusSeed()
	request.SeedSystemStatuses = models.SystemStatusSeed()
	request.SeedBasicStatuses = models.BasicStatusSeed()
	request.SeedGlobalStatuses = models.GlobalSeed()
	request.SeedChannels = models.ChannelSeed()
	tenant, appErr := t.tenantApp.Create(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewTenant(tenant))
}
