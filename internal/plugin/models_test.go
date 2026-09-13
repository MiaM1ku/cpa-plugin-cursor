package plugin

import (
	"context"
	"encoding/json"
	"testing"

	"cursorplugin/internal/cursorapi"
	"cursorplugin/internal/cursorauth"
	"cursorplugin/internal/cursorproto"

	"github.com/stretchr/testify/require"
)

func Test_Handler_ModelsForAuth_hides_models_disabled_by_cursor_plugin(t *testing.T) {
	handler := NewHandler(Dependencies{Cursor: fakeModelCursorClient{models: []string{"auto", "claude-4-sonnet", "gpt-5"}}})
	storage, err := cursorauth.MarshalCredentials(cursorauth.Credentials{
		AccessToken:    "access",
		RefreshToken:   "refresh",
		Type:           "cursor",
		DisabledModels: []string{"cursor/gpt-5"},
	})
	require.NoError(t, err)
	rawRequest, err := json.Marshal(authModelRequest{StorageJSON: storage})
	require.NoError(t, err)

	result, err := handler.modelsForAuth(context.Background(), rawRequest)

	require.NoError(t, err)
	rawResult, err := json.Marshal(result)
	require.NoError(t, err)
	require.Contains(t, string(rawResult), `"ID":"cursor/auto"`)
	require.Contains(t, string(rawResult), `"SupportedInputModalities":["text","image"]`)
	require.Contains(t, string(rawResult), `"SupportedOutputModalities":["text","image"]`)
	require.Contains(t, string(rawResult), `"ID":"cursor/claude-4-sonnet"`)
	require.NotContains(t, string(rawResult), `"ID":"cursor/gpt-5"`)
}

type fakeModelCursorClient struct {
	models []string
}

func (client fakeModelCursorClient) Run(context.Context, cursorapi.RunInput, func(cursorproto.ServerEvent) error) (cursorapi.RunResult, error) {
	return cursorapi.RunResult{}, nil
}

func (client fakeModelCursorClient) DiscoverModels(context.Context, string) ([]string, error) {
	return append([]string(nil), client.models...), nil
}

func Test_collapseModels_keeps_one_id_per_family(t *testing.T) {
	collapsed := collapseModels([]string{
		"claude-4.5-sonnet-thinking-xhigh",
		"claude-4.5-sonnet-thinking",
		"claude-4.5-sonnet",
		"gpt-5.4-high",
		"auto",
	})
	require.Equal(t, []string{"auto", "claude-4.5-sonnet-thinking", "claude-4.5-sonnet", "gpt-5.4"}, collapsed)
}

func Test_collapseModels_collapses_fable_thinking_effort_suffixes(t *testing.T) {
	collapsed := collapseModels([]string{
		"claude-fable-5-thinking-high",
		"claude-fable-5-thinking-xhigh",
		"claude-fable-5",
		"claude-fable-5-thinking",
	})
	require.Equal(t, []string{"auto", "claude-fable-5-thinking", "claude-fable-5"}, collapsed)
}

func Test_collapseModels_collapses_max_mode_and_effort_variants(t *testing.T) {
	collapsed := collapseModels([]string{
		"claude-fable-5",
		"claude-fable-5-thinking",
		"claude-fable-5-thinking-high",
		"claude-fable-5-max",
		"claude-fable-5-thinking-max",
	})
	require.Equal(t, []string{"auto", "claude-fable-5-thinking", "claude-fable-5"}, collapsed)
}

func Test_collapseModels_skips_fast_variants(t *testing.T) {
	collapsed := collapseModels([]string{
		"gpt-5.2-high-fast",
		"gpt-5.2-medium-fast",
		"gpt-5.2-high",
	})
	require.Equal(t, []string{"auto", "gpt-5.2"}, collapsed)
}

func Test_filterDisabledModels_matches_collapsed_effort_ids(t *testing.T) {
	filtered := filterDisabledModels(
		[]string{"auto", "claude-fable-5", "gpt-5"},
		[]string{"claude-fable-5-thinking-high"},
	)
	require.Equal(t, []string{"auto", "gpt-5"}, filtered)
}
