package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cursorplugin/internal/cursorauth"
	"cursorplugin/internal/cursorusage"
)

const maxHTTPResponseBytes = 2 << 20

type executorHTTPRequest struct {
	AuthID       string      `json:"AuthID"`
	AuthProvider string      `json:"AuthProvider"`
	Method       string      `json:"Method"`
	URL          string      `json:"URL"`
	Headers      http.Header `json:"Headers"`
	Body         []byte      `json:"Body"`
	StorageJSON  []byte      `json:"StorageJSON"`
}

type executorHTTPResponse struct {
	StatusCode int         `json:"StatusCode"`
	Headers    http.Header `json:"Headers"`
	Body       []byte      `json:"Body"`
}

func (handler *Handler) httpRequest(ctx context.Context, raw []byte) (any, error) {
	var request executorHTTPRequest
	if err := json.Unmarshal(raw, &request); err != nil {
		return nil, fmt.Errorf("decode executor HTTP request: %w", err)
	}
	credentials, err := cursorauth.ParseCredentials(request.StorageJSON)
	if err != nil {
		return nil, err
	}
	method := strings.ToUpper(strings.TrimSpace(request.Method))
	if method == "" {
		method = http.MethodGet
	}
	target := strings.TrimSpace(request.URL)
	if target == "" {
		return nil, fmt.Errorf("Cursor HTTP passthrough requires a URL")
	}
	httpRequest, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(request.Body))
	if err != nil {
		return nil, fmt.Errorf("create Cursor HTTP request: %w", err)
	}
	httpRequest.Header = request.Headers.Clone()
	if httpRequest.Header == nil {
		httpRequest.Header = make(http.Header)
	}
	if cursorusage.IsDashboardURL(httpRequest.URL) {
		cursorusage.ApplyDashboardAuth(httpRequest, credentials.DashboardAccountID(), credentials.AccessToken)
	} else if httpRequest.Header.Get("Authorization") == "" {
		httpRequest.Header.Set("Authorization", "Bearer "+credentials.AccessToken)
	}
	client := &http.Client{Timeout: 60 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("Cursor HTTP passthrough failed: %w", err)
	}
	body, readErr := io.ReadAll(io.LimitReader(response.Body, maxHTTPResponseBytes+1))
	closeErr := response.Body.Close()
	if err := joinHTTPErrors(readErr, closeErr); err != nil {
		return nil, fmt.Errorf("read Cursor HTTP response: %w", err)
	}
	if len(body) > maxHTTPResponseBytes {
		return nil, fmt.Errorf("Cursor HTTP response exceeds 2 MiB")
	}
	return executorHTTPResponse{StatusCode: response.StatusCode, Headers: response.Header.Clone(), Body: body}, nil
}

func joinHTTPErrors(errs ...error) error {
	var first error
	for _, err := range errs {
		if err == nil {
			continue
		}
		if first == nil {
			first = err
			continue
		}
		first = fmt.Errorf("%v; %w", first, err)
	}
	return first
}
