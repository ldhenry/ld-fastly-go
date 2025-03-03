package main

import (
	"encoding/json"
	"sort"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	"github.com/launchdarkly/go-sdk-common/v3/ldreason"
)

type Response struct {
	Context        ldcontext.Context         `json:"context"`
	Animal         string                    `json:"animal"`
	Reason         ldreason.EvaluationDetail `json:"reason"`
	ServiceVersion string                    `json:"service_version"`
}

// MarshalJSON implements custom JSON marshaling with sorted keys
func (r Response) MarshalJSON() ([]byte, error) {
	// Create a map to store all fields
	m := make(map[string]interface{})

	// Convert context to map for sorting
	contextJSON, err := json.Marshal(r.Context)
	if err != nil {
		return nil, err
	}
	var contextMap map[string]interface{}
	if err := json.Unmarshal(contextJSON, &contextMap); err != nil {
		return nil, err
	}

	// Convert reason to map for sorting
	reasonJSON, err := json.Marshal(r.Reason)
	if err != nil {
		return nil, err
	}
	var reasonMap map[string]interface{}
	if err := json.Unmarshal(reasonJSON, &reasonMap); err != nil {
		return nil, err
	}

	m["animal"] = r.Animal
	m["context"] = sortMapRecursive(contextMap)
	m["reason"] = sortMapRecursive(reasonMap)
	m["service_version"] = r.ServiceVersion

	return json.Marshal(m)
}

// sortMapRecursive sorts map keys recursively
func sortMapRecursive(m map[string]interface{}) map[string]interface{} {
	sorted := make(map[string]interface{})
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		switch val := v.(type) {
		case map[string]interface{}:
			sorted[k] = sortMapRecursive(val)
		case []interface{}:
			sorted[k] = sortSliceRecursive(val)
		default:
			sorted[k] = v
		}
	}
	return sorted
}

// sortSliceRecursive sorts elements in slices that are maps
func sortSliceRecursive(s []interface{}) []interface{} {
	for i, v := range s {
		if m, ok := v.(map[string]interface{}); ok {
			s[i] = sortMapRecursive(m)
		}
	}
	return s
}
