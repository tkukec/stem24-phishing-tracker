package http

import (
	"github.com/asseco-voice/agent-management/shared"
	"github.com/asseco-voice/agent-management/shared/database"
	"github.com/asseco-voice/agent-management/shared/exceptions"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/asseco-voice/agent-management/application"
	"github.com/asseco-voice/agent-management/domain/models"
	"github.com/asseco-voice/agent-management/infrastructure/dto"
	"github.com/gin-gonic/gin"
)

// NewAgent constructor for Agent
func NewAgent(agentApp *application.Agent, controller Controller) *Agent {
	return &Agent{
		Controller: controller,
		agentApp:   agentApp,
	}
}

// Agent ....
type Agent struct {
	Controller
	agentApp *application.Agent
}

// CreateAgent
// @Summary Create new Agent
// @Description create new agent
// @Tags Agents
// @Accept  json
// @Produce  json
// @Param agent body application.CreateAgentRequest true "Agent"
// @Success 201 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /agents [post]
func (c *Agent) CreateAgent(ctx *gin.Context) {
	var err error
	var request application.CreateAgentRequest
	if err = ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	agent, appErr := c.agentApp.CreateAgent(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusCreated, dto.NewAgent(agent))
}

// UpdateAgent
// @Summary Update new Agent
// @Description update new agent
// @Tags Agents
// @Accept  json
// @Produce  json
// @Param agent body application.UpdateAgentRequest true "Agent"
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "User id"
// @Router /agents/{id} [put]
func (c *Agent) UpdateAgent(ctx *gin.Context) {
	var err error
	var request application.UpdateAgentRequest
	if err = ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request.UserID = ctx.Param("id")
	agent, appErr := c.agentApp.UpdateAgent(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// LoginAgent
// @Summary Log in | log out agent
// @Description log in | log out agent
// @Tags Agents
// @Param request body application.ChangeLoggedInStatusRequest true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /agents/login [post]
func (c *Agent) LoginAgent(ctx *gin.Context) {
	var request application.ChangeLoggedInStatusRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request.LoggedIn = true
	request.UserID = requestContext.User().Claims().UserID
	agent, appErr := c.agentApp.ChangeLoggedInStatus(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// UpdateConnectionInfo
// @Summary Log in | log out agent
// @Description log in | log out agent
// @Tags Agents
// @Param request body application.UpdateConnectionInfoRequest true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Agent id"
// @Router /agents/{id}/connection-info [post]
func (c *Agent) UpdateConnectionInfo(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request application.UpdateConnectionInfoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.AgentID = requestContext.User().ID()
	agent, appErr := c.agentApp.UpdateConnectionInfo(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// LogoutAgent
// @Summary Log in | log out agent
// @Description log in | log out agent
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /agents/logout [post]
func (c *Agent) LogoutAgent(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request := application.ChangeLoggedInStatusRequest{
		LoggedIn: false,
		UserID:   requestContext.User().ID(),
	}

	agent, appErr := c.agentApp.ChangeLoggedInStatus(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// GetAgent
// @Summary Get an agent
// @Description gets an agent by user_id
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "User id"
// @Router /agents/{id} [get]
func (c *Agent) GetAgent(ctx *gin.Context) {
	var agent *models.Agent
	var err error
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request := &application.GetAgentRequest{
		ID: ctx.Param("id"),
	}
	agent, appErr = c.agentApp.GetAgent(requestContext, request)
	if err != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// GetAgents
// @Summary Get all agents
// @Description get all agents
// @Tags Agents
// @Accept  json
// @Produce  json
// @Param        id   				query     string  false  "search by id"
// @Param        global_status_id   query     string  false  "search by global_status_id"
// @Param        extension   		query     string  false  "search by extension"
// @Param        number   			query     string  false  "search by number"
// @Param        display_name   	query     string  false  "search by display_name"
// @Param        skill_group_id   	query     string  false  "search by skill_group_id"
// @Param        channel_id      	query     string  false  "search by any channel id of a status"
// @Param        is_blocked     	query     bool    false  "search by agents that have any unblocked or blocked status"
// @Param        sort   			query	  string  false  "sort result eg. created_at DESC"
// @Param        page   			query     int     false  "page number for pagination, default to: 1"
// @Param        per_page  	 		query     int     false  "number of items per page, defaults to: 15"
// @Param        include_blocked    query     bool    false  "include bocked agents, defaults to true"
// @Param        channel_status_id	query     string  false  "search by channel status id"
// @Success 200 {array} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /agents [get]
func (c *Agent) GetAgents(ctx *gin.Context) {
	var agents []*models.Agent
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var requestParams application.GetAllAgentsRequest

	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	agents, appErr = c.agentApp.GetAllAgents(requestContext, &requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgents(agents))
}

// DeleteAgent
// @Summary Delete an agent
// @Description Delete an agent by user_id
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "User id"
// @Router /agents/{id} [delete]
func (c *Agent) DeleteAgent(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request := &application.DeleteAgentRequest{
		ID: ctx.Param("id"),
	}
	if appErr = c.agentApp.DeleteAgent(requestContext, request); appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	ctx.JSON(http.StatusOK, NewStringResponse("OK"))
}

// ChangeAgentStatus
// @Summary Change status of user
// @Description change status of user to one of predefined users
// @Tags Agents
// @Accept  json
// @Produce  json
// @Param status body application.ChangeStatusRequest true "GlobalStatus"
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "User id"
// @Router /agents/{id}/status [post]
func (c *Agent) ChangeAgentStatus(ctx *gin.Context) {
	var err error
	var request application.ChangeStatusRequest
	if err = ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.UserID = ctx.Param("id")
	request.ExecutedBy = shared.ExecutedByAgent
	if request.Reason == "" {
		request.Reason = shared.ReasonNoReasonByAgent
	}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	agent, appErr := c.agentApp.ChangeStatus(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// ChangeAgentGlobalStatus
// @Summary Change status of user
// @Description change status of user to one of predefined users
// @Tags Agents
// @Accept  json
// @Produce  json
// @Param status body application.ChangeStatusRequest true "GlobalStatus"
// @Success 200 {object} dto.Agent
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "User id"
// @Router /agents/{id}/status [post]
func (c *Agent) ChangeAgentGlobalStatus(ctx *gin.Context) {
	var err error
	var request application.ChangeStatusRequest
	if err = ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	request.UserID = ctx.Param("id")
	request.ExecutedBy = shared.ExecutedByAgent
	request.Reason = shared.ReasonNoReasonByAgent
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	agent, appErr := c.agentApp.ChangeGlobalStatus(requestContext, &request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewAgent(agent))
}

// AgentQueueStatistics
// @Summary Fetch statistics about all the queues the agent has access to
// @Description Fetch statistics about all the queues the agent has access to
// @Tags Queues
// @Accept  json
// @Produce  json
// @Success 200 {array} application.QueueStatisticsOfUserChannelResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Router /queues/me/statistics [get]
func (c *Agent) AgentQueueStatistics(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	request := &application.QueuesOfUserStatisticsRequest{
		UserId: requestContext.User().ID(),
	}
	queues, appErr := c.agentApp.QueuesOfUserStatistics(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, queues)
}

// ExtendedQueues
// @Summary Fetch list of items in the queue
// @Description Fetch list of items in the queue
// @Tags Queues
// @Accept  json
// @Produce  json
// @Success 200 {array} dto.QueuedItem
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param        id   query     string  false  "search by id"
// @Param        session_uuid   query     string  false  "search by session_uuid"
// @Param        queue_count   query     string  false  "search by queue_count"
// @Param        queued   query     bool  false  "search by queued"
// @Param        answered   query     bool  false  "search by answered"
// @Param        agent_id   query     string  false  "search by agent_id"
// @Param        sort   query     string  false  "sort result eg. created_at DESC"
// @Param        page   query     string  false  "page number for pagination, default to: 1"
// @Param        per_page   query     string  false  "number of items per page, defaults to: 15"
// @Param   id path string true "Queue id"
// @Router /queues/{id} [get]
func (c *Agent) ExtendedQueues(ctx *gin.Context) {
	request := &application.ExtendedQueuesRequest{QueueId: ctx.Param("id")}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}

	var requestParams database.GetAllQueuedItemsRequest

	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	queues, appErr := c.agentApp.ExtendedQueues(requestContext, request, requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, dto.NewQueuedItems(queues))
}

// AgentSkillGroups
// @Summary Get all agent skill-groups
// @Description get all agent skill-groups
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {array} dto.SkillGroup
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Agent id"
// @Router /agents/{id}/skill-groups [get]
func (c *Agent) AgentSkillGroups(ctx *gin.Context) {
	request := &application.SkillGroupsOfUserRequest{AgentId: ctx.Param("id")}
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	skillGroups, appErr := c.agentApp.SkillGroupsOfUser(requestContext, request)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	dtoSkillGroups := make([]*dto.SkillGroup, 0)
	for _, skillGroup := range skillGroups {
		dtoSkillGroups = append(dtoSkillGroups, dto.NewSkillGroup(skillGroup))
	}

	ctx.JSON(http.StatusOK, dtoSkillGroups)
}

// GenerateExtension
// @Summary Get all agent skill-groups
// @Description get all agent skill-groups
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {array} clients.Extension
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Agent id"
// @Router /agents/{id}/generate-extension [get]
func (c *Agent) GenerateExtension(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var request application.GenerateExtensionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}

	ext := request.Extension
	if ext == "" {
		ext = requestContext.User().Claims().PreferredUsername
	}

	extension, appErr := c.agentApp.GenerateExtension(
		requestContext, &application.GenerateExtensionRequest{
			AgentID:   ctx.Param("id"),
			Extension: ext,
			Temporary: true,
		})
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, extension)
}

// GetAgentGlobalStatusHistory
// @Summary Get agent global status history
// @Description Get agent global status history
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {array} database.PaginatedResult
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id query string false "search by items ids"
// @Param   agent_id query string false "search by agents ids"
// @Param   status_id query string false "search by status ids"
// @Param   channel_id query string false "search by channel ids"
// @Param   executed_by query string false "search by executed"
// @Param   reason query string false "search by reason"
// @Param   date_time_from query string false "from time ex. 2023-11-30T23:59:59.999Z"
// @Param   date_time_to query string false "to time ex. 2023-11-30T23:59:59.999Z"
// @Param   sort query string false "sort result ex. created_at DESC"
// @Param   page query int false "page number for pagination, default to: 1"
// @Param   per_page query int false "number of items per page, defaults to: 15"
// @Router /agents/status/history [get]
func (c *Agent) GetAgentGlobalStatusHistory(ctx *gin.Context) {
	var globalHistoryItems []*dto.AgentGlobalStatusHistoryItem
	var totalItems *int64
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var requestParams application.GetAllAgentGlobalHistoryItemsRequest
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	globalHistoryItems, totalItems, appErr = c.agentApp.GetAllAgentGlobalHistoryItems(requestContext, &requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, database.NewPaginatedResult(
		dto.NewAgentGlobalStatusHistoryItemCollection(globalHistoryItems), requestParams.Page, requestParams.PerPage, totalItems))
}

// GetAgentChannelStatusHistory
// @Summary Get agent channel status history
// @Description Get agent channel status history
// @Tags ChannelStatusHistory
// @Accept  json
// @Produce  json
// @Success 200 {array} database.PaginatedResult
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id query string false "search by items ids"
// @Param   agent_id query string false "search by agents ids"
// @Param   status_id query string false "search by status ids"
// @Param   channel_id query string false "search by channel ids"
// @Param   executed_by query string false "search by executed"
// @Param   reason query string false "search by reason"
// @Param   date_time_from query string false "from time ex. 2023-11-30T23:59:59.999Z"
// @Param   date_time_to query string false "to time ex. 2023-11-30T23:59:59.999Z"
// @Param   sort query string false "sort result ex. created_at DESC"
// @Param   page query int false "page number for pagination, default to: 1"
// @Param   per_page query int false "number of items per page, defaults to: 15"
// @Router /channel-status/history [get]
func (c *Agent) GetAgentChannelStatusHistory(ctx *gin.Context) {
	var channelHistoryItems []*dto.AgentChannelStatusHistoryItem
	var totalItems *int64
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var requestParams application.GetAllAgentChannelHistoryItemsRequest
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	channelHistoryItems, totalItems, appErr = c.agentApp.GetAllAgentChannelHistoryItems(requestContext, &requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	collection := dto.NewAgentChannelStatusHistoryItemCollection(channelHistoryItems)
	ctx.JSON(http.StatusOK, database.NewPaginatedResult(collection, requestParams.Page, requestParams.PerPage, totalItems))
}

// BlockAgentChannelStatus
// @Summary Block agent channel status
// @Description Block agent channel status
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Agent id"
// @Router /agents/{id}/block [post]
func (c *Agent) BlockAgentChannelStatus(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var requestParams application.BlockAgentChannelStatusRequest
	requestParams.AgentId = ctx.Param("id")
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	response, appErr := c.agentApp.BlockAgentChannelStatus(requestContext, &requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, NewStringResponse(response))
}

// UnBlockAgentChannelStatus
// @Summary Un-Block agent channel status
// @Description Un-Block agent channel status
// @Tags Agents
// @Accept  json
// @Produce  json
// @Success 200 {object} StringResponse
// @Failure 500 {object} exceptions.ApiError
// @Failure 404 {object} exceptions.ApiError
// @Failure 400 {object} exceptions.ApiError
// @Param   id path string true "Agent id"
// @Router /agents/{id}/un-block [post]
func (c *Agent) UnBlockAgentChannelStatus(ctx *gin.Context) {
	requestContext, appErr := c.BuildRequestContext(ctx)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	var requestParams application.UnBlockAgentChannelStatusRequest
	requestParams.AgentId = ctx.Param("id")
	if err := ctx.Bind(&requestParams); err != nil {
		exception := exceptions.UnprocessableEntity(c.ValidationErrors(err.(validator.ValidationErrors)), "")
		ctx.JSON(exception.Status(), exception.ToDto())
		return
	}
	response, appErr := c.agentApp.UnBlockAgentChannelStatus(requestContext, &requestParams)
	if appErr != nil {
		ctx.JSON(appErr.Status(), appErr.ToDto())
		return
	}
	ctx.JSON(http.StatusOK, NewStringResponse(response))
}
