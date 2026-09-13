package plugin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"cursorplugin/internal/cursorauth"
	"cursorplugin/internal/cursorusage"

	"github.com/stretchr/testify/require"
)

func Test_Handler_HTTPRequest_uses_dashboard_session_cookie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method)
		require.Equal(t, "/api/usage-summary", request.URL.Path)
		require.Equal(t, "https://cursor.com", request.Header.Get("origin"))
		cookie, err := url.QueryUnescape(request.Header.Get("cookie"))
		require.NoError(t, err)
		require.Equal(t, "WorkosCursorSessionToken=auth0|user-1::access-token", cookie)
		require.Empty(t, request.Header.Get("Authorization"))
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	credentials, err := cursorauth.MarshalCredentials(cursorauth.Credentials{
		AccessToken: "access-token", RefreshToken: "refresh-token", AccountID: "auth0|user-1", Type: "cursor",
	})
	require.NoError(t, err)
	raw, err := json.Marshal(executorHTTPRequest{
		Method: http.MethodGet, URL: server.URL + "/api/usage-summary",
		Headers:     http.Header{"Authorization": []string{"Bearer $TOKEN$"}},
		StorageJSON: credentials,
	})
	require.NoError(t, err)
	handler := NewHandler(Dependencies{})

	result, callErr := handler.httpRequest(context.Background(), raw)
	require.NoError(t, callErr)
	response := result.(executorHTTPResponse)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, string(response.Body), `"ok":true`)
}

func Test_Handler_HTTPRequest_translates_codex_wham_usage(t *testing.T) {
	credentials, err := cursorauth.MarshalCredentials(cursorauth.Credentials{
		AccessToken: "access-token", RefreshToken: "refresh-token", AccountID: "auth0|user-1", Type: "cursor",
	})
	require.NoError(t, err)
	handler := NewHandler(Dependencies{Usage: stubUsage{snapshot: cursorusage.Snapshot{
		MembershipType:  "pro",
		AutoPercentUsed: 12.5,
		APIPercentUsed:  80,
	}}})
	raw, err := json.Marshal(executorHTTPRequest{
		Method: http.MethodGet, URL: "https://chatgpt.com/backend-api/wham/usage",
		StorageJSON: credentials,
	})
	require.NoError(t, err)

	result, callErr := handler.httpRequest(context.Background(), raw)
	require.NoError(t, callErr)
	response := result.(executorHTTPResponse)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Contains(t, string(response.Body), `"plan_type":"pro"`)
	require.Contains(t, string(response.Body), `"primary_window"`)
	require.Contains(t, string(response.Body), `"used_percent":12.5`)
	require.Contains(t, string(response.Body), `"used_percent":80`)
}
