package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
)

const (
	PathEndpoint                               = "/api/path"
	MoveEndpoint                               = "/api/move"
	PathMetaEndpoint                           = "/api/path/meta"
	TransferEndpoint                           = "/api/transfer"
	ArchiveEndpoint                            = "/api/archive"
	StreamEndpoint                             = "/api/file"
	UnTransferredEndpoint                      = "/api/un-transferred"
	TransferUnTransferredEndpoint              = "/api/un-transferred/transfer"
	AppendAndTransferRelatedMediaFilesEndpoint = "/api/append"

	PathQueryKey    = "path"
	ForceQueryKey   = "force"
	RangeKey        = "Range"
	BytesKey        = "bytes"
	CreatePathQuery = "create_path"
)

type MediaFiles struct {
}

func NewMediaFiles() *MediaFiles {
	return &MediaFiles{}
}

func (m *MediaFiles) NewRequest(mediaServerUri string, method string) *RequestBuilder {
	return NewRequestBuilder(mediaServerUri, method)
}

type MediaServerFileError struct {
	Error string `json:"error"`
}

type CreatePathRequest struct {
	Location   string `json:"path" binding:"required"`
	Permission string `json:"permission"`
}

func (r *CreatePathRequest) ToJson() []byte {
	v, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return v
}

func (r *CreatePathRequest) ToJsonString() string {
	return string(r.ToJson())
}

type CreatePathResponse struct {
	Response
	result string
}

func (r *CreatePathResponse) Result() string {
	return r.result
}

func (m *MediaFiles) CreatePath(ctx *assecoContext.RequestContext, uri string, headers map[string]string, request CreatePathRequest) *CreatePathResponse {
	url := helpers.SanitizeUrl(uri, PathEndpoint)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.ToJson()))
	if err != nil {
		return &CreatePathResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &CreatePathResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &CreatePathResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &CreatePathResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &CreatePathResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}
	return &CreatePathResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: string(body),
	}

}

type PathMetaResponse struct {
	Response
	result DirectoryMeta
}

func (r *PathMetaResponse) Result() DirectoryMeta {
	return r.result
}

type DirectoryMeta struct {
	Location string          `json:"location"`
	Items    []DirectoryItem `json:"items"`
}

type DirectoryItem struct {
	Name       string `json:"name"`
	IsDir      bool   `json:"is_dir"`
	Permission string `json:"permission"`
	Size       string `json:"size"`
	SizeH      string `json:"size_human_readable"`
}

func BuildPathMetaQuery(uri, path string) string {
	return fmt.Sprintf("%s?%s", helpers.SanitizeUrl(uri, PathMetaEndpoint), composeQuery(PathQueryKey, path))
}

func BuildDeletePathQuery(uri, path, force string) string {
	return fmt.Sprintf("%s?%s&%s", helpers.SanitizeUrl(uri, PathEndpoint), composeQuery(PathQueryKey, path), composeQuery(ForceQueryKey, force))
}

func composeQuery(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

func (m *MediaFiles) PathMeta(ctx *assecoContext.RequestContext, uri string, headers map[string]string, path string) *PathMetaResponse {
	url := BuildPathMetaQuery(uri, path)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &PathMetaResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &PathMetaResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PathMetaResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &PathMetaResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &PathMetaResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}
	var result DirectoryMeta
	if err = json.Unmarshal(body, &result); err != nil {
		return &PathMetaResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}

	return &PathMetaResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type DeletePathResponse struct {
	Response
	result string
}

func (r *DeletePathResponse) Result() string {
	return r.result
}

func (m *MediaFiles) DeletePath(ctx *assecoContext.RequestContext, uri string, headers map[string]string, path string, force bool) *DeletePathResponse {
	url := BuildDeletePathQuery(uri, path, strconv.FormatBool(force))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return &DeletePathResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &DeletePathResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &DeletePathResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &DeletePathResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &DeletePathResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}

	return &DeletePathResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: string(body),
	}
}

type SendMediaToLiveRequest struct {
	Location string `json:"path" binding:"required"`
	Dir      string `json:"directory" binding:"required"`
	Purpose  string `json:"purpose" binding:"required"`
	Uuid     string `json:"session_uuid" binding:"required"`
	Remove   bool   `json:"remove"`
}

func (r *SendMediaToLiveRequest) ToJson() []byte {
	v, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return v
}

func (r *SendMediaToLiveRequest) ToJsonString() string {
	return string(r.ToJson())
}

