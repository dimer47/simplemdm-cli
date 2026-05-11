package response

import (
	"testing"
)

func TestFlattenItem_WithAttributes(t *testing.T) {
	item := map[string]interface{}{
		"type": "device",
		"id":   float64(121),
		"attributes": map[string]interface{}{
			"name":          "Test Device",
			"serial_number": "ABC123",
		},
	}

	result := FlattenItem(item)

	if result["id"] != float64(121) {
		t.Errorf("expected id 121, got %v", result["id"])
	}
	if result["type"] != "device" {
		t.Errorf("expected type device, got %v", result["type"])
	}
	if result["name"] != "Test Device" {
		t.Errorf("expected name Test Device, got %v", result["name"])
	}
	if result["serial_number"] != "ABC123" {
		t.Errorf("expected serial_number ABC123, got %v", result["serial_number"])
	}
}

func TestFlattenItem_WithoutAttributes(t *testing.T) {
	item := map[string]interface{}{
		"id":     float64(1),
		"status": "ok",
	}

	result := FlattenItem(item)

	if result["id"] != float64(1) {
		t.Errorf("expected id 1, got %v", result["id"])
	}
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %v", result["status"])
	}
}

func TestFlattenItem_Empty(t *testing.T) {
	item := map[string]interface{}{}
	result := FlattenItem(item)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
