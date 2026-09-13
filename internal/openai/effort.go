package openai

import "strings"

type ModelSelection struct {
	ID       string
	Effort   string
	MaxMode  bool
	Thinking bool
	Fast     bool
}

var effortSuffixes = []string{
	"thinking-extra-high", "thinking-xhigh", "thinking-high", "thinking-medium",
	"thinking-low", "thinking-minimal", "thinking-none", "thinking-max", "thinking",
	"nothinking", "extra-high", "xhigh", "high", "medium", "low", "minimal", "none", "max",
}

func ResolveModel(model, effort string, maxMode bool) ModelSelection {
	id := strings.TrimPrefix(strings.TrimSpace(model), "cursor/")
	id = strings.TrimSpace(id)
	explicit := normalizeEffort(effort)
	suffixEffort := ""
	thinking := false
	thinkingSet := false
	fast := false
	for {
		next, strippedFast := stripTaggedSuffix(id, "fast")
		if strippedFast {
			id = next
			fast = true
			continue
		}
		next, strippedMax := stripTaggedSuffix(id, "1m")
		if strippedMax {
			id = next
			maxMode = true
			continue
		}
		next, suffix, strippedEffort := stripEffortSuffix(id)
		if strippedEffort {
			id = next
			if suffix == "thinking" || strings.HasPrefix(suffix, "thinking-") {
				thinking = true
				thinkingSet = true
			}
			if suffix == "nothinking" {
				thinking = false
				thinkingSet = true
			}
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
	return ModelSelection{ID: id, Effort: resolved, MaxMode: maxMode, Thinking: thinking && thinkingSet, Fast: fast}
}

func CollapseModelID(model string) string {
	return ResolveModel(model, "", false).ID
}

func stripTaggedSuffix(id, suffix string) (string, bool) {
	lower := strings.ToLower(id)
	trimmed, found := strings.CutSuffix(lower, "-"+suffix)
	if !found || trimmed == "" {
		return id, false
	}
	return id[:len(trimmed)], true
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
	case "thinking-xhigh", "thinking-extra-high", "xhigh", "extra-high":
		return "xhigh"
	case "thinking-high", "high":
		return "high"
	case "thinking-medium", "medium":
		return "medium"
	case "thinking-low", "low":
		return "low"
	case "thinking-max", "max":
		return "max"
	case "thinking-minimal", "minimal":
		return "minimal"
	case "thinking-none", "none":
		return "none"
	default:
		return ""
	}
}

func normalizeEffort(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "none", "minimal", "low", "medium", "high", "xhigh", "max":
		return strings.ToLower(strings.TrimSpace(value))
	case "extra-high", "extra_high":
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
	return nil
}
