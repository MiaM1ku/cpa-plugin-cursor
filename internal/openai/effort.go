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

var maxModeSuffixes = []string{"1m", "max"}

func ResolveModel(model, effort string, maxMode bool) ModelSelection {
	id := strings.TrimPrefix(strings.TrimSpace(model), "cursor/")
	id = strings.TrimSpace(id)
	explicit := normalizeEffort(effort)
	suffixEffort := ""
	for {
		next, strippedMax := stripMaxModeSuffix(id)
		if strippedMax {
			id = next
			maxMode = true
			continue
		}
		next, suffix, strippedEffort := stripEffortSuffix(id)
		if strippedEffort {
			id = next
			if suffixEffort == "" {
				suffixEffort = effortFromSuffix(suffix)
			}
			continue
		}
		break
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

func stripMaxModeSuffix(id string) (string, bool) {
	lower := strings.ToLower(id)
	for _, suffix := range maxModeSuffixes {
		trimmed, found := strings.CutSuffix(lower, "-"+suffix)
		if !found || trimmed == "" {
			continue
		}
		return id[:len(trimmed)], true
	}
	return id, false
}

func stripEffortSuffix(id string) (string, string, bool) {
	for _, suffix := range effortSuffixes {
		trimmed, found := strings.CutSuffix(id, "-"+suffix)
		if found && trimmed != "" {
			return trimmed, suffix, true
		}
		trimmed, found = strings.CutSuffix(strings.ToLower(id), "-"+suffix)
		if found && trimmed != "" {
			return id[:len(trimmed)], suffix, true
		}
	}
	return id, "", false
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
