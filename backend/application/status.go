package application

import (
	"fmt"
	assecoExceptions "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/exceptions"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"github.com/asseco-voice/agent-management/domain/models"
	"github.com/asseco-voice/agent-management/infrastructure/repositories"
	"github.com/asseco-voice/agent-management/shared/exceptions"
)

// NewStatus constructor for GlobalStatus
func NewStatus(statusRepo repositories.StatusRepository, globalStatus repositories.GlobalStatusRepository) *Status {
	return &Status{
		globalStatus: globalStatus,
		statusRepo:   statusRepo,
	}
}

// Status ....
type Status struct {
	globalStatus repositories.GlobalStatusRepository
	statusRepo   repositories.StatusRepository
}

type CreateStatusRequest struct {
	System            bool   `json:"system"`
	Blocked           bool   `json:"blocked"`
	StartingStatus    bool   `json:"starting_status"`
	Reason            string `json:"reason" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Label             string `json:"label" binding:"required"`
	Timer             int64  `json:"timer"`
	TimerTransitionID string `json:"timer_transition_id"`
	ChannelID         string `json:"channel_id" binding:"required"`
	OnReject          bool   `json:"on_reject"`
	OnTimeout         bool   `json:"on_timeout"`
	DefaultBlocked    bool   `json:"default_blocked"`
	DefaultUnblocked  bool   `json:"default_unblocked"`
}

func (c *Status) CreateStatus(ctx *context.RequestContext, request *CreateStatusRequest) (*models.Status, assecoExceptions.ApplicationException) {
	var timerStatus *models.Status
	var err error

	if request.Timer != 0 {
		if timerStatus, err = c.statusRepo.Get(ctx.TenantID(), request.TimerTransitionID); err != nil {
			return nil, exceptions.BadRequest(map[string][]string{
				"timer_transition_id": {
					err.Error(),
				},
			}, "")
		}
	}

	if request.StartingStatus {
		oldStartingStatus, err := c.statusRepo.GetByChannelAndIsStarting(ctx.TenantID(), request.ChannelID, true)
		if err == nil && oldStartingStatus != nil {
			oldStartingStatus.StartingStatus = false
			if _, err = c.statusRepo.Update(ctx.TenantID(), oldStartingStatus); err != nil {
				return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
			}
		}
	}

	channelStatus, err := c.statusRepo.Persist(ctx.TenantID(),
		models.NewStatus(
			request.System,
			request.Blocked,
			request.StartingStatus,
			request.OnReject,
			request.OnTimeout,
			request.Reason,
			request.Name,
			request.Label,
			request.Timer,
			timerStatus,
			request.ChannelID,
			request.DefaultBlocked,
			request.DefaultUnblocked))

	if err != nil {
		return nil, exceptions.FailedPersisting(models.ChannelStatusModelName, err)
	}

	if err = c.syncSiblingStatuses(ctx, channelStatus); err != nil {
		return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
	}

	return channelStatus, nil
}

func (c *Status) syncSiblingStatuses(ctx *context.RequestContext, status *models.Status) error {
	statusesOfChannel, err := c.statusRepo.GetByChannel(ctx.TenantID(), status.ChannelID, "Transitions")
	if err != nil {
		return err
	}

	for _, channelStatus := range statusesOfChannel {
		if status.ID == channelStatus.ID {
			continue
		}

		if status.DefaultBlocked {
			channelStatus.DefaultBlocked = false
		}
		if status.DefaultUnblocked {
			channelStatus.DefaultUnblocked = false
		}
		_, err = c.statusRepo.Update(ctx.TenantID(), channelStatus)
		if err != nil {
			return err
		}
	}

	return nil
}

type UpdateStatusRequest struct {
	ID                string  `json:"-"`
	Name              string  `json:"name"`
	Label             string  `json:"label"`
	Blocked           bool    `json:"blocked"`
	Reason            string  `json:"reason"`
	OnReject          bool    `json:"on_reject"`
	OnTimeout         bool    `json:"on_timeout"`
	TimerTransitionID *string `json:"timer_transition_id"`
	Timer             int64   `json:"timer"`
	StartingStatus    bool    `json:"starting_status"`
	DefaultBlocked    bool    `json:"default_blocked"`
	DefaultUnblocked  bool    `json:"default_unblocked"`
}

func (r *UpdateStatusRequest) UpdateValues(status *models.Status) (*models.Status, error) {
	if status.Name != r.Name && !status.System {
		return nil, fmt.Errorf("can not change name of system status")
	}

	if r.DefaultBlocked && r.Blocked {
		status.DefaultBlocked = true
	}
	if r.DefaultUnblocked && !r.Blocked {
		status.DefaultUnblocked = true
	}

	status.Label = r.Label
	status.Blocked = r.Blocked
	status.Reason = r.Reason
	status.OnReject = r.OnReject
	status.OnTimeout = r.OnTimeout
	status.TimerTransitionID = r.TimerTransitionID
	status.Timer = r.Timer
	status.StartingStatus = r.StartingStatus
	return status, nil
}

func (c *Status) UpdateStatus(ctx *context.RequestContext, request *UpdateStatusRequest) (*models.Status, assecoExceptions.ApplicationException) {
	var status *models.Status
	var err error
	if status, err = c.statusRepo.Get(ctx.TenantID(), request.ID); err != nil {
		return nil, exceptions.ChannelStatusNotFound(err)
	}
	status, err = request.UpdateValues(status)
	if err != nil {
		return nil, exceptions.BadRequest(map[string][]string{
			"status": {
				err.Error(),
			},
		}, "")
	}

	if status.TimerTransitionID != nil {
		if _, err = c.statusRepo.Get(ctx.TenantID(), *status.TimerTransitionID); err != nil {
			return nil, exceptions.BadRequest(map[string][]string{
				"timer_transition_id": {
					err.Error(),
				},
			}, "")
		}
	}

	if status.OnReject {
		existingRejectStatus, err := c.statusRepo.GetOnRejectStatus(ctx.TenantID(), status.ChannelID)
		if err == nil {
			existingRejectStatus.OnReject = false
			if _, err = c.statusRepo.Update(ctx.TenantID(), existingRejectStatus); err != nil {
				return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
			}
		}
	}
	if status.OnTimeout {
		existingTimeoutStatus, err := c.statusRepo.GetOnTimeoutStatus(ctx.TenantID(), status.ChannelID)
		if err == nil {
			existingTimeoutStatus.OnTimeout = false
			if _, err = c.statusRepo.Update(ctx.TenantID(), existingTimeoutStatus); err != nil {
				return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
			}
		}
	}

	if request.StartingStatus {
		oldStartingStatus, err := c.statusRepo.GetByChannelAndIsStarting(ctx.TenantID(), status.ChannelID, true)
		if err == nil && oldStartingStatus != nil {
			oldStartingStatus.StartingStatus = false
			if _, err = c.statusRepo.Update(ctx.TenantID(), oldStartingStatus); err != nil {
				return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
			}
		}
	}

	if status, err = c.statusRepo.Update(ctx.TenantID(), status); err != nil {
		return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
	}

	if err = c.syncSiblingStatuses(ctx, status); err != nil {
		return nil, exceptions.FailedUpdating(models.ChannelStatusModelName, err)
	}

	return status, nil
}

func (c *Status) DeleteStatus(ctx *context.RequestContext, ID string) assecoExceptions.ApplicationException {
	var status *models.Status
	var err error
	if status, err = c.statusRepo.Get(ctx.TenantID(), ID); err != nil {
		return exceptions.ChannelStatusNotFound(err)
	}
	if err = c.statusRepo.Delete(ctx.TenantID(), status); err != nil {
		return exceptions.FailedDeleting(models.ChannelStatusModelName, err)
	}
	return nil
}

func (c *Status) GetStatus(ctx *context.RequestContext, ID string) (*models.Status, assecoExceptions.ApplicationException) {
	channelStatus, err := c.statusRepo.Get(ctx.TenantID(), ID, "Transitions.TimerTransition", "TimerTransition")
	if err != nil {
		return nil, exceptions.ChannelStatusNotFound(err)
	}
	return channelStatus, nil
}

type GetAllStatusesRequest struct {
	Channels []string `json:"channel_id" form:"channel_id"`
}

func (c *Status) GetAllStatuses(ctx *context.RequestContext, request *GetAllStatusesRequest) ([]*models.Status, assecoExceptions.ApplicationException) {
	if len(request.Channels) > 0 {
		statuses, err := c.statusRepo.GetByChannels(ctx.TenantID(), request.Channels, "Transitions.TimerTransition", "TimerTransition")
		if err != nil {
			return nil, exceptions.FailedQuerying(models.ChannelStatusModelName, err)
		}
		return statuses, nil
	}

	channelStatuses, err := c.statusRepo.GetAll(ctx.TenantID(), "Transitions.TimerTransition", "TimerTransition")
	if err != nil {
		return nil, exceptions.FailedQuerying(models.ChannelStatusModelName, err)
	}
	return channelStatuses, nil
}
