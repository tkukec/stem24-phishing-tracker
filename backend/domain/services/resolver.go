package services

import (
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/clients"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"time"
)

type Resolver struct {
	sipService        *clients.SipService
	telephoneService  *clients.Telephone
	videoChatService  *clients.VideoChat
	assecoChatService *clients.AssecoChat
	logger            zerolog.Logger
}

func NewResolver(
	sipService *clients.SipService,
	telephoneService *clients.Telephone,
	videoChatService *clients.VideoChat,
	assecoChatService *clients.AssecoChat,
	logger zerolog.Logger,
) Resolver {
	return Resolver{
		sipService:        sipService,
		telephoneService:  telephoneService,
		videoChatService:  videoChatService,
		assecoChatService: assecoChatService,
		logger:            logger,
	}
}

func logFinishTime(log zerolog.Logger, message string, start time.Time) {
	finish := time.Now()
	diff := finish.Sub(start)
	log.Debug().Msgf("finished %s at %s, finished in %d ms", message, finish.String(), diff.Milliseconds())
}

func (r *Resolver) GetSession(ctx *context.RequestContext, headers map[string]string, sessionUuid string) (*clients.Session, error) {
	log := ctx.BuildLog(r.logger, "services.Resolver.GetSession")
	start := time.Now()
	defer logFinishTime(log, fmt.Sprintf("%s %s", "fetching session", sessionUuid), start)
	log.Debug().Msgf("started fetching session %s at %s", sessionUuid, start.String())

	rCtx, cancel := ctx.WithTimeout(4 * time.Second)
	defer cancel()
	result := r.sipService.GetSession(rCtx, sessionUuid, headers)
	if result.Failed() {
		log.Error().Err(result.AsError()).Msg("failed fetching session data from sip service")
		return nil, result.AsError()
	}

	return result.Result(), nil
}

func (r *Resolver) GetSessions(ctx *context.RequestContext, headers map[string]string, sessionUuids []string) (*[]clients.Session, error) {
	log := ctx.BuildLog(r.logger, "services.Resolver.GetSessions")
	start := time.Now()
	defer logFinishTime(log, fmt.Sprintf("%s %s", "fetching session", strings.Join(sessionUuids, " ")), start)
	log.Debug().Msgf("started fetching session %s at %s", strings.Join(sessionUuids, " "), start.String())

	rCtx, cancel := ctx.WithTimeout(4 * time.Second)
	defer cancel()
	result := r.sipService.ActiveSessions(rCtx, clients.NewActiveSessionsRequest(sessionUuids), headers)
	if result.Failed() {
		log.Error().Err(result.AsError()).Msg("failed fetching session data from sip service")
		return nil, result.AsError()
	}

	return result.Result(), nil
}

func (r *Resolver) UpdateQueuePosition(ctx *context.RequestContext, headers map[string]string, extension string, position, time int) error {
	response := r.sipService.UpdateQueuePosition(ctx, clients.UpdateQueuePositionRequest{
		TargetExtension: extension,
		QueuePosition:   strconv.Itoa(position),
		QueueTime:       strconv.Itoa(time),
	}, headers)
	if response.Failed() {
		return response.AsError()
	}
	return nil
}

func (r *Resolver) UpdateLeg(ctx *context.RequestContext, headers map[string]string, extension string, position, time int, contextData map[string]interface{}) error {
	response := r.sipService.UpdateLeg(ctx, extension, clients.UpdateLegRequest{
		ContextData:   contextData,
		QueuePosition: position,
		QueueTime:     time,
	}, headers)
	if response.Failed() {
		return response.AsError()
	}
	return nil
}

func (r *Resolver) GetTelephoneActivityData(ctx *context.RequestContext, headers map[string]string, sessionUuid string) (clients.CallCollection, error) {
	log := ctx.BuildLog(r.logger, "services.Resolver.GetTelephoneActivityData")
	start := time.Now()
	defer logFinishTime(log, fmt.Sprintf("%s %s", "telephone activity by session_uuid", sessionUuid), start)
	log.Debug().Msgf("started fetching telephone activity by session_uuid %s at %s", sessionUuid, start.String())

	result := r.telephoneService.FindCallBySessionUuid(
		ctx,
		clients.FindCallBySessionUuidRequest{SessionUuid: sessionUuid},
		headers)
	if result.Failed() {
		log.Error().Err(result.Error()).Msg("failed fetching telephone activity")
		return nil, result.AsError()
	}
	return result.Result(), nil
}

func (r *Resolver) GetVideoChatActivityData(ctx *context.RequestContext, headers map[string]string, sessionUuid string) (clients.CallCollection, error) {
	log := ctx.BuildLog(r.logger, "services.Resolver.GetVideoChatActivityData")
	start := time.Now()
	defer logFinishTime(log, fmt.Sprintf("%s %s", "video chat activity by session_uuid", sessionUuid), start)
	log.Debug().Msgf("started fetching video chat activity by session_uuid %s at %s", sessionUuid, start.String())

	result := r.videoChatService.FindCallBySessionUuid(
		ctx,
		clients.FindCallBySessionUuidRequest{SessionUuid: sessionUuid},
		headers)
	if result.Failed() {
		log.Error().Err(result.Error()).Msg("failed fetching video chat activity")
		return nil, result.AsError()
	}
	return result.Result(), nil
}

func (r *Resolver) GetAssecoChatActivityData(ctx *context.RequestContext, headers map[string]string, sessionUuid string) (clients.SocialNetworksConversationCollection, error) {
	log := ctx.BuildLog(r.logger, "services.Resolver.GetAssecoChatActivityData")
	start := time.Now()
	defer logFinishTime(log, fmt.Sprintf("%s %s", "asseco chat activity by session_uuid", sessionUuid), start)
	log.Debug().Msgf("started fetching asseco chat activity by session_uuid %s at %s", sessionUuid, start.String())

	result := r.assecoChatService.FindConversationByProviderConversationId(
		ctx,
		sessionUuid,
		headers,
	)
	if result.Failed() {
		log.Error().Err(result.AsError()).Msg("failed fetching asseco chat activity")
		return nil, result.AsError()
	}
	return result.Result(), nil
}
