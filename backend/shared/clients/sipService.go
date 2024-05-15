package clients

import (
	"encoding/json"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"log"
	"net/http"
	"time"
)

const (
	SipServiceUriEnvKey = "SIP_SERVICE_URI"

	callEndpoint          = "/api/call"
	sessionsEndpoint      = "/api/sessions"
	transfersEndpoint     = "/api/transfers"
	conferencesEndpoint   = "/api/conferences/"
	ivrsEndpoint          = "/api/ivrs/"
	extensionsEndpoint    = "/api/extensions"
	extensionEndpointPart = "/extensions/"
	pauseCall             = "/pause"
	unPauseCall           = "/un-pause"
	joinConference        = "/join"
	transferCall          = "/transfer"
	requeueCall           = "/requeue"
	startRecording        = "/recording/start"
	stopRecording         = "/recording/stop"
	resumeRecording       = "/recording/resume"
	pauseRecording        = "/recording/pause"
	muteCall              = "/mute"
	unMuteCall            = "/un-mute"
	executeActions        = "/actions"
	initiateCall          = "/call"
	messagesEndpoint      = "/api/messages"
	legsEndpoint          = "/api/legs"
	answerQueueEndpoint   = "/api/answer-queue"
	plans                 = "/plans"
	blindTransfer         = "/blind-transfer"
	conferences           = "conferences"
	reject                = "reject"
)

type SipService struct {
	uri string
}

func NewSipService() *SipService {
	return &SipService{
		uri: helpers.UrlParse(runtimebag.GetEnvString(
			SipServiceUriEnvKey,
			"http://sip-service:8080",
		)),
	}
}

func (s *SipService) SetUri(uri string) *SipService {
	s.uri = helpers.UrlParse(uri)
	return s
}

func (s *SipService) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(s.uri, method)
}

type InitiateCallRequest struct {
	Caller      *ConferenceAttendee    `json:"caller"`
	Callee      *ConferenceAttendee    `json:"callee"`
	ContextData map[string]interface{} `json:"context_data"`
	UseAmd      bool                   `json:"use_amd"`
}

type CallResponse struct {
	Message     string `json:"message"`
	SessionUuid string `json:"session_uuid"`
}

type InitiateCallResponse struct {
	Response
	result *CallResponse
}

func (r *InitiateCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) InitiateCall(ctx *assecoContext.RequestContext, request *InitiateCallRequest, headers map[string]string) *InitiateCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(callEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &InitiateCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type EndCallRequest struct {
	SessionUuid string `json:"session_uuid"`
}

type EndCallResponse struct {
	Response
	result *CallResponse
}

func (r *EndCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) EndCall(ctx *assecoContext.RequestContext, request *EndCallRequest, headers map[string]string) *EndCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodDelete).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &EndCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewActiveSessionsRequest(sessionUuids []string) *ActiveSessionsRequest {
	return &ActiveSessionsRequest{SessionUuids: sessionUuids}
}

type ActiveSessionsRequest struct {
	SessionUuids []string `json:"session_uuids,omitempty"`
}

type ActiveSessionsResponse struct {
	Response
	result *[]Session
}

func (r *ActiveSessionsResponse) Result() *[]Session {
	return r.result
}

func (s *SipService) ActiveSessions(ctx *assecoContext.RequestContext, request *ActiveSessionsRequest, headers map[string]string) *ActiveSessionsResponse {
	var result *[]Session
	req := s.NewRequest(http.MethodGet).
		AddPath(sessionsEndpoint).
		AddMultipleQueryValues("id", request.SessionUuids).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &ActiveSessionsResponse{
		Response: req.Response(),
		result:   result,
	}
}

type CloseConferenceRequest struct {
	SessionUuid    string `json:"session_uuid"`
	ConferenceUuid string `json:"conference_uuid"`
}

type CloseConferenceResponse struct {
	Response
	result *CallResponse
}

func (r *CloseConferenceResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) CloseConference(ctx *assecoContext.RequestContext, request *CloseConferenceRequest, headers map[string]string) *CloseConferenceResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodDelete).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath("conferences").
		AddPath(request.ConferenceUuid).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &CloseConferenceResponse{
		Response: req.Response(),
		result:   result,
	}
}

type LeaveCallRequest struct {
	SessionUuid string `json:"session_uuid"`
	Extension   string `json:"extension"`
}

type LeaveCallResponse struct {
	Response
	result *CallResponse
}

