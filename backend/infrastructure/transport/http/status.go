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

// NewStatus constructor for Status
func NewStatus(statusApplication *application.Status, controller Controller) *Status {
	return &Status{
		Controller:        controller,
		statusApplication: statusApplication,
	}
}

// Status ....
type Status struct {
	Controller
	statusApplication *application.Status
}

// Create
// @Summary Create new Status
// @Description create new status
// @Tags Statuss
// @Accept  json
// @Produce  json
// @Param status body application.CreateStatusRequest true "Status"
// @Success 201 {object} dto.Status
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /statuss [post]
func (c *Status) Create(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.CreateStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	status, appErr := c.statusApplication.Create(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewStatus(status))
}

// Get
// @Summary Get a status
// @Description get a status by id
// @Tags Statuss
// @Produce  json
// @Success 200 {object} dto.Status
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Status id"
// @Router /statuss/{id} [get]
func (c *Status) Get(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	status, appErr := c.statusApplication.Get(requestContext, ctx.Param("id"))
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewStatus(status))
}

// GetAll
// @Summary Get all status
// @Description get all status by id
// @Tags Statuss
// @Produce  json
// @Success 200 {array} dto.Status
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /statuss [get]
func (c *Status) GetAll(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	var requestParams database.GetAllStatusesRequest
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	records, appErr := c.statusApplication.GetAll(requestContext, requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	Statuss := make([]*dto.Status, 0)
	for _, v := range records {
		Statuss = append(Statuss, dto.NewStatus(v))
	}
	ctx.JSON(http.StatusOK, Statuss)
}

// Delete
// @Summary Delete status
// @Description Delete status by id
// @Tags Statuss
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Status id"
// @Router /statuss/{id} [delete]
func (c *Status) Delete(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	id := ctx.Param("id")
	appErr = c.statusApplication.Delete(requestContext, id)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, helpers.ToJson("deleted"))
}

// Update
// @Summary Update a status
// @Description Update a status by id
// @Tags Statuss
// @Produce  json
// @Param status body application.UpdateStatusRequest true "Status"
// @Success 200 {object} dto.Status
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Status id"
// @Router /statuss/{id} [put]
func (c *Status) Update(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.ID = ctx.Param("id")

	status, appErr := c.statusApplication.Update(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewStatus(status))
}
