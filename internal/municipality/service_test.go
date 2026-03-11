package municipality

import (
	"testing"
)

func TestNormalize_PolishCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Łódź", "lodz"},
		{"Kraków", "krakow"},
		{"Gdańsk", "gdansk"},
		{"Wrocław", "wroclaw"},
		{"Białystok", "bialystok"},
		{"Szczecin", "szczecin"},
		{"Gorzów Wielkopolski", "gorzow wielkopolski"},
		{"Żółć", "zolc"},
		{"ŁÓDŹ", "lodz"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalize(tt.input)
			if result != tt.expected {
				t.Errorf("normalize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateEmail_WithExistingEmail(t *testing.T) {
	svc := &Service{}
	unit := TerritorialUnit{
		Name:  "Warszawa",
		Type:  "miasto",
		Email: "kontakt@warszawa.pl",
	}
	result := svc.GenerateEmail(unit)
	if result != "kontakt@warszawa.pl" {
		t.Errorf("expected existing email, got %q", result)
	}
}

func TestGenerateEmail_Wojewodztwo(t *testing.T) {
	svc := &Service{}
	unit := TerritorialUnit{
		Name: "Małopolskie",
		Type: "wojewodztwo",
	}
	result := svc.GenerateEmail(unit)
	if result != "kontakt@malopolskie.uw.gov.pl" {
		t.Errorf("expected kontakt@malopolskie.uw.gov.pl, got %q", result)
	}
}

func TestGenerateEmail_Powiat(t *testing.T) {
	svc := &Service{}
	unit := TerritorialUnit{
		Name: "Powiat Krakowski",
		Type: "powiat",
	}
	result := svc.GenerateEmail(unit)
	if result != "starostwo@krakowski.pl" {
		t.Errorf("expected starostwo@krakowski.pl, got %q", result)
	}
}

func TestGenerateEmail_Miasto(t *testing.T) {
	svc := &Service{}
	unit := TerritorialUnit{
		Name: "Miasto Kraków",
		Type: "miasto",
	}
	result := svc.GenerateEmail(unit)
	if result != "urzad@um.krakow.pl" {
		t.Errorf("expected urzad@um.krakow.pl, got %q", result)
	}
}

func TestGenerateEmail_Gmina(t *testing.T) {
	svc := &Service{}
	unit := TerritorialUnit{
		Name: "Gmina Wieliczka",
		Type: "gmina",
	}
	result := svc.GenerateEmail(unit)
	if result != "ug@wieliczka.pl" {
		t.Errorf("expected ug@wieliczka.pl, got %q", result)
	}
}

func TestSearch_QueryTooShort(t *testing.T) {
	svc := &Service{}
	result := svc.Search("a", "")
	if result != nil {
		t.Errorf("expected nil for short query, got %v", result)
	}
}
