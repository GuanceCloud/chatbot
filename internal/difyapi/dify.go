package difyapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	EnvDifyAPIServer = "ENV_DIFY_API_SERVER"
	EnvDifyAPIKey    = "ENV_DIFY_API_KEY"
	difyAPIServer    = "http://localhost/v1"

	datasetsAPI     = "/datasets"
	creatByTextAPI  = "/datasets/%s/document/create_by_text"
	updateByTextAPI = "/datasets/%s/documents/%s/update_by_text"
	// createByFileAPI   = "/datasets/%s/document/create_by_file"
	// updateByFileAPI   = "/datasets/%s/documents/%s/update_by_file"
	// indexingStatusAPI = "/datasets/%s/documents/%s/indexing-status"
	// deleteDocumentAPI = "/datasets/%s/documents/%s"
	// listDocumentAPI   = "/datasets/%s/documents"

	GET  = "GET"
	POST = "POST"
)

func StatusMapping() map[int]map[string]string {
	// TODO: new error
	return map[int]map[string]string{
		400: {
			"no_file_uploaded":          "Please upload your file.",
			"too_many_files":            "Only one file is allowed.",
			"high_quality_dataset_only": "Current operation only supports 'high-quality' datasets.",
			"dataset_not_initialized":   "The dataset is still being initialized or indexing. Please wait a moment.",
			"invalid_action":            "Invalid action.",
			"document_already_finished": "The document has been processed. Please refresh the page or go to the document details.",
			"document_indexing":         "The document is being processed and cannot be edited.",
			"invalid_metadata":          "The metadata content is incorrect. Please check and verify.",
		},
		403: {
			"archived_document_immutable": "The archived document is not editable.",
		},
		409: {
			"dataset_name_duplicate": "The dataset name already exists. Please modify your dataset name.",
		},
		413: {
			"file_too_large": "File size exceeded.",
		},
		415: {
			"unsupported_file_type": "File type not allowed.",
		},
	}
}

type Datasets struct {
	apiServer string
	apiKey    string

	cli *http.Client
}

func GetAPI() (string, string) {
	apiKey := os.Getenv(EnvDifyAPIKey)
	apiServer := os.Getenv(EnvDifyAPIServer)
	if apiServer == "" {
		apiServer = difyAPIServer
	}
	return apiServer, apiKey
}

func NewDataset(server, key string) *Datasets {
	return &Datasets{
		apiServer: server,
		apiKey:    key,
		cli: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    100,
				IdleConnTimeout: 90 * time.Second,
			}},
	}
}

func (ag *Datasets) req(method, url string, body any) (int, any, error) {
	v, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(v))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", "Bearer "+ag.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := ag.cli.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != 200 {
		return resp.StatusCode, buf, fmt.Errorf("%s", string(buf))
	}
	var mp map[string]any
	if err := json.Unmarshal(buf, &mp); err != nil {
		return resp.StatusCode, buf, err
	}

	return resp.StatusCode, mp, nil
}

func (ag *Datasets) ListDatasets(page, limit int) (map[string]any, error) {
	u, err := url.JoinPath(ag.apiServer, datasetsAPI)
	if err != nil {
		return nil, err
	}
	u += fmt.Sprintf(
		"?limit=%d&page=%d", limit, page)

	_, val, err := ag.req(GET, u, nil)
	if err != nil {
		return nil, err
	}

	if val, ok := val.(map[string]any); ok {
		return val, nil
	}

	return nil, fmt.Errorf("failed")
}

func (ds *Datasets) CreateDocByText(datasetID, indexTech, name, text string) (map[string]any, error) {
	baseURL, err := url.JoinPath(ds.apiServer, fmt.Sprintf(creatByTextAPI, datasetID))
	if err != nil {
		return nil, err
	}
	if indexTech == "" {
		indexTech = "high_quality"
	}
	body := map[string]any{
		"name":               name,
		"text":               text,
		"indexing_technique": indexTech,
		"process_rule": map[string]string{
			"mode": "automatic",
		},
	}

	_, val, err := ds.req(POST, baseURL, body)
	if err != nil {
		return nil, err
	}
	if val, ok := val.(map[string]any); ok {
		return val, nil
	}

	return nil, fmt.Errorf("failed")
}

func (ag *Datasets) UpdateDocByText(datasetID, docID, name, text string) (map[string]any, error) {
	baseURL := filepath.Join(ag.apiServer, fmt.Sprintf(updateByTextAPI, datasetID, docID))
	body := map[string]any{}
	if name != "" {
		body["name"] = name
	}
	if text != "" {
		body["text"] = text
	}
	_, val, err := ag.req(POST, baseURL, body)
	if err != nil {
		return nil, err
	}
	if val, ok := val.(map[string]any); ok {
		return val, nil
	}

	return nil, fmt.Errorf("failed")
}
