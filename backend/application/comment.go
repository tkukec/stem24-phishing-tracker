package application

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/exceptions"
	"github.com/rs/zerolog"
)

type Comment struct {
	commentRepo repositories.CommentRepository
	logger      zerolog.Logger
}

func NewComment(commentRepo repositories.CommentRepository,
	logger zerolog.Logger,
) *Comment {
	return &Comment{commentRepo: commentRepo, logger: logger}
}

type CreateCommentRequest struct {
	Description string // Name of the comment
	EventID     string
}

func (a *Comment) Create(ctx *context.RequestContext, request *CreateCommentRequest) (*models.Comment, exceptions.ApplicationException) {
	comment := &models.Comment{
		Description: request.Description,
		EventID:     request.EventID,
	}

	comment, err := a.commentRepo.Persist(ctx.TenantID(), comment)
	if err != nil {
		return nil, exceptions.FailedPersisting(models.CommentModelName, err)
	}
	return comment, nil
}

type UpdateCommentRequest struct {
	Description string // Updated description of the comment
	ID          string // Updated description of the comment
}

func (request *UpdateCommentRequest) ApplyValues(comment *models.Comment) *models.Comment {
	if request.Description != "" {
		comment.Description = request.Description
	}
	// Add similar checks for other fields that can be updated

	return comment
}
func (a *Comment) Update(ctx *context.RequestContext, request *UpdateCommentRequest) (*models.Comment, exceptions.ApplicationException) {
	comment, err := a.commentRepo.Get(ctx.TenantID(), request.ID)
	if err != nil {
		return nil, exceptions.CommentNotFound(err)
	}
	comment, err = a.commentRepo.Update(ctx.TenantID(), request.ApplyValues(comment))
	if err != nil {
		return nil, exceptions.FailedUpdating(models.CommentModelName, err)
	}
	return comment, nil
}

func (a *Comment) Delete(ctx *context.RequestContext, ID string) exceptions.ApplicationException {
	comment, err := a.commentRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return exceptions.CommentNotFound(err)
	}
	if err = a.commentRepo.Delete(ctx.TenantID(), comment); err != nil {
		return exceptions.FailedDeleting(models.TenantModelName, err)
	}
	return nil
}

func (a *Comment) Get(ctx *context.RequestContext, ID string) (*models.Comment, exceptions.ApplicationException) {
	comment, err := a.commentRepo.Get(ctx.TenantID(), ID)
	if err != nil {
		return nil, exceptions.CommentNotFound(err)
	}
	return comment, nil
}

func (a *Comment) GetAll(ctx *context.RequestContext, request database.GetAllCommentsRequest) ([]*models.Comment, exceptions.ApplicationException) {
	comments, err := a.commentRepo.GetAll(ctx.TenantID(), request)
	if err != nil {
		return nil, exceptions.FailedQuerying(models.CommentModelName, err)
	}
	return comments, nil
}
