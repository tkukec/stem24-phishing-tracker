package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TelephoneRecordingsEndpoint = "/api/recordings"

	TelephoneUriEnvKey = "TELEPHONE_URI"
)

type Telephone struct {
	uri         string
	searchModel string
}

func NewTelephone() *Telephone {
	return &Telephone{
		uri: helpers.UrlParse(runtimebag.GetEnvString(
			TelephoneUriEnvKey,
			"http://telephone:8080",
		)),
		searchModel: "call",
	}
}

func (t *Telephone) SetUri(uri string) *Telephone {
	t.uri = helpers.UrlParse(uri)
	return t
}

func (t *Telephone) NewRequest(method string) *RequestBuilder {
	return NewRequestBuilder(t.uri, method)
}

type UploadFileRequest struct {
	File          *os.File
	RecordingType string
	SessionUuid   string
}

func (r *UploadFileRequest) ToJson() []byte {
	var v []byte
	var err error
	if v, err = json.Marshal(r); err != nil {
		log.Panic(err.Error())
	}
	return v
}

func (r *UploadFileRequest) ToJsonString() string {
	return string(r.ToJson())
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

type UploadFileResponse struct {
	Response
	result *TelephoneRecording
}

func (r *UploadFileResponse) Result() *TelephoneRecording {
	return r.result
}

func (t *Telephone) UploadFile(ctx *assecoContext.RequestContext, request *UploadFileRequest, headers map[string]string, timeout time.Duration) *UploadFileResponse {
	url := fmt.Sprintf("%s%s", TrimSuffix(t.uri, "/"), TelephoneRecordingsEndpoint)
	recordingType := request.RecordingType
	if recordingType == "" {
		recordingType = "f41f65a8-dd34-41ac-8f37-a3b74422713b"
	}
	name := strings.TrimSuffix(request.File.Name(), filepath.Ext(request.File.Name()))
	postBody := &bytes.Buffer{}
	writer := multipart.NewWriter(postBody)
	err := writer.WriteField("name", name)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding name field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("call_id", request.SessionUuid)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding call_id field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("recording_type_id", recordingType)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding recording_type_id field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	part, err := writer.CreateFormFile("file", filepath.Base(request.File.Name()))
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("failed creating form data for file field, error %s", err.Error()),
			},
			result: nil,
		}
	}
	_, err = io.Copy(part, request.File)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error copying src to dst file, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.Close()
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error closing writer, error %s", err.Error()),
			},
			result: nil,
		}
	}

	client := &http.Client{}
	clientRequest, err := http.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	clientRequest.Header.Add("Content-Type", writer.FormDataContentType())
	clientRequest.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		clientRequest.Header.Add(i, v)
	}
	rCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	defer cancel()
	resp, err := client.Do(clientRequest.WithContext(rCtx))

	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error reading body, error %s", err.Error()),
			},
			result: nil,
		}
	}

	if resp.StatusCode >= 300 {
		return &UploadFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	var result *TelephoneRecording
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &UploadFileResponse{
			Response: Response{
				error: fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
			result: nil,
		}
	}
	_ = resp.Body.Close()

	return &UploadFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type FindCallBySessionUuidRequest struct {
	SessionUuid string `json:"session_uuid"`
}

type FindCallBySessionUuidResponse struct {
	Response
	result CallCollection
}

func (r *FindCallBySessionUuidResponse) Result() CallCollection {
	return r.result
}

func (t *Telephone) FindCallBySessionUuid(ctx *assecoContext.RequestContext, request FindCallBySessionUuidRequest, headers map[string]string) *FindCallBySessionUuidResponse {
	var result CallCollection
	searchPayload := map[string]interface{}{
		"search": map[string]string{
			"session_uuid": fmt.Sprintf("=%s", request.SessionUuid),
		},
		"relations": []string{
			"skill_group",
			"remote_relations",
			"call_direction",
			"process_status",
			"connections.actions",
			"connections.end_reason",
		},
	}
	req := t.NewRequest(http.MethodPost).
		AddPath("api/search").
		AddPath(t.searchModel).
		SetBody(searchPayload).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)
	return &FindCallBySessionUuidResponse{
		Response: req.Response(),
		result:   result,
	}
}

type CallCollection []*Call

func (c CallCollection) First() *Call {
	if c != nil && len(c) > 0 {
		return c[0]
	}
	return nil
}

func (c CallCollection) IsEmpty() bool {
	if c == nil || len(c) == 0 {
		return true
	}
	return false
}

type Call struct {
	Id                 string         `json:"id"`
	IvrVersion         string         `json:"ivr_version"`
	IvrId              string         `json:"ivr_id"`
	PlanId             string         `json:"plan_id"`
	ServiceNumber      string         `json:"service_number"`
	ShortId            string         `json:"short_id"`
	ProcessStatusId    string         `json:"process_status_id"`
	CallDirectionId    string         `json:"call_direction_id"`
	CallbackId         interface{}    `json:"callback_id"`
	AbandonedCallId    interface{}    `json:"abandoned_call_id"`
	ResponsibleAgentId string         `json:"responsible_agent_id"`
	SessionUuid        string         `json:"session_uuid"`
	Successful         bool           `json:"successful"`
	StartedAt          string         `json:"started_at"`
	EndedAt            string         `json:"ended_at"`
	CreatedAt          string         `json:"created_at"`
	CreatedBy          interface{}    `json:"created_by"`
	CreatorType        string         `json:"creator_type"`
	UpdatedAt          string         `json:"updated_at"`
	UpdatedBy          interface{}    `json:"updated_by"`
	UpdaterType        string         `json:"updater_type"`
	DeletedAt          interface{}    `json:"deleted_at"`
	DeletedBy          interface{}    `json:"deleted_by"`
	DeleterType        interface{}    `json:"deleter_type"`
	RemoteRelations    []interface{}  `json:"remote_relations"`
	CallDirection      CallDirection  `json:"call_direction"`
	ProcessStatus      ProcessStatus  `json:"process_status"`
	Connections        []Connection   `json:"connections"`
	SkillGroup         CallSkillGroup `json:"skill_group"`
}

