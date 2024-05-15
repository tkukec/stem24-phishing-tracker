package repositories

import (
	"encoding/json"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"github.com/asseco-voice/agent-management/domain/models"
	"github.com/asseco-voice/agent-management/shared"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"os"
)

type TenantSeeder struct {
	tenantRepo         TenantRepository
	channelRepo        ChannelRepository
	globalStatusRepo   GlobalStatusRepository
	channelStatusRepo  StatusRepository
	signalRepo         SignalRepository
	activityStatusRepo ActivityStatusRepository
	ponderRepository   PonderRepository
	logger             zerolog.Logger
}

func NewTenantSeeder(
	tenantRepo TenantRepository,
	channelRepo ChannelRepository,
	globalStatusRepo GlobalStatusRepository,
	channelStatusRepo StatusRepository,
	signalRepo SignalRepository,
	activityStatusRepo ActivityStatusRepository,
	ponderRepository PonderRepository,
	logger zerolog.Logger) *TenantSeeder {
	return &TenantSeeder{
		tenantRepo:         tenantRepo,
		channelRepo:        channelRepo,
		globalStatusRepo:   globalStatusRepo,
		channelStatusRepo:  channelStatusRepo,
		signalRepo:         signalRepo,
		activityStatusRepo: activityStatusRepo,
		ponderRepository:   ponderRepository,
		logger:             logger}
}

func (t *TenantSeeder) Description() string {
	return fmt.Sprintf("seeds data from file, if file does not exists seed will be for %s", shared.DefaultTenant)
}

