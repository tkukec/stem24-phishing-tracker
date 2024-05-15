package clients

import (
	"bytes"
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
	"strconv"
	"strings"
	"time"
)

const (
	ContentUriEnvKey = "CONTENT_URI"
)

type Content struct {
	uri string
}

func NewContent() *Content {
	return &Content{uri: helpers.UrlParse(runtimebag.GetEnvString(
		ContentUriEnvKey,
		"http://content:8080/v1/content",
	))}
}

func (c *Content) SetUri(uri string) *Content {
	c.uri = helpers.UrlParse(uri)
	return c
}

type RepositoriesResponse struct {
	Response
	result []*Repository
}

func (r *RepositoriesResponse) Result() []*Repository {
	return r.result
}

func (c *Content) Repositories(ctx *assecoContext.RequestContext, headers map[string]string) *RepositoriesResponse {
	url := helpers.SanitizeUrl(c.uri, "repositories")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &RepositoriesResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))

	if err != nil {
		return &RepositoriesResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RepositoriesResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
			result: nil,
		}
	}
	if resp.StatusCode == 404 {
		return &RepositoriesResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &RepositoriesResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	var result []*Repository
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &RepositoriesResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
			result: nil,
		}
	}
	_ = resp.Body.Close()
	return &RepositoriesResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type FoldersRequest struct {
	Repo, Dir string
}

type FoldersResponse struct {
	Response
	result *Folder
}

func (r *FoldersResponse) Result() *Folder {
	return r.result
}

func (c *Content) Folder(ctx *assecoContext.RequestContext, request *FoldersRequest, headers map[string]string) *FoldersResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, request.Dir)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &FoldersResponse{
			Response: Response{
				error: fmt.Errorf("error creating request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	req.Header.Add(constants.XCorrelationID, ctx.XCorrelationID())
	for i, v := range headers {
		req.Header.Add(i, v)
	}
	resp, err := client.Do(req.WithContext(ctx.Context()))

	if err != nil {
		return &FoldersResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
			result: nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &FoldersResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode == 404 {
		return &FoldersResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &FoldersResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	var result *Folder
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &FoldersResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
			result: nil,
		}
	}
	_ = resp.Body.Close()

	return &FoldersResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type FolderMetadataResponse struct {
	Response
	result *Item
}

func (r *FolderMetadataResponse) Result() *Item {
	return r.result
}

func (c *Content) FolderMetadata(ctx *assecoContext.RequestContext, request *FoldersRequest, headers map[string]string) *FolderMetadataResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, request.Dir, "metadata")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &FolderMetadataResponse{
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
		return &FolderMetadataResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &FolderMetadataResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode == 404 {
		return &FolderMetadataResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &FolderMetadataResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
		}
	}

	var result *Item
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &FolderMetadataResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
		}
	}
	_ = resp.Body.Close()

	return &FolderMetadataResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type FileMetadataRequest struct {
	Repo, Location string
}

type FileMetadataResponse struct {
	Response
	result *FileMetadata
}

func (r *FileMetadataResponse) Result() *FileMetadata {
	return r.result
}