func (r *LeaveCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) LeaveCall(ctx *assecoContext.RequestContext, request *LeaveCallRequest, headers map[string]string) *LeaveCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodDelete).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(extensionEndpointPart).
		AddPath(request.Extension).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &LeaveCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type PauseCallRequest struct {
	SessionUuid   string `json:"-"`
	Number        string `json:"number"`
	PlayHoldMusic bool   `json:"play_hold_music"`
}

type PauseCallResponse struct {
	Response
	result *CallResponse
}

func (r *PauseCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) PauseCall(ctx *assecoContext.RequestContext, request *PauseCallRequest, headers map[string]string) *PauseCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(pauseCall).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &PauseCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type UnPauseCallRequest struct {
	SessionUuid string `json:"-"`
	Number      string `json:"number"`
}

type UnPauseCallResponse struct {
	Response
	result *CallResponse
}

func (r *UnPauseCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) UnPauseCall(ctx *assecoContext.RequestContext, request *UnPauseCallRequest, headers map[string]string) *UnPauseCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(unPauseCall).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &UnPauseCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type AnswerQueueRequest struct {
	ContextData map[string]interface{} `json:"context_data"`
	Agent       *ConferenceAttendee    `json:"agent"`
	SessionUuid string                 `json:"session_uuid"`
	AgentInfo   *AgentData             `json:"agent_info"`
}

type AgentData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AgentSipName string `json:"agent_sip_name"`
	AgentId      string `json:"agent_id"`
}

type AnswerQueueResponse struct {
	Response
	result *CallResponse
}

func (r *AnswerQueueResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) AnswerQueue(ctx *assecoContext.RequestContext, request *AnswerQueueRequest, headers map[string]string) *AnswerQueueResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(answerQueueEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &AnswerQueueResponse{
		Response: req.Response(),
		result:   result,
	}
}

type JoinConferenceRequest struct {
	ConferenceUuid   string              `json:"conference_uuid"`
	Agent            *ConferenceAttendee `json:"agent"`
	SkillGroupNumber string              `json:"skill_group_number"`
}

type JoinConferenceResponse struct {
	Response
	result *Conference
}

func (r *JoinConferenceResponse) Result() *Conference {
	return r.result
}

func (s *SipService) JoinConference(ctx *assecoContext.RequestContext, request *JoinConferenceRequest, headers map[string]string) *JoinConferenceResponse {
	var result *Conference
	req := s.NewRequest(http.MethodPost).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(joinConference).
		SetHeaders(headers).
		SetBody(request).
		Bind(&result).
		Print().
		Run(ctx)

	return &JoinConferenceResponse{
		Response: req.Response(),
		result:   result,
	}
}

type GetSessionResponse struct {
	Response
	result *Session
}

func (r *GetSessionResponse) Result() *Session {
	return r.result
}

func (s *SipService) GetSession(ctx *assecoContext.RequestContext, sessionUuid string, headers map[string]string) *GetSessionResponse {
	var result *Session
	req := s.NewRequest(http.MethodGet).
		AddPath(sessionsEndpoint).
		AddPath(sessionUuid).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &GetSessionResponse{
		Response: req.Response(),
		result:   result,
	}
}

type SendMessageRequest struct {
	TargetExtension string `json:"target_extension" binding:"required"`
	From            string `json:"from" binding:"required"`
	Message         string `json:"message" binding:"required"`
}

type SendMessageResponse struct {
	Response
	result *StringResponse
}

func (r *SendMessageResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) SendMessage(ctx *assecoContext.RequestContext, request SendMessageRequest, headers map[string]string) *SendMessageResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(messagesEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &SendMessageResponse{
		Response: req.Response(),
		result:   result,
	}
}

type UpdateLegRequest struct {
	ContextData   map[string]interface{} `json:"context_data"`
	QueuePosition int                    `json:"queue_position"`
	QueueTime     int                    `json:"queue_time"`
}

type UpdateLegResponse struct {
	Response
	result *StringResponse
}

func (r *UpdateLegResponse) Result() *StringResponse {
	return r.result
}

// UpdateLeg identifier can be leg uuid, extension or number
func (s *SipService) UpdateLeg(ctx *assecoContext.RequestContext, identifier string, request UpdateLegRequest, headers map[string]string) *UpdateLegResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(legsEndpoint).
		AddPath(identifier).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &UpdateLegResponse{
		Response: req.Response(),
		result:   result,
	}
}

type UpdateQueuePositionRequest struct {
	TargetExtension string `json:"target_extension" binding:"required"`
	QueuePosition   string `json:"queue_position"`
	QueueTime       string `json:"queue_time"`
}

