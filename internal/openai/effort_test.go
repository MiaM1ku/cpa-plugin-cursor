package openai

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveModel_collapses_effort_suffixes_and_prefers_explicit_effort(t *testing.T) {
	selection := ResolveModel("cursor/claude-4.5-sonnet-thinking-xhigh", "", false)
	require.Equal(t, "claude-4.5-sonnet", selection.ID)
	require.Equal(t, "xhigh", selection.Effort)

	selection = ResolveModel("gpt-5.4-high", "low", false)
	require.Equal(t, "gpt-5.4", selection.ID)
	require.Equal(t, "low", selection.Effort)

	selection = ResolveModel("composer-1m", "", false)
	require.Equal(t, "composer", selection.ID)
	require.True(t, selection.MaxMode)

	selection = ResolveModel("claude-fable-5-thinking-high", "", false)
	require.Equal(t, "claude-fable-5", selection.ID)
	require.Equal(t, "high", selection.Effort)

	selection = ResolveModel("claude-fable-5-thinking-xhigh", "medium", false)
	require.Equal(t, "claude-fable-5", selection.ID)
	require.Equal(t, "medium", selection.Effort)

	selection = ResolveModel("claude-fable-5-max", "", false)
	require.Equal(t, "claude-fable-5", selection.ID)
	require.True(t, selection.MaxMode)
	require.Empty(t, selection.Effort)

	selection = ResolveModel("claude-fable-5-thinking-max", "", false)
	require.Equal(t, "claude-fable-5", selection.ID)
	require.Equal(t, "high", selection.Effort)
	require.True(t, selection.MaxMode)

	selection = ResolveModel("claude-fable-5-max-thinking-xhigh", "", false)
	require.Equal(t, "claude-fable-5", selection.ID)
	require.Equal(t, "xhigh", selection.Effort)
	require.True(t, selection.MaxMode)
}

func Test_ParseChatRequest_reads_reasoning_effort(t *testing.T) {
	request, err := ParseChatRequest([]byte(`{
		"model":"cursor/claude-4.5-sonnet",
		"reasoning_effort":"high",
		"messages":[{"role":"user","content":"hello"}]
	}`))
	require.NoError(t, err)
	require.Equal(t, "claude-4.5-sonnet", request.Model)
	require.Equal(t, "high", request.Effort)
}
