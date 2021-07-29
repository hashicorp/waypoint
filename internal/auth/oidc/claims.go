package oidc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/pointerstructure"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// SelectorData returns the data for go-bexpr for selector evaluation.
// This is useful for server-side OIDC implementations, not client.
func SelectorData(
	am *pb.AuthMethod_OIDC,
	idClaims, userClaims json.RawMessage,
) (map[string]interface{}, error) {
	// Extract the claims into a map[string]interface{}
	var all map[string]interface{}
	if err := json.Unmarshal([]byte(idClaims), &all); err != nil {
		return nil, err
	}
	if len(userClaims) > 0 {
		// Keep these cause we never let these get overwritten
		iss, issOk := all["iss"]
		sub, subOk := all["sub"]

		if err := json.Unmarshal([]byte(userClaims), &all); err != nil {
			return nil, err
		}

		if issOk {
			all["iss"] = iss
		}
		if subOk {
			all["sub"] = sub
		}
	}

	// I expect SelectorData will do more in the future which is why
	// this is just calling this other function directly and not doing
	// anything else today.
	return extractClaims(am, all)
}

// extractClaims takes the claim mapping configuration of the OIDC
// auth method, extracts the claims, and returns a map of data that can
// be used with go-bexpr.
func extractClaims(
	am *pb.AuthMethod_OIDC,
	all map[string]interface{},
) (map[string]interface{}, error) {
	values, err := extractMappings(all, am.ClaimMappings)
	if err != nil {
		return nil, err
	}

	list, err := extractListMappings(all, am.ListClaimMappings)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"value": values,
		"list":  list,
	}, nil
}

// extractMappings extracts the string value mappings.
func extractMappings(
	all map[string]interface{},
	mapping map[string]string,
) (map[string]string, error) {
	result := make(map[string]string)
	for source, target := range mapping {
		rawValue := getClaim(all, source)
		if rawValue == nil {
			continue
		}

		strValue, ok := stringifyClaimValue(rawValue)
		if !ok {
			return nil, fmt.Errorf(
				"error converting claim '%s' to string from unknown type %T",
				source, rawValue)
		}

		result[target] = strValue
	}

	return result, nil
}

// extractListMappings builds a metadata map of string list values from a set
// of claims and claims mappings.  The referenced claims must be strings and
// the claims mappings must be of the structure:
//
//   {
//       "/some/claim/pointer": "metadata_key1",
//       "another_claim": "metadata_key2",
//        ...
//   }
func extractListMappings(
	all map[string]interface{}, mappings map[string]string,
) (map[string][]string, error) {
	result := make(map[string][]string)
	for source, target := range mappings {
		rawValue := getClaim(all, source)
		if rawValue == nil {
			continue
		}

		rawList, ok := normalizeList(rawValue)
		if !ok {
			return nil, fmt.Errorf("%q list claim could not be converted to string list", source)
		}

		list := make([]string, 0, len(rawList))
		for _, raw := range rawList {
			value, ok := stringifyClaimValue(raw)
			if !ok {
				return nil, fmt.Errorf(
					"value %v in %q list claim could not be parsed as string",
					raw, source)
			}

			if value == "" {
				continue
			}

			list = append(list, value)
		}

		result[target] = list
	}

	return result, nil
}

// getClaim returns a claim value from allClaims given a provided claim string.
// If this string is a valid JSONPointer, it will be interpreted as such to
// locate the claim. Otherwise, the claim string will be used directly.
//
// There is no fixup done to the returned data type here. That happens a layer
// up in the caller.
func getClaim(all map[string]interface{}, claim string) interface{} {
	if !strings.HasPrefix(claim, "/") {
		return all[claim]
	}

	val, err := pointerstructure.Get(all, claim)
	if err != nil {
		// We silently drop the error since keys that are invalid
		// just have no values.
		return nil
	}

	return val
}

// stringifyClaimValue will try to convert the provided raw value into a
// faithful string representation of that value per these rules:
//
// - strings      => unchanged
// - bool         => "true" / "false"
// - json.Number  => String()
// - float32/64   => truncated to int64 and then formatted as an ascii string
// - intXX/uintXX => casted to int64 and then formatted as an ascii string
//
// If successful the string value and true are returned. otherwise an empty
// string and false are returned.
func stringifyClaimValue(rawValue interface{}) (string, bool) {
	switch v := rawValue.(type) {
	case string:
		return v, true
	case bool:
		return strconv.FormatBool(v), true
	case json.Number:
		return v.String(), true
	case float64:
		// The claims unmarshalled by go-oidc don't use UseNumber, so
		// they'll come in as float64 instead of an integer or json.Number.
		return strconv.FormatInt(int64(v), 10), true

		// The numerical type cases following here are only here for the sake
		// of numerical type completion. Everything is truncated to an integer
		// before being stringified.
	case float32:
		return strconv.FormatInt(int64(v), 10), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case int:
		return strconv.FormatInt(int64(v), 10), true
	case uint8:
		return strconv.FormatInt(int64(v), 10), true
	case uint16:
		return strconv.FormatInt(int64(v), 10), true
	case uint32:
		return strconv.FormatInt(int64(v), 10), true
	case uint64:
		return strconv.FormatInt(int64(v), 10), true
	case uint:
		return strconv.FormatInt(int64(v), 10), true
	default:
		return "", false
	}
}

// normalizeList takes an item or a slice and returns a slice. This is useful
// when providers are expected to return a list (typically of strings) but
// reduce it to a non-slice type when the list count is 1.
//
// There is no fixup done to elements of the returned slice here. That happens
// a layer up in the caller.
func normalizeList(raw interface{}) ([]interface{}, bool) {
	switch v := raw.(type) {
	case []interface{}:
		return v, true
	case string, // note: this list should be the same as stringifyClaimValue
		bool,
		json.Number,
		float64,
		float32,
		int8,
		int16,
		int32,
		int64,
		int,
		uint8,
		uint16,
		uint32,
		uint64,
		uint:
		return []interface{}{v}, true
	default:
		return nil, false
	}

}