func (s *SipService) UpdateQueuePosition(ctx *assecoContext.RequestContext, request UpdateQueuePositionRequest, headers map[string]string) *SendMessageResponse {
	msg := map[string]interface{}{
		"info": map[string]interface{}{
			"queueorder": request.QueuePosition,
			"queuetime":  request.QueueTime,
		},
	}

	return s.SendMessage(ctx, SendMessageRequest{
		TargetExtension: request.TargetExtension,
		From:            "Service",
		Message:         helpers.ToString(msg),
	}, headers)
}

type GetExtensionsResponse struct {
	Response
	result []*Extension
}

func (r *GetExtensionsResponse) Result() []*Extension {
	return r.result
}

func (s *SipService) GetExtensions(ctx *assecoContext.RequestContext, extensions []string, headers map[string]string) *GetExtensionsResponse {
	var result []*Extension
	req := s.NewRequest(http.MethodGet).
		AddPath(extensionsEndpoint).
		AddMultipleQueryValues("extension", extensions).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &GetExtensionsResponse{
		Response: req.Response(),
		result:   result,
	}
}

type InviteRequest struct {
	Name             string                `json:"name"`
	Type             string                `json:"type"`
	SessionUuid      string                `json:"session_uuid"`
	ConferenceUuid   string                `json:"conference_uuid"`
	Attendees        []*ConferenceAttendee `json:"attendees"`
	TransferExisting bool                  `json:"transfer_existing"`
}

type InviteCallResponse struct {
	Message        string `json:"message"`
	SessionUuid    string `json:"session_uuid"`
	ConferenceUuid string `json:"conference_uuid"`
}

type InviteToConferenceResponse struct {
	Response
	result *InviteCallResponse
}

func (r *InviteToConferenceResponse) Result() *InviteCallResponse {
	return r.result
}

