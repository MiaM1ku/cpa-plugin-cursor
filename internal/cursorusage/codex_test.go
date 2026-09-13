package cursorusage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_CodexPayload_maps_cursor_percentages_to_primary_and_weekly_windows(t *testing.T) {
	end := time.Date(2026, 9, 24, 7, 46, 2, 0, time.UTC)
	payload := CodexPayload(Snapshot{
		MembershipType:  "pro",
		AutoPercentUsed: 49.05,
		APIPercentUsed:  100,
		BillingCycleEnd: end,
	}, "auth0|user-1", "a@example.com", time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC))

	require.Equal(t, "pro", payload.PlanType)
	require.Equal(t, "auth0|user-1", payload.AccountID)
	require.True(t, payload.RateLimit.LimitReached)
	require.False(t, payload.RateLimit.Allowed)
	require.InDelta(t, 49.05, payload.RateLimit.Primary.UsedPercent, 1e-9)
	require.EqualValues(t, 18000, payload.RateLimit.Primary.LimitWindowSeconds)
	require.InDelta(t, 100, payload.RateLimit.Secondary.UsedPercent, 1e-9)
	require.EqualValues(t, 604800, payload.RateLimit.Secondary.LimitWindowSeconds)
	require.Equal(t, end.Unix(), payload.RateLimit.Primary.ResetAt)
}

func Test_IsCodexUsageURL(t *testing.T) {
	require.True(t, IsCodexUsageURL("https://chatgpt.com/backend-api/wham/usage"))
	require.True(t, IsCodexResetCreditsURL("https://chatgpt.com/backend-api/wham/rate-limit-reset-credits"))
	require.False(t, IsCodexUsageURL("https://cursor.com/api/usage-summary"))
}
