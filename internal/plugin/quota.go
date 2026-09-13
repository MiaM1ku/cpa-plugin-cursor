package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cursorplugin/internal/cursorauth"
)

type quotaFetchRequest struct {
	AuthIndex   string            `json:"auth_index"`
	AuthID      string            `json:"auth_id"`
	Provider    string            `json:"provider"`
	StorageJSON []byte            `json:"storage_json"`
	Attributes  map[string]string `json:"attributes"`
}

func quotaIdentifier() map[string]string {
	return map[string]string{"identifier": "cursor"}
}

func quotaDescribe() map[string]any {
	return map[string]any{
		"supported_providers": []string{"cursor", "codex"},
		"display_name":        "Cursor",
		"supports_reset":      false,
	}
}

func (handler *Handler) fetchQuota(ctx context.Context, raw []byte) (any, error) {
	var request quotaFetchRequest
	if err := json.Unmarshal(raw, &request); err != nil {
		return nil, fmt.Errorf("decode quota fetch request: %w", err)
	}
	credentials, err := cursorauth.ParseCredentials(request.StorageJSON)
	if err != nil {
		return nil, err
	}
	if handler.usageAPI == nil {
		return nil, fmt.Errorf("Cursor dashboard usage is unavailable")
	}
	snapshot, err := handler.usageAPI.Fetch(ctx, credentials.AccessToken, credentials.DashboardAccountID())
	if err != nil {
		return nil, err
	}
	reset := ""
	if !snapshot.BillingCycleEnd.IsZero() {
		reset = snapshot.BillingCycleEnd.UTC().Format(time.RFC3339)
	}
	remaining := func(used float64) float64 {
		if used < 0 {
			return 1
		}
		if used > 100 {
			return 0
		}
		return (100 - used) / 100
	}
	plan := strings.TrimSpace(snapshot.MembershipType)
	return map[string]any{
		"subscription": map[string]any{"plan": plan, "tierName": plan},
		"groups": []map[string]any{{
			"displayName": "Cursor",
			"buckets": []map[string]any{
				{
					"window":            "primary",
					"remainingFraction": remaining(snapshot.AutoPercentUsed),
					"resetTime":         reset,
					"description":       "Cursor Models",
				},
				{
					"window":            "weekly",
					"remainingFraction": remaining(snapshot.APIPercentUsed),
					"resetTime":         reset,
					"description":       "Other Models",
				},
			},
		}},
	}, nil
}

func quotaReset() map[string]any {
	return map[string]any{"success": false, "message": "Cursor dashboard usage cannot be reset by the plugin"}
}