func (s *SipService) InviteToConference(ctx *assecoContext.RequestContext, request *InviteRequest, headers map[string]string) *InviteToConferenceResponse {
	var result *InviteCallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(conferencesEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &InviteToConferenceResponse{
		Response: req.Response(),
		result:   result,
	}
}

type TransferCallRequest struct {
	ChannelUuid        string `json:"channel_uuid"`
	Number             string `json:"number"`
	ConferenceUuid     string `json:"-"`
	StartRecording     bool   `json:"start_recording"`
	ConferenceType     string `json:"conference_type"`
	ResponsibleAgentId string `json:"responsible_agent_id"`
}

type TransferCallResponse struct {
	Response
	result *CallResponse
}

func (r *TransferCallResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) TransferCall(ctx *assecoContext.RequestContext, request *TransferCallRequest, headers map[string]string) *TransferCallResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(transferCall).
		SetHeaders(headers).
		SetBody(request).
		Bind(&result).
		Print().
		Run(ctx)
	return &TransferCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type TransferToSessionRequest struct {
	ChannelUuid string    `json:"channel_uuid"`
	Number      string    `json:"number"`
	SessionUuid string    `json:"session_uuid"`
	AgentInfo   AgentData `json:"agent_info"`
}

type TransferToSessionResponse struct {
	Response
	result *Session
}

func (r *TransferToSessionResponse) Result() *Session {
	return r.result
}

func (s *SipService) TransferToSession(ctx *assecoContext.RequestContext, request *TransferToSessionRequest, headers map[string]string) *TransferToSessionResponse {
	var result *Session
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(transferCall).
		SetHeaders(headers).
		SetBody(request).
		Bind(&result).
		Print().
		Run(ctx)
	return &TransferToSessionResponse{
		Response: req.Response(),
		result:   result,
	}
}

type RequeueRequest struct {
	SessionUuid   string                 `json:"session_uuid"`
	Number        string                 `json:"number" binding:"required"`
	PlayHoldMusic bool                   `json:"play_hold_music"`
	Data          map[string]interface{} `json:"data"`
}

func (r *RequeueRequest) ToJson() []byte {
	var v []byte
	var err error
	if v, err = json.Marshal(r); err != nil {
		log.Panic(err.Error())
	}
	return v
}

func (r *RequeueRequest) ToJsonString() string {
	return string(r.ToJson())
}

type StringResponse struct {
	Message string `json:"message"`
}

type RequeueCallResponse struct {
	Response
	result *StringResponse
}

func (r *RequeueCallResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) RequeueCall(ctx *assecoContext.RequestContext, request *RequeueRequest, headers map[string]string) *RequeueCallResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(requeueCall).
		SetHeaders(headers).
		SetBody(request).
		Bind(&result).
		Print().
		Run(ctx)
	return &RequeueCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

type CreateOneExtensionRequest struct {
	Extension       string   `json:"extension"`
	Password        string   `json:"password"`
	Gateway         string   `json:"gateway"`
	Hostname        string   `json:"hostname"`
	Temporary       bool     `json:"temporary"`
	VideoCompatible bool     `json:"video_compatible"`
	Directories     []string `json:"directories"`
}

type CreateOneRequestResponse struct {
	Response
	result *Extension
}

func (r *CreateOneRequestResponse) Result() *Extension {
	return r.result
}

func (s *SipService) CreateExtension(ctx *assecoContext.RequestContext, headers map[string]string, request *CreateOneExtensionRequest) *CreateOneRequestResponse {
	var result *Extension
	req := s.NewRequest(http.MethodPost).
		AddPath(extensionsEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &CreateOneRequestResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewStartRecordingRequest(conferenceUuid string) *StartRecordingRequest {
	return &StartRecordingRequest{ConferenceUuid: conferenceUuid}
}

type StartRecordingRequest struct {
	ConferenceUuid string `json:"conference_uuid"`
}

type StartRecordingResponse struct {
	Response
	result *StringResponse
}

func (r *StartRecordingResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) StartRecording(ctx *assecoContext.RequestContext, request *StartRecordingRequest, headers map[string]string) *StartRecordingResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(startRecording).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &StartRecordingResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewStopRecordingRequest(conferenceUuid string) *StopRecordingRequest {
	return &StopRecordingRequest{ConferenceUuid: conferenceUuid}
}

type StopRecordingRequest struct {
	ConferenceUuid string `json:"conference_uuid"`
}

type StopRecordingResponse struct {
	Response
	result *StringResponse
}

func (r *StopRecordingResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) StopRecording(ctx *assecoContext.RequestContext, request *StopRecordingRequest, headers map[string]string) *StopRecordingResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(stopRecording).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &StopRecordingResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewPauseRecordingRequest(conferenceUuid string) *PauseRecordingRequest {
	return &PauseRecordingRequest{ConferenceUuid: conferenceUuid}
}

type PauseRecordingRequest struct {
	ConferenceUuid string `json:"conference_uuid"`
}

type PauseRecordingResponse struct {
	Response
	result *StringResponse
}

func (r *PauseRecordingResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) PauseRecording(ctx *assecoContext.RequestContext, request *PauseRecordingRequest, headers map[string]string) *PauseRecordingResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(pauseRecording).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &PauseRecordingResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewResumeRecordingRequest(conferenceUuid string) *ResumeRecordingRequest {
	return &ResumeRecordingRequest{ConferenceUuid: conferenceUuid}
}

type ResumeRecordingRequest struct {
	ConferenceUuid string `json:"conference_uuid"`
}

type ResumeRecordingResponse struct {
	Response
	result *StringResponse
}

func (r *ResumeRecordingResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) ResumeRecording(ctx *assecoContext.RequestContext, request *ResumeRecordingRequest, headers map[string]string) *ResumeRecordingResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(conferencesEndpoint).
		AddPath(request.ConferenceUuid).
		AddPath(resumeRecording).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &ResumeRecordingResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewMuteCallRequest(sessionUuid string, number string) *MuteCallRequest {
	return &MuteCallRequest{SessionUuid: sessionUuid, Number: number}
}

type MuteCallRequest struct {
	SessionUuid string `json:"session_uuid"`
	Number      string `json:"number"`
}

type MuteCallResponse struct {
	Response
	result *StringResponse
}

func (r *MuteCallResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) MuteCall(ctx *assecoContext.RequestContext, request *MuteCallRequest, headers map[string]string) *MuteCallResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(extensionEndpointPart).
		AddPath(request.Number).
		AddPath(muteCall).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &MuteCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewUnMuteCallRequest(sessionUuid string, number string) *UnMuteCallRequest {
	return &UnMuteCallRequest{SessionUuid: sessionUuid, Number: number}
}

type UnMuteCallRequest struct {
	SessionUuid string `json:"session_uuid"`
	Number      string `json:"number"`
}

type UnMuteCallResponse struct {
	Response
	result *StringResponse
}

func (r *UnMuteCallResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) UnMuteCall(ctx *assecoContext.RequestContext, request *UnMuteCallRequest, headers map[string]string) *UnMuteCallResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(extensionEndpointPart).
		AddPath(request.Number).
		AddPath(unMuteCall).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &UnMuteCallResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewInitiateCallWithIvrRequest(initiateCallRequest InitiateCallRequest, ivrId string) *InitiateCallWithIvrRequest {
	return &InitiateCallWithIvrRequest{InitiateCallRequest: initiateCallRequest, IvrId: ivrId}
}

type InitiateCallWithIvrRequest struct {
	InitiateCallRequest
	IvrId string
}

type InitiateCallWithIvrResponse struct {
	Response
	result *CallResponse
}

func (r *InitiateCallWithIvrResponse) Result() *CallResponse {
	return r.result
}

func (s *SipService) InitiateCallWithIvr(ctx *assecoContext.RequestContext, request *InitiateCallWithIvrRequest, headers map[string]string) *InitiateCallWithIvrResponse {
	var result *CallResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(ivrsEndpoint).
		AddPath(request.IvrId).
		AddPath(initiateCall).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &InitiateCallWithIvrResponse{
		Response: req.Response(),
		result:   result,
	}
}

type ConferenceAttendee struct {
	ViewNumber      string `json:"view_number"`
	Number          string `json:"number" binding:"required"`
	AgentId         string `json:"agent_id"`
	Initiator       bool   `json:"initiator"`
	Silent          bool   `json:"silent"`
	Whisper         bool   `json:"whisper"`
	Supervisor      bool   `json:"supervisor"`
	IsExternalAgent bool   `json:"is_external_agent"`
}

func NewExecuteIvrActionsRequest(sessionUuid string, legNumber string, actions []interface{}) *ExecuteIvrActionsRequest {
	return &ExecuteIvrActionsRequest{SessionUuid: sessionUuid, LegNumber: legNumber, Actions: actions}
}

type ExecuteIvrActionsRequest struct {
	SessionUuid string        `json:"session_uuid"`
	LegNumber   string        `json:"leg_number"`
	Actions     []interface{} `json:"actions"`
}

func (r *ExecuteIvrActionsRequest) ToJson() []byte {
	var v []byte
	var err error
	if v, err = json.Marshal(r); err != nil {
		log.Panic(err.Error())
	}
	return v
}

type ExecuteIvrActionsResponse struct {
	Response
	result *StringResponse
}

func (r *ExecuteIvrActionsResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) ExecuteIvrActions(ctx *assecoContext.RequestContext, request *ExecuteIvrActionsRequest, headers map[string]string) *ExecuteIvrActionsResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(extensionEndpointPart).
		AddPath(request.LegNumber).
		AddPath(executeActions).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &ExecuteIvrActionsResponse{
		Response: req.Response(),
		result:   result,
	}
}

func NewBlindTransferRequest(
	sessionUuid string,
	clientNumber string,
	transferNumber string,
	skillGroupNumber string,
	data map[string]interface{},
) *BlindTransferRequest {
	return &BlindTransferRequest{
		SessionUuid:      sessionUuid,
		ClientNumber:     clientNumber,
		TransferNumber:   transferNumber,
		SkillGroupNumber: skillGroupNumber,
		Data:             data,
	}
}

type BlindTransferRequest struct {
	SessionUuid      string                 `json:"-"`
	ClientNumber     string                 `json:"number" binding:"required"`
	TransferNumber   string                 `json:"blind_number" binding:"required"`
	SkillGroupNumber string                 `json:"skill_group_number" binding:"required"`
	Data             map[string]interface{} `json:"data"`
}

type BlindTransferResponse struct {
	Response
	result *StringResponse
}

func (r *BlindTransferResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) BlindTransfer(ctx *assecoContext.RequestContext, request *BlindTransferRequest, headers map[string]string) *BlindTransferResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPost).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(blindTransfer).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &BlindTransferResponse{
		Response: req.Response(),
		result:   result,
	}
}

