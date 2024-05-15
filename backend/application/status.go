package application

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/rs/zerolog"
)

type Status struct {
	statusRepo repositories.StatusRepository
	logger     zerolog.Logger
}

func NewStatus(statusRepo repositories.StatusRepository,
	logger zerolog.Logger,
) *Status {
	return &Status{statusRepo: statusRepo, logger: logger}
}

type CreateStatusRequest struct {
	Name    string // Name of the status
	EventID string
}

func (a *Status) Create(ctx *context.RequestContext, request *CreateStatusRequest) (*models.Status, exceptions.ApplicationException) {
	status := &models.Status{
		Name: request.Name,
	}

	status, err := a.statusRepo.Persist(ctx.TenantID(), status)
	if err != nil {
		return nil, exceptions.FailedPersisting(models.StatusModelName, err)
	}
	return status, nil
}

type UpdateStatusRequest struct {
	Name string // Updated description of the status
	ID   string // Updated description of the status
}

func (request *UpdateStatusRequest) ApplyValues(status *models.Status) *models.Status {
	if request.Name != "" {
		status.Name = request.Name
	}
	// Add similar checks for other fields that can be updated

	return status
}
func (a *Status) Update(ctx *context.RequestContext, request *UpdateStatusRequest) (*models.Status, exceptions.ApplicationException) {
	status, err := a.statusRepo.Get(ctx.TenantID(), request.ID)
	if err != nil {
		return nil, exceptions.StatusNotFound(err)
	}
	status, err = a.statusRepo.Update(ctx.TenantID(), request.ApplyValues(status))
	if err != nil {
		return nil, exceptions.FailedUpdating(models.StatusModelName, err)
	}
	return status, nil
}

func (a *Status) Delete(ctx *context.RequestContext, ID string) exceptions.ApplicationException {
	status, err := a.statusRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return exceptions.StatusNotFound(err)
	}
	if err = a.statusRepo.Delete(ctx.TenantID(), status); err != nil {
		return exceptions.FailedDeleting(models.TenantModelName, err)
	}
	return nil
}

func (a *Status) Get(ctx *context.RequestContext, ID string) (*models.Status, exceptions.ApplicationException) {
	status, err := a.statusRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return nil, exceptions.StatusNotFound(err)
	}
	return status, nil
}

func (a *Status) GetAll(ctx *context.RequestContext, request database.GetAllStatusesRequest) ([]*models.Status, exceptions.ApplicationException) {
	statuss, err := a.statusRepo.GetAll(ctx.TenantID(), request)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.StatusModelName, err)
	}
	return statuss, nil
}