func (c *Call) HasConnections() bool {
	if c.Connections == nil || len(c.Connections) == 0 {
		return false
	}
	return true
}

type CallDirection struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	CreatedAt   string      `json:"created_at"`
	CreatedBy   interface{} `json:"created_by"`
	CreatorType interface{} `json:"creator_type"`
	UpdatedAt   string      `json:"updated_at"`
	UpdatedBy   interface{} `json:"updated_by"`
	UpdaterType interface{} `json:"updater_type"`
	DeletedAt   interface{} `json:"deleted_at"`
	DeletedBy   interface{} `json:"deleted_by"`
	DeleterType interface{} `json:"deleter_type"`
}

type ProcessStatus struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	CreatedAt   string      `json:"created_at"`
	CreatedBy   interface{} `json:"created_by"`
	CreatorType interface{} `json:"creator_type"`
	UpdatedAt   string      `json:"updated_at"`
	UpdatedBy   interface{} `json:"updated_by"`
	UpdaterType interface{} `json:"updater_type"`
	DeletedAt   interface{} `json:"deleted_at"`
	DeletedBy   interface{} `json:"deleted_by"`
	DeleterType interface{} `json:"deleter_type"`
}

type Connection struct {
	Id                    string              `json:"id"`
	CallId                string              `json:"call_id"`
	ChannelUuid           string              `json:"channel_uuid"`
	Number                string              `json:"number"`
	Extension             string              `json:"extension"`
	AgentId               *string             `json:"agent_id"`
	Initiator             bool                `json:"initiator"`
	ConnectionEndReasonId string              `json:"connection_end_reason_id"`
	HappenedAt            string              `json:"happened_at"`
	CreatedAt             string              `json:"created_at"`
	CreatedBy             interface{}         `json:"created_by"`
	CreatorType           string              `json:"creator_type"`
	UpdatedAt             string              `json:"updated_at"`
	UpdatedBy             interface{}         `json:"updated_by"`
	UpdaterType           string              `json:"updater_type"`
	DeletedAt             interface{}         `json:"deleted_at"`
	DeletedBy             interface{}         `json:"deleted_by"`
	DeleterType           interface{}         `json:"deleter_type"`
	Actions               []ConnectionAction  `json:"actions"`
	EndReason             ConnectionEndReason `json:"end_reason"`
}

type ConnectionAction struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	Type           string      `json:"type"`
	ConferenceUuid *string     `json:"conference_uuid"`
	ConferenceType *string     `json:"conference_type"`
	ConferenceName *string     `json:"conference_name"`
	HappenedAt     string      `json:"happened_at"`
	ConnectionId   string      `json:"connection_id"`
	ContextData    interface{} `json:"context_data"`
	IvrData        interface{} `json:"ivr_data"`
	SkillGroupId   *string     `json:"skill_group_id"`
	CreatedAt      string      `json:"created_at"`
	CreatedBy      interface{} `json:"created_by"`
	CreatorType    string      `json:"creator_type"`
	UpdatedAt      string      `json:"updated_at"`
	UpdatedBy      interface{} `json:"updated_by"`
	UpdaterType    string      `json:"updater_type"`
	DeletedAt      interface{} `json:"deleted_at"`
	DeletedBy      interface{} `json:"deleted_by"`
	DeleterType    interface{} `json:"deleter_type"`
}

type CallSkillGroup struct {
	Id              string      `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Priority        int         `json:"priority"`
	PromptTimeout   int         `json:"prompt_timeout"`
	QueueTimeout    int         `json:"queue_timeout"`
	TelephoneNumber string      `json:"telephone_number"`
	CreatedAt       string      `json:"created_at"`
	CreatedBy       interface{} `json:"created_by"`
	CreatorType     string      `json:"creator_type"`
	UpdatedAt       string      `json:"updated_at"`
	UpdatedBy       string      `json:"updated_by"`
	UpdaterType     string      `json:"updater_type"`
	DeletedAt       interface{} `json:"deleted_at"`
	DeletedBy       interface{} `json:"deleted_by"`
	DeleterType     string      `json:"deleter_type"`
}

type ConnectionEndReason struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
	Cause       string `json:"cause"`
	Message     string `json:"message"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TelephoneRecording struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	CallId          string `json:"call_id"`
	RecordingTypeId string `json:"recording_type_id"`
	CreatedAt       string `json:"created_at"`
	CreatedBy       string `json:"created_by"`
	CreatorType     string `json:"creator_type"`
	UpdatedAt       string `json:"updated_at"`
	UpdatedBy       string `json:"updated_by"`
	UpdaterType     string `json:"updater_type"`
	DeletedAt       string `json:"deleted_at"`
	DeletedBy       string `json:"deleted_by"`
	DeleterType     string `json:"deleter_type"`
}