func (t *TenantSeeder) Execute(location string) error {
	if location == "" {
		_, err := t.Run(context.Background(), &NewTenantRequest{
			Name:               shared.DefaultTenant,
			SeedChannels:       models.ChannelSeed(),
			SeedSystemStatuses: models.SystemStatusSeed(),
			SeedBasicStatuses:  models.BasicStatusSeed(),
			SeedGlobalStatuses: models.GlobalSeed(),
			SeedSignals:        models.SignalSeed(),
			ActivityStatus:     models.ActivityStatusSeed(),
			Ponders:            append(models.AgentPonders(), models.SkillGroupPonders()...),
		})
		if err != nil {
			return err
		}
		return nil
	}

	jsonFile, err := os.Open(location)
	if err != nil {
		log.Panicln(err.Error())
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var tenants []*models.Tenant
	err = json.Unmarshal(byteValue, &tenants)
	if err != nil {
		log.Panicln(err.Error())
	}

	for _, tenant := range tenants {
		_, err := t.Run(context.Background(), &NewTenantRequest{
			Name:               tenant.Name,
			SeedChannels:       models.ChannelSeed(),
			SeedSystemStatuses: models.SystemStatusSeed(),
			SeedBasicStatuses:  models.BasicStatusSeed(),
			SeedGlobalStatuses: models.GlobalSeed(),
			SeedSignals:        models.SignalSeed(),
			ActivityStatus:     models.ActivityStatusSeed(),
			Ponders:            append(models.AgentPonders(), models.SkillGroupPonders()...),
		})
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

type NewTenantRequest struct {
	ID                 string `json:"-"`
	Name               string `json:"name"`
	SeedChannels       []*models.Channel
	SeedSystemStatuses []*models.Status
	SeedBasicStatuses  []*models.Status
	SeedGlobalStatuses []*models.GlobalStatus
	SeedSignals        []*models.Signal
	ActivityStatus     []*models.ActivityStatus
	Ponders            []*models.Ponder
}

func (t *TenantSeeder) Run(ctx *context.RequestContext, request *NewTenantRequest) (*models.Tenant, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.ChangeGlobalStatus")

	log.Debug().Msgf("Creating new tenet %s", request.Name)
	if tenant, err := t.tenantRepo.GetByName(request.Name); err == nil && tenant != nil {
		log.Debug().Msgf("tenant %s (%s) exists, skipping creation....", request.Name, tenant.ID)
		return tenant, nil
	}

	tenant, err := t.tenantRepo.Persist(&models.Tenant{
		ID:   request.ID,
		Name: request.Name,
	})
	if err != nil {
		log.Debug().Msgf("failed creating new tenant %s with error %s", request.Name, err.Error())
		return nil, fmt.Errorf("failed creating tenant %s with error %s", request.Name, err.Error())
	}
	log.Debug().Msgf("tenet %s created. Starting seed procedure for new tenant....", tenant.ID)

	globalStatuses, err := t.seedGlobalStatuses(ctx, request.SeedGlobalStatuses, tenant, ctx.XCorrelationID())
	if err != nil {
		log.Error().Err(err).Msg("failed seeding global statuses")
		return nil, fmt.Errorf("failed seeding global statuses with error %s", err.Error())
	}
	log.Debug().Msgf("created global statuses %s", shared.ToJsonString(globalStatuses))

	activityStatuses, err := t.seedActivityStatuses(ctx, &seedActivityStatusesRequest{
		tenant:           tenant,
		statusCandidates: request.ActivityStatus,
		XCorrelationID:   ctx.XCorrelationID(),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed seeding activity statuses")
		return nil, fmt.Errorf("failed seeding activity statuses with error %s", err.Error())
	}
	log.Debug().Msgf("created activity statuses %s", shared.ToJsonString(activityStatuses))

	channels, err := t.seedChannels(ctx, &seedChannelsRequest{
		channelCandidates:         request.SeedChannels,
		tenant:                    tenant,
		seedChannelSystemStatuses: request.SeedSystemStatuses,
		seedChannelBasicStatuses:  request.SeedBasicStatuses,
		seedChannelSignals:        request.SeedSignals,
		XCorrelationID:            ctx.XCorrelationID(),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed seeding channels")
		return nil, fmt.Errorf("failed seeding channels with errot %s", err.Error())
	}
	log.Debug().Msgf("created channels %s", shared.ToJsonString(channels))

	t.seedPonders(ctx, &PonderSeedRequest{
		tenant:   tenant,
		ponders:  request.Ponders,
		channels: channels,
	})

	return tenant, nil
}

func (t *TenantSeeder) seedGlobalStatuses(
	ctx *context.RequestContext,
	globalStatusCandidates []*models.GlobalStatus,
	tenant *models.Tenant,
	XCorrelationID string,
) ([]*models.GlobalStatus, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.seedGlobalStatuses")

	globalStatuses := make([]*models.GlobalStatus, 0)
	for _, globalStatusCandidate := range globalStatusCandidates {
		var globalStatus *models.GlobalStatus
		var err error
		globalStatus, err = t.globalStatusRepo.GetByName(tenant.ID, globalStatusCandidate.Name)
		if err != nil {
			log.Debug().Msgf("Global status %s does not exist. Creating it ....", globalStatusCandidate.Name)
			if len(globalStatusCandidate.Transitions) > 0 {
				transitionsStatuses, err := t.seedGlobalStatuses(ctx, globalStatusCandidate.Transitions, tenant, XCorrelationID)
				if err != nil {
					log.Error().Err(err).Msgf("failed creating transitions statuses for status %s", globalStatusCandidate.Name)
					return nil, fmt.Errorf("failed creating transitions statuses for status %s with error %s", globalStatusCandidate.Name, err.Error())
				}
				globalStatusCandidate.Transitions = transitionsStatuses
			}

			if globalStatusCandidate.TimerTransition != nil {
				timerStatus, err := t.seedGlobalStatuses(ctx, []*models.GlobalStatus{globalStatusCandidate.TimerTransition}, tenant, XCorrelationID)
				if err != nil {
					log.Error().Err(err).Msgf("failed creating timer status %s", globalStatusCandidate.TimerTransition.Name)
					return nil, fmt.Errorf("failed creating timer status %s with error %s", globalStatusCandidate.TimerTransition.Name, err.Error())
				}
				globalStatusCandidate.TimerTransition = timerStatus[0]
				globalStatusCandidate.TimerTransitionID = &timerStatus[0].ID
			}

			globalStatusCandidate.Tenant = tenant
			globalStatusCandidate.TenantID = tenant.ID

			globalStatus, err = t.globalStatusRepo.GetByName(tenant.ID, globalStatusCandidate.Name)
			if err != nil {
				if globalStatus, err = t.globalStatusRepo.Persist(tenant.ID, globalStatusCandidate); err != nil {
					log.Error().Err(err).Msgf("failed persisting global status %s", globalStatusCandidate.Name)
					return nil, fmt.Errorf("failed persisting global status %s with error %s", globalStatusCandidate.Name, err.Error())
				}
			} else {
				if globalStatus, err = t.globalStatusRepo.Update(tenant.ID, globalStatusCandidate); err != nil {
					log.Error().Err(err).Msgf("failed persisting global status %s", globalStatusCandidate.Name)
					return nil, fmt.Errorf("failed persisting global status %s with error %s", globalStatusCandidate.Name, err.Error())
				}
			}
		} else {
			log.Debug().Msgf("Global status %s exist. Updating it ....", globalStatusCandidate.Name)
			if len(globalStatusCandidate.Transitions) > 0 {
				transitionsStatuses, err := t.seedGlobalStatuses(ctx, globalStatusCandidate.Transitions, tenant, XCorrelationID)
				if err != nil {
					log.Error().Err(err).Msgf("failed creating transitions statuses for status %s", globalStatusCandidate.Name)
					return nil, fmt.Errorf("failed creating transitions statuses for status %s with error %s", globalStatusCandidate.Name, err.Error())
				}
				globalStatus.Transitions = transitionsStatuses
			}
			if globalStatusCandidate.TimerTransition != nil {
				timerStatus, err := t.seedGlobalStatuses(ctx, []*models.GlobalStatus{globalStatusCandidate.TimerTransition}, tenant, XCorrelationID)
				if err != nil {
					log.Error().Err(err).Msgf("failed creating timer status %s", globalStatusCandidate.TimerTransition.Name)
					return nil, fmt.Errorf("failed creating timer status %s with error %s", globalStatusCandidate.TimerTransition.Name, err.Error())
				}
				globalStatus.TimerTransition = timerStatus[0]
				globalStatus.TimerTransitionID = &timerStatus[0].ID
			}
			if globalStatus, err = t.globalStatusRepo.Update(tenant.ID, globalStatus); err != nil {
				log.Error().Err(err).Msgf("failed persisting global status %s", globalStatus.Name)
				return nil, fmt.Errorf("failed persisting global status %s with error %s", globalStatus.Name, err.Error())
			}
		}
		globalStatuses = append(globalStatuses, globalStatus)
	}
	return globalStatuses, nil
}

type seedChannelsRequest struct {
	channelCandidates         []*models.Channel
	tenant                    *models.Tenant
	seedChannelSystemStatuses []*models.Status
	seedChannelBasicStatuses  []*models.Status
	seedChannelSignals        []*models.Signal
	XCorrelationID            string
}

func (t *TenantSeeder) seedChannels(ctx *context.RequestContext, request *seedChannelsRequest) ([]*models.Channel, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.seedChannels")
	log.Debug().Msgf("seeding channels %s", shared.ToJsonString(request.channelCandidates))
	channels := make([]*models.Channel, 0)
	for _, channelCandidate := range request.channelCandidates {
		var channel *models.Channel
		var err error
		if channel, err = t.channelRepo.GetByName(request.tenant.ID, channelCandidate.Name); err != nil {
			channelCandidate.TenantID = request.tenant.ID
			channelCandidate.Tenant = request.tenant
			if channel, err = t.channelRepo.Persist(request.tenant.ID, channelCandidate); err != nil {
				log.Error().Err(err).Msgf("failed persisting channel %s", channelCandidate.Name)
				return nil, fmt.Errorf("failed persisting channel %s with error %s", channelCandidate.Name, err.Error())
			}
			log.Debug().Msgf("channel %s created on tenant %s", channel.ID, channel.TenantID)
			var statusCandidates []*models.Status
			if shared.IsRtcChannel(channelCandidate.Name) {
				if err = copier.Copy(&statusCandidates, request.seedChannelSystemStatuses); err != nil {
					log.Error().Err(err).Msg("failed copying struct")
				}
			} else {
				if err = copier.Copy(&statusCandidates, request.seedChannelBasicStatuses); err != nil {
					log.Error().Err(err).Msg("failed copying struct")
				}
			}
			log.Printf(shared.ToJsonString(statusCandidates))
			channelStatuses, err := t.seedChannelStatuses(ctx, &seedChannelStatusesRequest{
				channelStatusCandidates: statusCandidates,
				channel:                 channel,
				tenant:                  request.tenant,
				XCorrelationID:          request.XCorrelationID,
			})
			if err != nil {
				log.Error().Err(err).Msgf("failed creating channel %s statuses", channel.ID)
				return nil, fmt.Errorf("failed creating channel %s statuses with error  %s", channel.ID, err.Error())
			}
			log.Debug().Msgf("created channel %s statuses %s", channel.ID, shared.ToJsonString(channelStatuses))

			if shared.IsRtcChannel(channelCandidate.Name) {
				signals, err := t.seedChannelSignals(ctx, &seedChannelSignalsRequest{
					signalCandidates: request.seedChannelSignals,
					channel:          channel,
					tenant:           request.tenant,
					XCorrelationID:   request.XCorrelationID,
				})
				if err != nil {
					log.Error().Err(err).Msgf("failed creating channel %s signals", channel.ID)
					return nil, fmt.Errorf("failed creating channel %s signals with error  %s", channel.ID, err.Error())
				}
				log.Debug().Msgf("created channel %s signals %s", channel.ID, shared.ToJsonString(signals))
			}
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

type seedChannelStatusesRequest struct {
	channelStatusCandidates []*models.Status
	channel                 *models.Channel
	tenant                  *models.Tenant
	XCorrelationID          string
}

func (t *TenantSeeder) seedChannelStatuses(ctx *context.RequestContext, request *seedChannelStatusesRequest) ([]*models.Status, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.seedChannelStatuses")
	log.Debug().Msgf("seeding channel %s statuses %s", request.channel.ID, shared.ToJsonString(request.channelStatusCandidates))

	statuses := make([]*models.Status, 0)
	for _, channelStatusCandidate := range request.channelStatusCandidates {
		var status *models.Status
		var err error
		if status, err = t.channelStatusRepo.GetByNameAndChannel(request.tenant.ID, channelStatusCandidate.Name, request.channel.ID); err != nil {
			log.Debug().Msgf("status %s does not exits. creating it .....", channelStatusCandidate.Name)
			channelStatusCandidate.ID = uuid.NewString()
			channelStatusCandidate.TenantID = request.tenant.ID
			channelStatusCandidate.Tenant = request.tenant
			channelStatusCandidate.ChannelID = request.channel.ID
			channelStatusCandidate.Channel = request.channel

			transitionCandidateStatuses := channelStatusCandidate.Transitions
			channelStatusCandidate.Transitions = nil
			transitionTimerStatus := channelStatusCandidate.TimerTransition
			channelStatusCandidate.TimerTransition = nil

			status, err = t.channelStatusRepo.Persist(request.tenant.ID, channelStatusCandidate)
			if err != nil {
				log.Error().Err(err).Msgf("failed persisting status %s", channelStatusCandidate.Name)
				return nil, fmt.Errorf("failed persisting status %s with error %s", channelStatusCandidate.Name, err.Error())
			}
			var transitionStatuses []*models.Status
			if transitionCandidateStatuses != nil && len(transitionCandidateStatuses) > 0 {
				transitionStatuses, err = t.seedChannelStatuses(ctx, &seedChannelStatusesRequest{
					channelStatusCandidates: transitionCandidateStatuses,
					channel:                 request.channel,
					tenant:                  request.tenant,
					XCorrelationID:          request.XCorrelationID,
				})
				if err != nil {
					log.Error().Err(err).Msgf("failed creating transition statuses for status %s", channelStatusCandidate.Name)
					return nil, fmt.Errorf("failed creating transition statuses for status %s with error %s", channelStatusCandidate.Name, err.Error())
				}
			}

			if transitionTimerStatus != nil {
				timerStatus, err := t.seedChannelStatuses(ctx, &seedChannelStatusesRequest{
					channelStatusCandidates: []*models.Status{transitionTimerStatus},
					channel:                 request.channel,
					tenant:                  request.tenant,
					XCorrelationID:          request.XCorrelationID,
				})
				if err != nil {
					log.Error().Err(err).Msgf("failed creating timer transition %s", channelStatusCandidate.TimerTransition.Name)
					return nil, fmt.Errorf("failed creating timer transition %s with error %s", channelStatusCandidate.TimerTransition.Name, err.Error())
				}
				status.TimerTransition = timerStatus[0]
			}

			status.Transitions = transitionStatuses
			status, err = t.channelStatusRepo.Update(request.tenant.ID, status)
			if err != nil {
				log.Error().Err(err).Msgf("failed persisting channel status %s", channelStatusCandidate.Name)
				return nil, fmt.Errorf("failed persisting channel status %s with error %s", channelStatusCandidate.Name, err.Error())
			}
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

type seedChannelSignalsRequest struct {
	signalCandidates []*models.Signal
	channel          *models.Channel
	tenant           *models.Tenant
	XCorrelationID   string
}

func (t *TenantSeeder) seedChannelSignals(ctx *context.RequestContext, request *seedChannelSignalsRequest) ([]*models.Signal, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.seedChannelSignals")
	log.Debug().Msgf("seeding channel %s signals %s", request.channel.ID, shared.ToJsonString(request.signalCandidates))
	signals := make([]*models.Signal, 0)
	for _, signalCandidate := range request.signalCandidates {
		var signal *models.Signal
		var err error
		if signal, err = t.signalRepo.GetByServiceModelActionAndChannel(
			request.tenant.ID,
			signalCandidate.Service,
			signalCandidate.ModelName,
			signalCandidate.Action,
			request.channel.ID); err != nil {

			if signalCandidate.Status != nil {
				status, err := t.channelStatusRepo.GetByNameAndChannel(request.tenant.ID, signalCandidate.Status.Name, request.channel.ID)
				if err != nil {
					log.Error().Err(err).Msgf("failed to fetch signal status %s", signalCandidate.Status.Name)
					return nil, fmt.Errorf("failed to fetch signal status %s with error %s", signalCandidate.Status.Name, err.Error())
				}
				signalCandidate.StatusID = &status.ID
				signalCandidate.Status = status
			}
			signalCandidate.ID = uuid.NewString()
			signalCandidate.Channel = request.channel
			signalCandidate.ChannelID = request.channel.ID
			signalCandidate.Tenant = request.tenant
			signalCandidate.TenantID = request.tenant.ID

			signal, err = t.signalRepo.Persist(request.tenant.ID, signalCandidate)
			if err != nil {
				log.Error().Err(err).Msgf("failed persisting signal %s for channel %s", signalCandidate.SignalName, request.channel.ID)
				return nil, fmt.Errorf("failed persisting signal %s for channel %s with error %s", signalCandidate.SignalName, request.channel.ID, err.Error())
			}
		}
		signals = append(signals, signal)
	}
	return signals, nil
}

type seedActivityStatusesRequest struct {
	tenant           *models.Tenant
	statusCandidates []*models.ActivityStatus
	XCorrelationID   string
}

func (t *TenantSeeder) seedActivityStatuses(ctx *context.RequestContext, request *seedActivityStatusesRequest) ([]*models.ActivityStatus, error) {
	log := ctx.BuildLog(t.logger, "services.Tenant.seedActivityStatuses")
	activityStatuses := make([]*models.ActivityStatus, 0)
	for _, statusCandidate := range request.statusCandidates {
		var err error
		var activityStatus *models.ActivityStatus
		if activityStatus, err = t.activityStatusRepo.GetByName(request.tenant.ID, statusCandidate.Name); err != nil {
			activityStatus, err = t.activityStatusRepo.Persist(request.tenant.ID, statusCandidate)
			if err != nil {
				log.Error().Err(err).Msgf("failed persisting activity status %s", statusCandidate.Name)
				return nil, fmt.Errorf("failed persisting activity status %s with error %s", statusCandidate.Name, err.Error())
			}
		}
		activityStatuses = append(activityStatuses, activityStatus)
	}
	return activityStatuses, nil
}

type PonderSeedRequest struct {
	tenant   *models.Tenant
	ponders  []*models.Ponder
	channels []*models.Channel
}

func (t *TenantSeeder) seedPonders(ctx *context.RequestContext, request *PonderSeedRequest) {
	log := ctx.BuildLog(t.logger, "service.Tenant.seedPonders")
	for _, ponder := range request.ponders {
		for _, channel := range request.channels {
			persist, err := t.ponderRepository.Persist(request.tenant.ID, models.NewPonder(channel, nil, ponder.Object, ponder.Label, ponder.Name, ponder.Value, ponder.Enabled))
			if err != nil {
				log.Error().Err(err).Msg("failed persisting ponder with error")
				continue
			}
			log.Debug().Msgf("persisted ponder %s", shared.ToJsonString(persist))
			for _, skillGroup := range channel.SkillGroups {
				persist, err = t.ponderRepository.Persist(request.tenant.ID, models.NewPonder(channel, skillGroup, ponder.Object, ponder.Label, ponder.Name, ponder.Value, false))
				if err != nil {
					log.Error().Err(err).Msg("failed persisting ponder with error")
					continue
				}
				log.Debug().Msgf("persisted ponder %s", shared.ToJsonString(persist))
			}
		}
	}
}
