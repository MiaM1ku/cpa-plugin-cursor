package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cursorplugin/internal/cursorauth"
	"cursorplugin/internal/openai"
)

type authModelRequest struct {
	StorageJSON []byte `json:"StorageJSON"`
}

func (handler *Handler) modelsForAuth(ctx context.Context, raw []byte) (any, error) {
	var request authModelRequest
	if err := json.Unmarshal(raw, &request); err != nil {
		return nil, fmt.Errorf("decode auth model request: %w", err)
	}
	credentials, err := cursorauth.ParseCredentials(request.StorageJSON)
	if err != nil {
		return nil, err
	}
	models, err := handler.cursor.DiscoverModels(ctx, credentials.AccessToken)
	if err != nil {
		return nil, err
	}
	return modelResponse(filterDisabledModels(collapseModels(models), credentials.DisabledModels)), nil
}

func filterDisabledModels(models, disabled []string) []string {
	if len(disabled) == 0 {
		return append([]string(nil), models...)
	}
	blocked := make(map[string]struct{}, len(disabled))
	for _, id := range disabled {
		normalized := strings.TrimPrefix(strings.TrimSpace(id), "cursor/")
		if normalized != "" {
			blocked[normalized] = struct{}{}
		}
	}
	filtered := make([]string, 0, len(models))
	for _, id := range models {
		normalized := strings.TrimPrefix(strings.TrimSpace(id), "cursor/")
		if _, found := blocked[normalized]; !found && normalized != "" {
			filtered = append(filtered, normalized)
		}
	}
	return filtered
}

func collapseModels(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	collapsed := make([]string, 0, len(ids))
	for _, id := range ids {
		base := openai.CollapseModelID(id)
		if base == "" {
			continue
		}
		if _, exists := seen[base]; exists {
			continue
		}
		seen[base] = struct{}{}
		collapsed = append(collapsed, base)
	}
	ordered := make([]string, 0, len(collapsed)+1)
	ordered = append(ordered, "auto")
	for _, id := range collapsed {
		if id == "auto" {
			continue
		}
		ordered = append(ordered, id)
	}
	return ordered
}

func modelResponse(ids []string) any {
	models := make([]modelInfo, 0, len(ids))
	for _, rawID := range ids {
		id := strings.TrimSpace(rawID)
		if id == "" {
			continue
		}
		models = append(models, modelInfo{
			ID:                         "cursor/" + id,
			Object:                     "model",
			OwnedBy:                    "cursor",
			DisplayName:                "Cursor " + id,
			SupportedGenerationMethods: []string{"chat"},
			SupportedInputModalities:   []string{"text", "image"},
			SupportedOutputModalities:  []string{"text", "image"},
			ContextLength:              200000,
			MaxCompletionTokens:        32768,
			UserDefined:                true,
		})
	}
	return struct {
		Provider string      `json:"Provider"`
		Models   []modelInfo `json:"Models"`
	}{Provider: "cursor", Models: models}
}
