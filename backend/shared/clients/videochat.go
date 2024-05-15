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
	VideoChatRecordingsEndpoint = "/api/recordings"

	VideoChatUriEnvKey = "VIDEO_CHAT_URI"
)

type VideoChat struct {
	Telephone
	url string
}

func NewVideoChat() *VideoChat {
	return &VideoChat{
		url: helpers.UrlParse(runtimebag.GetEnvString(
			VideoChatUriEnvKey,
			"http://video-chat:8080",
		)),
		Telephone: Telephone{
			uri: helpers.UrlParse(runtimebag.GetEnvString(
				VideoChatUriEnvKey,
				"http://video-chat:8080",
			)),
			searchModel: "video-call",
		},
	}
}

func (t *VideoChat) SetUri(uri string) *VideoChat {
	t.uri = helpers.UrlParse(uri)
	t.url = helpers.UrlParse(uri)
	return t
}

type UploadVideoFileRequest struct {
	File          *os.File
	RecordingType string
	SessionUuid   string
}

func (r *UploadVideoFileRequest) ToJson() []byte {
	var v []byte
	var err error
	if v, err = json.Marshal(r); err != nil {
		log.Panic(err.Error())
	}
	return v
}

func (r *UploadVideoFileRequest) ToJsonString() string {
	return string(r.ToJson())
}

type UploadVideoChatFileResponse struct {
	Response
	result *VideoChatRecording
}

func (r *UploadVideoChatFileResponse) Result() *VideoChatRecording {
	return r.result
}

func (t *VideoChat) UploadFile(ctx *assecoContext.RequestContext, request *UploadVideoFileRequest, headers map[string]string, timeout time.Duration) *UploadVideoChatFileResponse {
	url := fmt.Sprintf("%s%s", TrimSuffix(t.url, "/"), VideoChatRecordingsEndpoint)
	recordingType := request.RecordingType
	if recordingType == "" {
		recordingType = "09cf6fa2-d352-431f-b525-65c50cfc3049"
	}
	name := strings.TrimSuffix(request.File.Name(), filepath.Ext(request.File.Name()))
	postBody := &bytes.Buffer{}
	writer := multipart.NewWriter(postBody)
	err := writer.WriteField("name", name)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding name field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("video_call_id", request.SessionUuid)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding video_call_id field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("recording_type_id", recordingType)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("failed adding recording_type_id field to writer, error %s", err.Error()),
			},
			result: nil,
		}
	}
	part, err := writer.CreateFormFile("file", filepath.Base(request.File.Name()))
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("failed creating form data for file field, error %s", err.Error()),
			},
			result: nil,
		}
	}
	_, err = io.Copy(part, request.File)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("error copying src to dst file, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.Close()
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("error closing writer, error %s", err.Error()),
			},
			result: nil,
		}
	}

	clientRequest, err := http.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		return &UploadVideoChatFileResponse{
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
	client := &http.Client{}
	rCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	defer cancel()
	resp, err := client.Do(clientRequest.WithContext(rCtx))

	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("error reading body, error %s", err.Error()),
			},
			result: nil,
		}
	}

	if resp.StatusCode >= 300 {
		return &UploadVideoChatFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	var result *VideoChatRecording
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &UploadVideoChatFileResponse{
			Response: Response{
				error: fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
			result: nil,
		}
	}
	_ = resp.Body.Close()

	return &UploadVideoChatFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type VideoChatRecording struct {
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
