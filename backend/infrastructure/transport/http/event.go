package http

import (
	"github.com/andrezz-b/stem24-phishing-tracker/application"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/dto"
	helpers "github.com/andrezz-b/stem24-phishing-tracker/shared"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// NewEvent constructor for Event
func NewEvent(eventApplication *application.Event, controller Controller) *Event {
	return &Event{
		Controller:       controller,
		eventApplication: eventApplication,
	}
}

// Event ....
type Event struct {
	Controller
	eventApplication *application.Event
}

// Create
// @Summary Create new Event
// @Description create new event
// @Tags Events
// @Accept  json
// @Produce  json
// @Param event body application.CreateEventRequest true "Event"
// @Success 201 {object} dto.Event
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /events [post]
func (c *Event) Create(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.CreateEventRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	event, appErr := c.eventApplication.Create(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewEvent(event))
}

// Get
// @Summary Get a event
// @Description get a event by id
// @Tags Events
// @Produce  json
// @Success 200 {object} dto.Event
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Event id"
// @Router /events/{id} [get]
func (c *Event) Get(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	event, appErr := c.eventApplication.Get(requestContext, ctx.Param("id"))
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewEvent(event))
}

// GetAll
// @Summary Get all event
// @Description get all event by id
// @Tags Events
// @Produce  json
// @Success 200 {array} dto.Event
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /events [get]
func (c *Event) GetAll(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	var requestParams database.GetAllEventsRequest
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	records, appErr := c.eventApplication.GetAll(requestContext, requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	Events := make([]*dto.Event, 0)
	for _, v := range records {
		Events = append(Events, dto.NewEvent(v))
	}
	ctx.JSON(http.StatusOK, Events)
}

// Delete
// @Summary Delete event
// @Description Delete event by id
// @Tags Events
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Event id"
// @Router /events/{id} [delete]
func (c *Event) Delete(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	appErr = c.eventApplication.Delete(requestContext, ctx.Param("id"))
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, helpers.ToJson("deleted"))
}

// Update
// @Summary Update a event
// @Description Update a event by id
// @Tags Events
// @Produce  json
// @Param event body application.UpdateEventRequest true "Event"
// @Success 200 {object} dto.Event
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Event id"
// @Router /events/{id} [put]
func (c *Event) Update(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.ID = ctx.Param("id")

	event, appErr := c.eventApplication.Update(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewEvent(event))
}
