package openai

import "strings"

const defaultEffort = "medium"

func ResolveWireID(sel ModelSelection, catalog []string) string {
	sel.ID = strings.TrimSpace(sel.ID)
	if sel.ID == "" {
		sel.ID = "auto"
	}
	if sel.ID == "auto" || sel.ID == "default" {
		return sel.ID
	}
	catalogSet := make(map[string]struct{}, len(catalog))
	for _, id := range catalog {
		id = strings.TrimSpace(id)
		if id == "" || strings.HasSuffix(strings.ToLower(id), "-fast") {
			continue
		}
		catalogSet[id] = struct{}{}
	}
	candidates := wireCandidates(sel)
	for _, candidate := range candidates {
		if _, ok := catalogSet[candidate]; ok {
			return candidate
		}
	}
	if len(candidates) == 0 {
		return sel.ID
	}
	return candidates[0]
}

func wireCandidates(sel ModelSelection) []string {
	out := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	add := func(id string) {
		if id == "" {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	for _, effort := range effortSpellings(sel.Effort) {
		if sel.Thinking {
			add(sel.ID + "-thinking-" + effort)
		} else {
			add(sel.ID + "-" + effort)
		}
	}
	if sel.Thinking {
		add(sel.ID + "-thinking")
	}
	add(sel.ID)
	return out
}

func effortSpellings(effort string) []string {
	effort = normalizeEffort(effort)
	if effort == "" {
		effort = defaultEffort
	}
	if effort == "xhigh" {
		return []string{"xhigh", "extra-high"}
	}
	return []string{effort}
}
