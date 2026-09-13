package openai

import "strings"

type ModelSelection struct {
	ID      string
	Effort  string
	MaxMode bool
}

var effortSuffixes = []string{
	"thinking-xhigh", "thinking-high", "thinking-medium", "thinking-low",
	"thinking", "xhigh", "extra-high", "high", "medium", "low", "minimal", "none", "fast",
}

func ResolveModel(model, effort string, maxMode bool) ModelSelection {
	id := strings.TrimPrefix(strings.TrimSpace(model), "cursor/")
	id = strings.TrimSpace(id)
	explicit := normalizeEffort(effort)
	suffixEffort := ""
	for _, suffix := range effortSuffixes {
		trimmed, found := strings.CutSuffix(id, "-"+suffix)
		if !found || trimmed == "" {
			continue
		}
		id = trimmed
		suffixEffort = effortFromSuffix(suffix)
		break
	}
	if strings.HasSuffix(strings.ToLower(id), "-1m") {
		id = strings.TrimSuffix(id, "-1m")
		id = strings.TrimSuffix(id, "-1M")
		maxMode = true
	}
	if id == "" {
		id = "auto"
	}
	resolved := explicit
	if resolved == "" {
		resolved = suffixEffort
	}
	return ModelSelection{ID: id, Effort: resolved, MaxMode: maxMode}
}

func CollapseModelID(model string) string {
	return ResolveModel(model, "", false).ID
}

func effortFromSuffix(suffix string) string {
	switch suffix {
	case "thinking-xhigh", "xhigh", "extra-high":
		return "xhigh"
	case "thinking-high", "high":
		return "high"
	case "thinking-medium", "medium":
		return "medium"
	case "thinking-low", "low":
		return "low"
	case "thinking":
		return "high"
	case "minimal", "none":
		return suffix
	default:
		return ""
	}
}

func normalizeEffort(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "none", "minimal", "low", "medium", "high", "xhigh":
		return strings.ToLower(strings.TrimSpace(value))
	case "extra-high", "extra_high", "max":
		return "xhigh"
	case "x-high", "x_high":
		return "xhigh"
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func EffortParameters(effort string) []struct{ ID, Value string } {
	effort = normalizeEffort(effort)
	if effort == "" {
		return nil
	}
	return []struct{ ID, Value string }{{ID: "effort", Value: effort}}
}