type TransferItem struct {
	ProcessId string `json:"process_id"`
	ItemName  string `json:"item_name"`
	Id        string `json:"id"`
	Sent      bool   `json:"sent"`
	Error     string `json:"error"`
	Deleted   bool   `json:"deleted"`
}

type TransferItemResponse struct {
	Response
	result *StringResponse
}

func (r *TransferItemResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) AddTransferredItem(ctx *assecoContext.RequestContext, item TransferItem, headers map[string]string) *TransferItemResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(transfersEndpoint).
		SetBody(item).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &TransferItemResponse{
		Response: req.Response(),
		result:   result,
	}
}

type RejectConferenceRequest struct {
	SessionUuid    string `json:"session_uuid"`
	ConferenceUuid string `json:"conference_uuid"`
	PlayBusyMusic  bool   `json:"play_busy_music"`
}

type RejectConferenceResponse struct {
	Response
	result *StringResponse
}

func (r *RejectConferenceResponse) Result() *StringResponse {
	return r.result
}

func (s *SipService) RejectConference(ctx *assecoContext.RequestContext, request *RejectConferenceRequest, headers map[string]string) *RejectConferenceResponse {
	var result *StringResponse
	req := s.NewRequest(http.MethodPut).
		AddPath(sessionsEndpoint).
		AddPath(request.SessionUuid).
		AddPath(conferences).
		AddPath(request.ConferenceUuid).
		AddPath(reject).
		SetHeaders(headers).
		SetBody(request).
		Bind(&result).
		Print().
		Run(ctx)
	return &RejectConferenceResponse{
		Response: req.Response(),
		result:   result,
	}
}

