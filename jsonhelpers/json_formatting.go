package jsonhelpers

import (
	"encoding/json"
	"sort"
	"strings"
)

// ToJSON is just a shortcut for calling json.Marshal and taking only the first result.
func ToJSON(value any) []byte {
	ret, _ := json.Marshal(value)
	return ret
}

// ToJSONString calls json.Marshal and returns the result as a string.
func ToJSONString(value any) string { return string(ToJSON(value)) }

// CanonicalizeJSON reformats a JSON value so that object properties are alphabetized,
// making comparisons predictable and making it easier for a human reader to find a property.
func CanonicalizeJSON(originalJSON []byte) []byte {
	var rootValue any
	if err := json.Unmarshal(originalJSON, &rootValue); err != nil {
		return originalJSON
	}
	return []byte(toCanonicalizedString(rootValue))
}

func toCanonicalizedString(value any) string {
	switch v := value.(type) {
	case []any:
		items := make([]string, 0, len(v))
		for _, element := range v {
			items = append(items, toCanonicalizedString(element))
		}
		return "[" + strings.Join(items, ",") + "]"

	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		items := make([]string, 0, len(v))
		for _, key := range keys {
			items = append(items, ToJSONString(key)+":"+toCanonicalizedString(v[key]))
		}
		return "{" + strings.Join(items, ",") + "}"

	default:
		return ToJSONString(v)
	}
}
