package main

import (
	"reflect"
	"testing"
)

func TestNormalizeNullableUnionInlinesPrimitive(t *testing.T) {
	node := map[string]any{
		"oneOf": []any{
			map[string]any{"type": "string", "format": "date-time"},
			map[string]any{"type": "null"},
		},
	}
	normalizeNullableUnion(node, "oneOf")
	want := map[string]any{"type": "string", "format": "date-time", "nullable": true}
	if !reflect.DeepEqual(node, want) {
		t.Fatalf("got %v, want %v", node, want)
	}
}

func TestNormalizeNullableUnionWrapsRefInAllOf(t *testing.T) {
	// `nullable` next to `$ref` is ignored by oapi-codegen, which then emits a
	// non-pointer field; wrapping in allOf keeps the field nullable.
	ref := map[string]any{"$ref": "#/components/schemas/StateEnum"}
	node := map[string]any{
		"readOnly": true,
		"oneOf":    []any{ref, map[string]any{"type": "null"}},
	}
	normalizeNullableUnion(node, "oneOf")
	want := map[string]any{"readOnly": true, "allOf": []any{ref}, "nullable": true}
	if !reflect.DeepEqual(node, want) {
		t.Fatalf("got %v, want %v", node, want)
	}
}

func TestNormalizeNullableUnionLeavesNonNullableUnions(t *testing.T) {
	node := map[string]any{
		"oneOf": []any{
			map[string]any{"type": "string"},
			map[string]any{"type": "integer"},
		},
	}
	before := map[string]any{"oneOf": node["oneOf"]}
	normalizeNullableUnion(node, "oneOf")
	if !reflect.DeepEqual(node, before) {
		t.Fatalf("expected no change, got %v", node)
	}
}