type Session struct {
	SessionUuid string        `json:"session_uuid"`
	Conferences []*Conference `json:"conferences"`
	FreeLegs    []*Leg        `json:"free_legs"`
	Timestamp   int64         `json:"timestamp"`
	Hostname    string        `json:"hostname"`
}

func (r *Session) FreeLegClientNumber() string {
	if len(r.FreeLegs) == 0 {
		return ""
	}

	l := r.FreeLegs[0]

	if l == nil {
		return ""
	}

	return l.ExtensionData.Number
}

func (r *Session) LegByNumber(number string) *Leg {
	for _, conference := range r.Conferences {
		for _, leg := range conference.Legs {
			if leg.ExtensionData.Number == number {
				return leg
			}
		}
	}
	return nil
}

func (r *Session) Contact() interface{} {
	for _, l := range r.FreeLegs {
		contact := l.Contact()
		if contact != nil {
			return contact
		}
	}
	for _, c := range r.Conferences {
		for _, cl := range c.Legs {
			contact := cl.Contact()
			if contact != nil {
				return contact
			}
		}
	}
	return nil
}

type Conference struct {
	Name              string `json:"name"`
	Type              string `json:"type"`
	SessionUuid       string `json:"session_uuid"`
	ConferenceUuid    string `json:"conference_uuid"`
	Domain            string `json:"domain"`
	Legs              []*Leg `json:"legs"`
	Timestamp         int64  `json:"timestamp"`
	RecordingLocation string `json:"recording_location"`
}

type Leg struct {
	ChannelUuid     string                 `json:"channel_uuid"`
	Initiator       bool                   `json:"is_initiator"`
	State           string                 `json:"state"`
	ExtensionData   *ExtensionData         `json:"extension_data"`
	ContextData     map[string]interface{} `json:"context_data"`
	IvrData         map[string]interface{} `json:"ivr_data"`
	Queued          bool                   `json:"queued"`
	ExecutedActions []*ExecutedAction      `json:"executed_actions"`
}

func (l *Leg) Contact() interface{} {
	contact, ok := l.IvrData["contact"]
	if !ok {
		return nil
	}
	return contact
}

type ExtensionData struct {
	Number     string `json:"number"`
	AgentId    string `json:"agent_id"`
	Extension  string `json:"extension"`
	ViewNumber string `json:"view_number"`
}

type ExecutedAction struct {
	Name      string                 `json:"name"`
	Result    map[string]interface{} `json:"result"`
	Timestamp int64                  `json:"timestamp"`
}

type Extension struct {
	ID          string       `json:"id"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	DeletedAt   string       `json:"deleted_at"`
	Extension   string       `json:"extension"`
	Gateway     string       `json:"gateway"`
	Password    string       `json:"password"`
	Hostname    string       `json:"hostname"`
	Temporary   bool         `json:"temporary"`
	Directories []*Directory `json:"directories"`
}

type Directory struct {
	XML  string `json:"xml"`
	Rule *Rule  `json:"rule"`
}

type Rule struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`

	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Evaluations []*Evaluation `json:"evaluations"`
	Actions     []*Action     `json:"actions"`
	Dialplan    *Dialplan     `json:"dialplan"`
	Directory   *Directory    `json:"directory"`
	InEffect    bool          `json:"in_effect"`
}

type Evaluation struct {
	Field     string `json:"field"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type Action struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

type Dialplan struct {
	XML string `json:"xml"`
}
