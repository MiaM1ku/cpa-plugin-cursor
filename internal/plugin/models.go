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
	blocked := normalizedModelSet(disabled)
	filtered := make([]string, 0, len(models))
	for _, id := range models {
		normalized := normalizeModelID(id)
		if _, found := blocked[normalized]; !found && normalized != "" {
			filtered = append(filtered, id)
		}
	}
	return filtered
}

func collapseModels(ids []string) []string {
	type family struct {
		thinking   bool
		nothinking bool
	}
	families := make(map[string]*family, len(ids))
	order := make([]string, 0, len(ids))
	for _, id := range ids {
		if strings.HasSuffix(strings.ToLower(strings.TrimSpace(id)), "-fast") {
			continue
		}
		selection := openai.ResolveModel(id, "", false)
		if selection.ID == "" {
			continue
		}
		current := families[selection.ID]
		if current == nil {
			current = &family{}
			families[selection.ID] = current
			if selection.ID != "auto" {
				order = append(order, selection.ID)
			}
		}
		if selection.Thinking {
			current.thinking = true
		} else {
			current.nothinking = true
		}
	}
	ordered := make([]string, 0, len(order)+1)
	ordered = append(ordered, "auto")
	for _, familyID := range order {
		current := families[familyID]
		if current.thinking {
			ordered = append(ordered, familyID+"-thinking")
		}
		if current.nothinking {
			ordered = append(ordered, familyID)
		}
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
