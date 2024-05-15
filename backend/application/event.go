package application

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/rs/zerolog"
	"time"
)

type Event struct {
	eventRepo repositories.EventRepository
	logger    zerolog.Logger
}

func NewEvent(eventRepo repositories.EventRepository,
	logger zerolog.Logger,
) *Event {
	return &Event{eventRepo: eventRepo, logger: logger}
}

type CreateEventRequest struct {
	Name             string    // Name of the event
	Date             time.Time // Date of the event
	Brand            string    // Brand associated with the event
	Description      string    // Description of the event
	MalURL           string    // Malicious URL associated with the event
	MalDomainRegDate time.Time // Registration date of the malicious domain
	DNSRecord        string    // DNS record associated with the event
	Keywords         []string  // Keywords related to the event
}

func (a *Event) Create(ctx *context.RequestContext, request *CreateEventRequest) (*models.Event, exceptions.ApplicationException) {
	event := &models.Event{
		Name:             request.Name,
		Date:             request.Date,
		Brand:            request.Brand,
		Description:      request.Description,
		MalURL:           request.MalURL,
		MalDomainRegDate: request.MalDomainRegDate,
		DNSRecord:        request.DNSRecord,
		Keywords:         request.Keywords,
	}

	event, err := a.eventRepo.Persist(ctx.TenantID(), event)
	if err != nil {
		return nil, exceptions.FailedPersisting(models.EventModelName, err)
	}
	return event, nil
}

type UpdateEventRequest struct {
	ID               string    // ID of the event to be updated
	Name             string    // Updated name of the event
	Date             time.Time // Updated date of the event
	Brand            string    // Updated brand associated with the event
	Description      string    // Updated description of the event
	MalURL           string    // Updated malicious URL associated with the event
	MalDomainRegDate time.Time // Updated registration date of the malicious domain
	DNSRecord        string    // Updated DNS record associated with the event
	Keywords         []string  // Updated keywords related to the event
	// Add other fields as needed
}

func (request *UpdateEventRequest) ApplyValues(event *models.Event) *models.Event {
	if request.Name != "" {
		event.Name = request.Name
	}
	if !request.Date.IsZero() {
		event.Date = request.Date
	}
	if request.Brand != "" {
		event.Brand = request.Brand
	}
	if request.Description != "" {
		event.Description = request.Description
	}
	if request.MalURL != "" {
		event.MalURL = request.MalURL
	}
	if !request.MalDomainRegDate.IsZero() {
		event.MalDomainRegDate = request.MalDomainRegDate
	}
	if request.DNSRecord != "" {
		event.DNSRecord = request.DNSRecord
	}
	if len(request.Keywords) > 0 {
		event.Keywords = request.Keywords
	}
	// Add similar checks for other fields that can be updated

	return event
}
func (a *Event) Update(ctx *context.RequestContext, request *UpdateEventRequest) (*models.Event, exceptions.ApplicationException) {
	event, err := a.eventRepo.Get(ctx.TenantID(), request.ID)
	if err != nil {
		return nil, exceptions.EventNotFound(err)
	}
	event, err = a.eventRepo.Update(ctx.TenantID(), request.ApplyValues(event))
	if err != nil {
		return nil, exceptions.FailedUpdating(models.EventModelName, err)
	}
	return event, nil
}

func (a *Event) Delete(ctx *context.RequestContext, ID string) exceptions.ApplicationException {
	event, err := a.eventRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return exceptions.EventNotFound(err)
	}
	if err = a.eventRepo.Delete(ctx.TenantID(), event); err != nil {
		return exceptions.FailedDeleting(models.TenantModelName, err)
	}
	return nil
}

func (a *Event) Get(ctx *context.RequestContext, ID string) (*models.Event, exceptions.ApplicationException) {
	event, err := a.eventRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return nil, exceptions.EventNotFound(err)
	}
	return event, nil
}

func (a *Event) GetAll(ctx *context.RequestContext, request database.GetAllEventsRequest) ([]*models.Event, exceptions.ApplicationException) {
	events, err := a.eventRepo.GetAll(ctx.TenantID(), request)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.EventModelName, err)
	}
	return events, nil
}
