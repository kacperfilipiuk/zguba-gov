package model

import (
	"database/sql"
	"testing"
	"time"
)

func TestToResponse_Basic(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	fi := FoundItem{
		ID:                "abc-123",
		MunicipalityName:  "Warszawa",
		MunicipalityType:  "miasto",
		MunicipalityEmail: "urzad@warszawa.pl",
		ItemName:          "Portfel",
		ItemCategory:      "dokumenty",
		ItemDate:          "2025-01-10",
		ItemLocation:      "Park Łazienkowski",
		ItemStatus:        "available",
		ItemDescription:   sql.NullString{String: "Czarny portfel skórzany", Valid: true},
		PickupDeadline:    30,
		PickupLocation:    "Biuro rzeczy znalezionych",
		PickupHours:       sql.NullString{String: "8:00-16:00", Valid: true},
		PickupContact:     sql.NullString{String: "Jan Kowalski", Valid: true},
		Categories:        sql.NullString{String: `["dokumenty","portfele"]`, Valid: true},
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	resp := fi.ToResponse()

	if resp.ID != "abc-123" {
		t.Errorf("expected ID abc-123, got %s", resp.ID)
	}
	if resp.Municipality.Name != "Warszawa" {
		t.Errorf("expected municipality Warszawa, got %s", resp.Municipality.Name)
	}
	if resp.Item.Name != "Portfel" {
		t.Errorf("expected item name Portfel, got %s", resp.Item.Name)
	}
	if resp.Item.Description != "Czarny portfel skórzany" {
		t.Errorf("unexpected description: %s", resp.Item.Description)
	}
	if resp.Pickup.Deadline != 30 {
		t.Errorf("expected deadline 30, got %d", resp.Pickup.Deadline)
	}
	if resp.Pickup.Hours != "8:00-16:00" {
		t.Errorf("unexpected hours: %s", resp.Pickup.Hours)
	}
	if len(resp.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(resp.Categories))
	}
	if resp.Categories[0] != "dokumenty" || resp.Categories[1] != "portfele" {
		t.Errorf("unexpected categories: %v", resp.Categories)
	}
	if resp.CreatedAt != "2025-01-15T10:30:00Z" {
		t.Errorf("unexpected createdAt: %s", resp.CreatedAt)
	}
}

func TestToResponse_NullFields(t *testing.T) {
	fi := FoundItem{
		ID:              "def-456",
		ItemDescription: sql.NullString{},
		PickupHours:     sql.NullString{},
		PickupContact:   sql.NullString{},
		Categories:      sql.NullString{},
	}

	resp := fi.ToResponse()

	if resp.Item.Description != "" {
		t.Errorf("expected empty description, got %q", resp.Item.Description)
	}
	if resp.Pickup.Hours != "" {
		t.Errorf("expected empty hours, got %q", resp.Pickup.Hours)
	}
	if resp.Categories == nil {
		t.Error("categories should not be nil")
	}
	if len(resp.Categories) != 0 {
		t.Errorf("expected empty categories, got %v", resp.Categories)
	}
}

func TestToResponse_InvalidCategoriesJSON(t *testing.T) {
	fi := FoundItem{
		Categories: sql.NullString{String: "not-json", Valid: true},
	}

	resp := fi.ToResponse()

	if len(resp.Categories) != 0 {
		t.Errorf("expected empty categories for invalid JSON, got %v", resp.Categories)
	}
}
