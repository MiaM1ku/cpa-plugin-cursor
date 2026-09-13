package cursorusage

import (
	"strings"
	"time"
)

const (
	codexFiveHourSeconds = 18000
	codexWeekSeconds     = 604800
)

type CodexUsagePayload struct {
	PlanType  string         `json:"plan_type,omitempty"`
	AccountID string         `json:"account_id,omitempty"`
	Email     string         `json:"email,omitempty"`
	RateLimit CodexRateLimit `json:"rate_limit"`
	ResetAt   string         `json:"subscription_active_until,omitempty"`
}

type CodexRateLimit struct {
	Allowed      bool             `json:"allowed"`
	LimitReached bool             `json:"limit_reached"`
	Primary      CodexUsageWindow `json:"primary_window"`
	Secondary    CodexUsageWindow `json:"secondary_window"`
}

type CodexUsageWindow struct {
	UsedPercent        float64 `json:"used_percent"`
	LimitWindowSeconds int64   `json:"limit_window_seconds"`
	ResetAfterSeconds  int64   `json:"reset_after_seconds,omitempty"`
	ResetAt            int64   `json:"reset_at,omitempty"`
}

func CodexPayload(snapshot Snapshot, accountID, email string, now time.Time) CodexUsagePayload {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	resetAt := snapshot.BillingCycleEnd
	resetUnix := int64(0)
	resetAfter := int64(0)
	if !resetAt.IsZero() {
		resetUnix = resetAt.Unix()
		if resetAt.After(now) {
			resetAfter = int64(resetAt.Sub(now).Seconds())
		}
	}
	primaryUsed := snapshot.AutoPercentUsed
	if snapshot.CursorModels.UsagePercent != 0 {
		primaryUsed = snapshot.CursorModels.UsagePercent
	}
	secondaryUsed := snapshot.APIPercentUsed
	if snapshot.OtherModels.UsagePercent != 0 {
		secondaryUsed = snapshot.OtherModels.UsagePercent
	}
	limitReached := primaryUsed >= 100 || secondaryUsed >= 100
	plan := strings.TrimSpace(snapshot.MembershipType)
	if plan == "" {
		plan = "pro"
	}
	return CodexUsagePayload{
		PlanType:  plan,
		AccountID: strings.TrimSpace(accountID),
		Email:     strings.TrimSpace(email),
		ResetAt:   resetAt.Format(time.RFC3339),
		RateLimit: CodexRateLimit{
			Allowed:      !limitReached,
			LimitReached: limitReached,
			Primary: CodexUsageWindow{
				UsedPercent:        primaryUsed,
				LimitWindowSeconds: codexFiveHourSeconds,
				ResetAfterSeconds:  resetAfter,
				ResetAt:            resetUnix,
			},
			Secondary: CodexUsageWindow{
				UsedPercent:        secondaryUsed,
				LimitWindowSeconds: codexWeekSeconds,
				ResetAfterSeconds:  resetAfter,
				ResetAt:            resetUnix,
			},
		},
	}
}

func IsCodexUsageURL(rawURL string) bool {
	value := strings.ToLower(rawURL)
	return strings.Contains(value, "chatgpt.com/backend-api/wham/usage") ||
		strings.Contains(value, "/backend-api/wham/usage")
}

func IsCodexResetCreditsURL(rawURL string) bool {
	value := strings.ToLower(rawURL)
	return strings.Contains(value, "rate-limit-reset-credits")
}