func (c *Content) FileMetadataByPath(ctx *assecoContext.RequestContext, request *FileMetadataRequest, headers map[string]string) *FileMetadataResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, request.Location, "metadata")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &FileMetadataResponse{
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
		return &FileMetadataResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &FileMetadataResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode == 404 {
		return &FileMetadataResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &FileMetadataResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
		}
	}

	var result *FileMetadata
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &FileMetadataResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
		}
	}
	_ = resp.Body.Close()

	return &FileMetadataResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type CreateFolderRequest struct {
	Repo          string `json:"repo"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	Kind          string `json:"kind"`
	FolderPurpose string `json:"folder-purpose"`
}

func (r *CreateFolderRequest) ToJson() []byte {
	var v []byte
	var err error
	if v, err = json.Marshal(r); err != nil {
		log.Panic(err.Error())
	}
	return v
}

func (r *CreateFolderRequest) ToJsonString() string {
	return string(r.ToJson())
}

func NewCreateFolderRequest(repo string, name string, path string) *CreateFolderRequest {
	return &CreateFolderRequest{Repo: repo, Name: name, Path: path}
}

type CreateFolderResponse struct {
	Response
	result *Item
}

func (r *CreateFolderResponse) Result() *Item {
	return r.result
}

func (c *Content) CreateFolder(ctx *assecoContext.RequestContext, request *CreateFolderRequest, headers map[string]string) *CreateFolderResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, "folders")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(request.ToJson()))
	if err != nil {
		return &CreateFolderResponse{
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
		return &CreateFolderResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &CreateFolderResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
		}
	}
	if resp.StatusCode == 404 {
		return &CreateFolderResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &CreateFolderResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
		}
	}

	var result *Item
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &CreateFolderResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
		}
	}
	_ = resp.Body.Close()

	return &CreateFolderResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

func (c *Content) CreatePath(ctx *assecoContext.RequestContext, request *CreateFolderRequest, headers map[string]string) *CreateFolderResponse {
	pathString := request.Path
	path := strings.Split(pathString, "/")
	current := "/"
	var response *CreateFolderResponse
	for _, item := range path {
		if item == "" {
			continue
		}
		log.Printf("creating %s ....", item)
		request.Name = item
		request.Path = current
		response = c.CreateFolder(ctx, request, headers)
		if response.Failed() && !okFailure(response.Status()) {
			return response

		}
		log.Printf("item %s already exists", item)
		if current == "/" {
			current = fmt.Sprintf("%s%s", current, item)
		} else {
			current = fmt.Sprintf("%s/%s", current, item)
		}
	}
	return response
}

func okFailure(status int) bool {
	return status == 409 || status == 440
}

type UploadContentFileRequest struct {
	//Repo repository to use
	Repo string `json:"repo"`
	//FolderId the id of the folder in which to upload file
	FolderId string `json:"-"`

	File              io.Reader
	Name              string `json:"name"`
	MediaType         string `json:"media-type"`
	FilingPurpose     string `json:"filing-purpose"`
	FilingCaseNumber  string `json:"filing-case-number"`
	OverwriteIfExists bool   `json:"overwrite-if-exists"`
	Extended          string `json:"extended"`
}

func NewUploadContentFileRequest(
	repo, folderId string,
	file io.Reader,
	name, mediaType, filingPurpose, filingCaseNumber string,
	overwriteIfExists bool,
	extended string,
) *UploadContentFileRequest {
	return &UploadContentFileRequest{
		Repo:              repo,
		FolderId:          folderId,
		File:              file,
		Name:              name,
		MediaType:         mediaType,
		FilingPurpose:     filingPurpose,
		FilingCaseNumber:  filingCaseNumber,
		OverwriteIfExists: overwriteIfExists,
		Extended:          extended,
	}
}

type UploadContentFileResponse struct {
	Response
	result *UploadedFile
}

func (r *UploadContentFileResponse) Result() *UploadedFile {
	return r.result
}

func (c *Content) UploadFile(ctx *assecoContext.RequestContext, request *UploadContentFileRequest, headers map[string]string) *UploadContentFileResponse {
	postBody := &bytes.Buffer{}
	writer := multipart.NewWriter(postBody)
	part, err := writer.CreateFormFile("content-stream", request.Name)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error creating form file, error %s", err.Error()),
			},
			result: nil,
		}
	}
	_, err = io.Copy(part, request.File)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error copying from src to dst file, error %s", err.Error()),
			},
			result: nil,
		}
	}

	err = writer.WriteField("name", request.Name)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field name, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("media-type", request.MediaType)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field media-type, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("filing-purpose", request.FilingPurpose)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field filing-purpose, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("filing-case-number", request.FilingCaseNumber)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field filing-case-number, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("overwrite-if-exists", strconv.FormatBool(request.OverwriteIfExists))
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field overwrite-if-exists, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.WriteField("extended", request.Extended)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error writing a field extended, error %s", err.Error()),
			},
			result: nil,
		}
	}
	err = writer.Close()
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: fmt.Errorf("error closing writer, error %s", err.Error()),
			},
			result: nil,
		}
	}

	url := helpers.SanitizeUrl(c.uri, request.Repo, "folders", request.FolderId)
	client := &http.Client{}
	clientRequest, err := http.NewRequest(http.MethodPost, url, postBody)
	if err != nil {
		return &UploadContentFileResponse{
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
	resp, err := client.Do(clientRequest.WithContext(ctx.Context()))

	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				error: err,
			},
			result: nil,
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error reading body, error %s", err.Error()),
			},
			result: nil,
		}
	}
	if resp.StatusCode == 404 {
		return &UploadContentFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
				responseError: mapError(body),
			},
			result: nil,
		}
	}
	if resp.StatusCode >= 300 {
		return &UploadContentFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	var result *UploadedFile
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &UploadContentFileResponse{
			Response: Response{
				status: resp.StatusCode,
				error:  fmt.Errorf("error un-marshaling body, error %s body %s", err.Error(), string(body)),
			},
			result: nil,
		}
	}
	_ = resp.Body.Close()

	return &UploadContentFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: result,
	}
}

type GetFileRequest struct {
	Repo     string `json:"repo"`
	Location string `json:"location"`
}

func NewGetFileRequest(repo string, location string) *GetFileRequest {
	return &GetFileRequest{Repo: repo, Location: location}
}

type GetFileResponse struct {
	Response
	result []byte
}

func (r *GetFileResponse) Result() []byte {
	return r.result
}

func (c *Content) GetFile(ctx *assecoContext.RequestContext, request *GetFileRequest, headers map[string]string) *GetFileResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, request.Location)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &GetFileResponse{
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
		return &GetFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &GetFileResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
				result: nil,
			}
		}
		if resp.StatusCode == 404 {
			return &GetFileResponse{
				Response: Response{
					status:        resp.StatusCode,
					error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
					responseError: mapError(body),
				},
				result: nil,
			}
		}
		return &GetFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: nil,
		}
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return &GetFileResponse{
			Response: Response{
				error:  err,
				status: resp.StatusCode,
			},
		}
	}

	return &GetFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: buf.Bytes(),
	}
}

type DeleteFileRequest struct {
	Repo string `json:"repo"`
	ID   string `json:"id"`
}

type DeleteFileResponse struct {
	Response
	result bool
}

func (r *DeleteFileResponse) Result() bool {
	return r.result
}

func (c *Content) DeleteFile(ctx *assecoContext.RequestContext, request *DeleteFileRequest, headers map[string]string) *DeleteFileResponse {
	url := helpers.SanitizeUrl(c.uri, request.Repo, "documents", request.ID)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return &DeleteFileResponse{
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
		return &DeleteFileResponse{
			Response: Response{
				error: fmt.Errorf("error executing request, error %s", err.Error()),
			},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &DeleteFileResponse{
				Response: Response{
					status: resp.StatusCode,
					error:  fmt.Errorf("error reading body, error %s", err.Error()),
				},
				result: false,
			}
		}
		if resp.StatusCode == 404 {
			return &DeleteFileResponse{
				Response: Response{
					status:        resp.StatusCode,
					error:         fmt.Errorf("uri %s or its content deleted, moved or written incrorrectly", url),
					responseError: mapError(body),
				},
				result: false,
			}
		}
		return &DeleteFileResponse{
			Response: Response{
				status:        resp.StatusCode,
				responseError: mapError(body),
			},
			result: false,
		}
	}

	return &DeleteFileResponse{
		Response: Response{
			status: resp.StatusCode,
		},
		result: true,
	}
}

type Repository struct {
	ID   string `json:"repository-id"`
	Name string `json:"repository-name"`
}

type Folder struct {
	TotalCount int     `json:"total-count"`
	PageSize   int     `json:"page-size"`
	Page       int     `json:"page"`
	TotalPages int     `json:"total-pages"`
	Items      []*Item `json:"items"`
}

type Item struct {
	FolderPurpose string    `json:"folder-purpose"`
	Id            string    `json:"id"`
	ChangedOn     time.Time `json:"changed-on"`
	CreatedOn     time.Time `json:"created-on"`
	CreatedBy     string    `json:"created-by"`
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Kind          string    `json:"kind"`
	Extended      struct{}  `json:"extended"`
}

type FileMetadata struct {
	MediaType        string      `json:"media-type"`
	FilingPurpose    string      `json:"filing-purpose"`
	FilingCaseNumber string      `json:"filing-case-number"`
	SizeInBytes      int         `json:"size-in-bytes"`
	Id               string      `json:"id"`
	ChangedOn        time.Time   `json:"changed-on"`
	CreatedOn        time.Time   `json:"created-on"`
	CreatedBy        string      `json:"created-by"`
	Name             string      `json:"name"`
	Path             string      `json:"path"`
	Kind             string      `json:"kind"`
	Extended         interface{} `json:"extended"`
}

type UploadedFile struct {
	MediaType        string                 `json:"media-type"`
	FilingPurpose    string                 `json:"filing-purpose"`
	FilingCaseNumber string                 `json:"filing-case-number"`
	SizeInBytes      int                    `json:"size-in-bytes"`
	Id               string                 `json:"id"`
	ChangedOn        time.Time              `json:"changed-on"`
	CreatedOn        time.Time              `json:"created-on"`
	CreatedBy        string                 `json:"created-by"`
	Name             string                 `json:"name"`
	Path             string                 `json:"path"`
	Kind             string                 `json:"kind"`
	Extended         map[string]interface{} `json:"extended"`
}
