package openai

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveWireID_uses_thinking_suffix_and_effort(t *testing.T) {
	catalog := []string{
		"claude-fable-5-high",
		"claude-fable-5-medium",
		"claude-fable-5-thinking-high",
		"claude-fable-5-thinking-medium",
		"claude-fable-5-thinking-high-fast",
		"gpt-5.4-high",
		"gpt-5.4-medium",
		"auto",
	}

	require.Equal(t, "claude-fable-5-thinking-high", ResolveWireID(ModelSelection{ID: "claude-fable-5", Effort: "high", Thinking: true}, catalog))
	require.Equal(t, "claude-fable-5-high", ResolveWireID(ModelSelection{ID: "claude-fable-5", Effort: "high"}, catalog))
	require.Equal(t, "claude-fable-5-thinking-medium", ResolveWireID(ModelSelection{ID: "claude-fable-5", Thinking: true}, catalog))
	require.Equal(t, "claude-fable-5-medium", ResolveWireID(ModelSelection{ID: "claude-fable-5"}, catalog))
	require.Equal(t, "gpt-5.4-medium", ResolveWireID(ModelSelection{ID: "gpt-5.4"}, catalog))
	require.Equal(t, "auto", ResolveWireID(ModelSelection{ID: "auto"}, catalog))
}

func Test_ResolveWireID_ignores_fast_catalog_rows(t *testing.T) {
	catalog := []string{"gpt-5.2-high-fast", "gpt-5.2-medium-fast"}
	require.Equal(t, "gpt-5.2-medium", ResolveWireID(ModelSelection{ID: "gpt-5.2"}, catalog))
}