type ContentFileData struct {
	CaseNumber string `json:"case_number"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Id         string `json:"id"`
}

type SendMediaToLiveResponse struct {
	Response
	result ContentFileData
}

func (r SendMediaToLiveResponse) Result() ContentFileData {
	return r.result
}

func (m *MediaFiles) TransferMediaToLiveContent(ctx *assecoContext.RequestContext, uri string, headers map[string]string, request SendMediaToLiveRequest) *SendMediaToLiveResponse {
	url := helpers.SanitizeUrl(uri, TransferEndpoint)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.ToJson()))
	if err != nil {
		return &SendMediaToLiveResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &SendMediaToLiveResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &SendMediaToLiveResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &SendMediaToLiveResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &SendMediaToLiveResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}
	var result ContentFileData
	if err = json.Unmarshal(body, &result); err != nil {
		return &SendMediaToLiveResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}

	return &SendMediaToLiveResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type ArchiveRequest struct {
	WhatToArchive  string `json:"path"`
	WhereToArchive string `json:"location"`
}

func (r *ArchiveRequest) ToJson() []byte {
	v, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return v
}

func (r *ArchiveRequest) ToJsonString() string {
	return string(r.ToJson())
}

type ArchiveResponse struct {
	Response
	result string
}

func (r *ArchiveResponse) Result() string {
	return r.result
}

func (m *MediaFiles) Archive(ctx *assecoContext.RequestContext, uri string, headers map[string]string, request ArchiveRequest) *ArchiveResponse {
	url := helpers.SanitizeUrl(uri, ArchiveEndpoint)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.ToJson()))
	if err != nil {
		return &ArchiveResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &ArchiveResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ArchiveResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &ArchiveResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &ArchiveResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}

	return &ArchiveResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: string(body),
	}
}

type StreamMediaFileRequest struct {
	Range string
	Path  string
}

type StreamMediaFileResponse struct {
	Response
	Result *http.Response
}

func (m *MediaFiles) StreamMediaFile(
	ctx *assecoContext.RequestContext, uri string, headers map[string]string, request StreamMediaFileRequest,
) *StreamMediaFileResponse {
	url := fmt.Sprintf("%s?%s", helpers.SanitizeUrl(uri, StreamEndpoint), composeQuery(PathQueryKey, request.Path))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &StreamMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}

	if request.Range != "" {
		req.Header.Add(RangeKey, fmt.Sprintf("%s=%s", BytesKey, request.Range))
	}
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &StreamMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &StreamMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &StreamMediaFileResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &StreamMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}

	return &StreamMediaFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		Result: resp,
	}
}

type UploadMediaFileRequest struct {
	Path                    string
	CreatePathIfDoesntExist bool
	File                    io.Reader
	Name                    string
}

type UploadMediaFileResponse struct {
	Response
	result string
}

func (m *MediaFiles) UploadMediaFile(
	ctx *assecoContext.RequestContext, uri string, headers map[string]string, request UploadMediaFileRequest,
) *UploadMediaFileResponse {
	postBody := &bytes.Buffer{}
	writer := multipart.NewWriter(postBody)
	part, err := writer.CreateFormFile("file", request.Name)
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating form file, error %s", err.Error()),
			},
		}
	}
	_, err = io.Copy(part, request.File)
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error copying from src to dst file, error %s", err.Error()),
			},
		}
	}

	err = writer.Close()
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error closing writer, error %s", err.Error()),
			},
		}
	}

	createPath := "false"
	if request.CreatePathIfDoesntExist {
		createPath = "true"
	}
	url := fmt.Sprintf("%s?%s&%s", helpers.SanitizeUrl(uri, PathEndpoint), composeQuery(PathQueryKey, request.Path), composeQuery(CreatePathQuery, createPath))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, postBody)
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UploadMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &UploadMediaFileResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &UploadMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}

	return &UploadMediaFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: string(body),
	}
}

type DownloadMediaFileRequest struct {
	Path string
}

type DownloadedFileData struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type DownloadMediaFileResponse struct {
	Response
	Result DownloadedFileData
}

func (m *MediaFiles) DownloadMediaFile(
	ctx *assecoContext.RequestContext, uri string, headers map[string]string, request DownloadMediaFileRequest,
) *DownloadMediaFileResponse {
	url := fmt.Sprintf("%s?%s", helpers.SanitizeUrl(uri, PathEndpoint), composeQuery(PathQueryKey, request.Path))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &DownloadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}

	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &DownloadMediaFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &DownloadMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}

	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &DownloadMediaFileResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &DownloadMediaFileResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}

	return &DownloadMediaFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		Result: DownloadedFileData{
			Name:    path.Base(request.Path),
			Content: body,
		},
	}
}

type MoveRequest struct {
	Targets map[string]string `json:"files"`
}

func (r *MoveRequest) ToJson() []byte {
	v, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return v
}

func (r *MoveRequest) ToJsonString() string {
	return string(r.ToJson())
}

type MoveItemsResponse struct {
	Response
	result map[string]interface{}
}

func (r *MoveItemsResponse) Result() map[string]interface{} {
	return r.result
}

func (m *MediaFiles) MoveItems(ctx *assecoContext.RequestContext, uri string, headers map[string]string, request MoveRequest) *MoveItemsResponse {
	url := helpers.SanitizeUrl(uri, MoveEndpoint)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.ToJson()))
	if err != nil {
		return &MoveItemsResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))
	if err != nil {
		return &MoveItemsResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &MoveItemsResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode >= 300 {
		var mediaError MediaServerFileError
		if err = json.Unmarshal(body, &mediaError); err != nil {
			return &MoveItemsResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
			}
		}

		return &MoveItemsResponse{
			Response: Response{
				status: resp.StatusCode,
				responseError: &ApiError{
					Message: mediaError.Error,
					Data:    mediaError,
					Code:    resp.Status,
				},
			},
		}
	}
	var response map[string]interface{}
	if err = json.Unmarshal(body, &response); err != nil {
		return &MoveItemsResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	return &MoveItemsResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: response,
	}

}

type FindUnTransferredResponse struct {
	Response
	result *DirectoryMeta
}

func (r *FindUnTransferredResponse) Result() *DirectoryMeta {
	return r.result
}

func (m *MediaFiles) FindUnTransferred(ctx *assecoContext.RequestContext, uri string, headers map[string]string, path string) *FindUnTransferredResponse {
	var result *DirectoryMeta
	req := m.NewRequest(uri, http.MethodGet).
		AddPath(UnTransferredEndpoint).
		AddQuery("path", path).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &FindUnTransferredResponse{
		Response: req.Response(),
		result:   result,
	}
}

type TransferUnTransferredResponse struct {
	Response
	result *TransferUnTransferredResponseData
}

func (r *TransferUnTransferredResponse) Result() *TransferUnTransferredResponseData {
	return r.result
}

type TransferUnTransferredRequest struct {
	Location string `json:"path" binding:"required"`
	Purpose  string `json:"purpose" binding:"required"`
	//limit the number of recordings you wish to transfer (FIFO system)
	SendCount int `json:"count"`
	//delete the recordings after they are transferred
	Remove bool `json:"delete"`
}

type TransferUnTransferredResponseData struct {
	Message   string          `json:"message"`
	ProcessId string          `json:"process_id"`
	Count     int             `json:"count"`
	Items     []DirectoryItem `json:"items"`
}

func (m *MediaFiles) TransferUnTransferred(ctx *assecoContext.RequestContext, uri string, headers map[string]string, request TransferUnTransferredRequest) *TransferUnTransferredResponse {
	var result *TransferUnTransferredResponseData
	req := m.NewRequest(uri, http.MethodPost).
		AddPath(TransferUnTransferredEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &TransferUnTransferredResponse{
		Response: req.Response(),
		result:   result,
	}
}

type AppendAndTransferRelatedMediaFilesRequest struct {
	FilePath         string `json:"path" binding:"required"`
	DelAfterTransfer bool   `json:"delete_after_transfer" binding:"required"`
}

type AppendAndTransferRelatedMediaFilesResponse struct {
	Response
	result string
}

func (m *MediaFiles) AppendAndTransferRelatedMediaFiles(
	ctx *assecoContext.RequestContext, uri string, headers map[string]string, request AppendAndTransferRelatedMediaFilesRequest,
) *AppendAndTransferRelatedMediaFilesResponse {
	var result string
	req := m.NewRequest(uri, http.MethodPost).
		AddPath(AppendAndTransferRelatedMediaFilesEndpoint).
		SetBody(request).
		SetHeaders(headers).
		Bind(&result).
		Print().
		Run(ctx)

	return &AppendAndTransferRelatedMediaFilesResponse{
		Response: req.Response(),
		result:   result,
	}
}
