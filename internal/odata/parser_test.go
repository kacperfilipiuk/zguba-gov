package odata

import (
	"testing"
)

func TestParseFilter_Empty(t *testing.T) {
	fc, err := ParseFilter("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Where != "" {
		t.Errorf("expected empty where, got %q", fc.Where)
	}
	if len(fc.Args) != 0 {
		t.Errorf("expected no args, got %v", fc.Args)
	}
}

func TestParseFilter_Eq(t *testing.T) {
	fc, err := ParseFilter("item_status eq 'available'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Where != "item_status = ?" {
		t.Errorf("expected 'item_status = ?', got %q", fc.Where)
	}
	if len(fc.Args) != 1 || fc.Args[0] != "available" {
		t.Errorf("expected args [available], got %v", fc.Args)
	}
}

func TestParseFilter_Contains(t *testing.T) {
	fc, err := ParseFilter("contains(item_name, 'portfel')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Where != "item_name LIKE ?" {
		t.Errorf("expected 'item_name LIKE ?', got %q", fc.Where)
	}
	if len(fc.Args) != 1 || fc.Args[0] != "%portfel%" {
		t.Errorf("expected args [%%portfel%%], got %v", fc.Args)
	}
}

func TestParseFilter_StartsWith(t *testing.T) {
	fc, err := ParseFilter("startswith(municipality_name, 'Warsz')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Where != "municipality_name LIKE ?" {
		t.Errorf("expected 'municipality_name LIKE ?', got %q", fc.Where)
	}
	if len(fc.Args) != 1 || fc.Args[0] != "Warsz%" {
		t.Errorf("expected args [Warsz%%], got %v", fc.Args)
	}
}

func TestParseFilter_MultipleConditions(t *testing.T) {
	fc, err := ParseFilter("item_status eq 'available' and contains(item_name, 'klucz')")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Where != "item_status = ? AND item_name LIKE ?" {
		t.Errorf("unexpected where: %q", fc.Where)
	}
	if len(fc.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(fc.Args))
	}
}

func TestParseFilter_UnsupportedField(t *testing.T) {
	_, err := ParseFilter("unknown_field eq 'value'")
	if err == nil {
		t.Fatal("expected error for unsupported field")
	}
}

func TestParseOrderBy_Empty(t *testing.T) {
	result := ParseOrderBy("")
	if result != "created_at DESC" {
		t.Errorf("expected default order, got %q", result)
	}
}

func TestParseOrderBy_ValidAsc(t *testing.T) {
	result := ParseOrderBy("item_name asc")
	if result != "item_name ASC" {
		t.Errorf("expected 'item_name ASC', got %q", result)
	}
}

func TestParseOrderBy_ValidDesc(t *testing.T) {
	result := ParseOrderBy("created_at desc")
	if result != "created_at DESC" {
		t.Errorf("expected 'created_at DESC', got %q", result)
	}
}

func TestParseOrderBy_DefaultDirection(t *testing.T) {
	result := ParseOrderBy("item_date")
	if result != "item_date ASC" {
		t.Errorf("expected 'item_date ASC', got %q", result)
	}
}

func TestParseOrderBy_UnsupportedField(t *testing.T) {
	result := ParseOrderBy("unknown_field desc")
	if result != "created_at DESC" {
		t.Errorf("expected default fallback, got %q", result)
	}
}
