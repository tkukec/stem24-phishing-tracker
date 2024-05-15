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

// NewComments constructor for Comments
func NewComments(commentsApplication *application.Comment, controller Controller) *Comments {
	return &Comments{
		Controller:          controller,
		commentsApplication: commentsApplication,
	}
}

// Comments ....
type Comments struct {
	Controller
	commentsApplication *application.Comment
}

// Create
// @Summary Create new Comments
// @Description create new comments
// @Tags Comments
// @Accept  json
// @Produce  json
// @Param comments body application.CreateCommentsRequest true "Comments"
// @Success 201 {object} dto.Comments
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /commentss [post]
func (c *Comments) Create(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	comments, appErr := c.commentsApplication.Create(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewComment(comments))
}

// Get
// @Summary Get a comments
// @Description get a comments by id
// @Tags Comments
// @Produce  json
// @Success 200 {object} dto.Comments
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Comments id"
// @Router /commentss/{id} [get]
func (c *Comments) Get(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	comments, appErr := c.commentsApplication.Get(requestContext, ctx.Param("id"))
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewComment(comments))
}

// GetAll
// @Summary Get all comments
// @Description get all comments by id
// @Tags Comments
// @Produce  json
// @Success 200 {array} dto.Comments
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /commentss [get]
func (c *Comments) GetAll(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	var requestParams database.GetAllCommentsRequest
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	records, appErr := c.commentsApplication.GetAll(requestContext, requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	Commentss := make([]*dto.Comment, 0)
	for _, v := range records {
		Commentss = append(Commentss, dto.NewComment(v))
	}
	ctx.JSON(http.StatusOK, Commentss)
}

// Delete
// @Summary Delete comments
// @Description Delete comments by id
// @Tags Comments
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Comments id"
// @Router /commentss/{id} [delete]
func (c *Comments) Delete(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	id := ctx.Param("id")
	appErr = c.commentsApplication.Delete(requestContext, id)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, helpers.ToJson("deleted"))
}

// Update
// @Summary Update a comments
// @Description Update a comments by id
// @Tags Comments
// @Produce  json
// @Param comments body application.UpdateCommentsRequest true "Comments"
// @Success 200 {object} dto.Comments
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Comments id"
// @Router /commentss/{id} [put]
func (c *Comments) Update(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request *application.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.ID = ctx.Param("id")

	comments, appErr := c.commentsApplication.Update(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewComment(comments))
}
