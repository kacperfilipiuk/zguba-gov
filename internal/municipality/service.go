package municipality

import (
	"database/sql"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type TerritorialUnit struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Email       string `json:"email,omitempty"`
	OfficeName  string `json:"officeName,omitempty"`
	Voivodeship string `json:"voivodeship,omitempty"`
	County      string `json:"county,omitempty"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Search(query, unitType string) []TerritorialUnit {
	if len(query) < 2 {
		return nil
	}

	q := normalize(query)

	var rows *sql.Rows
	var err error

	if unitType != "" {
		rows, err = s.db.Query(
			"SELECT id, name, type, email, office_name, voivodeship, county FROM territorial_units WHERE type = ? AND name_normalized LIKE ? ORDER BY LOCATE(?, name_normalized), name LIMIT 20",
			unitType, "%"+q+"%", q,
		)
	} else {
		rows, err = s.db.Query(
			"SELECT id, name, type, email, office_name, voivodeship, county FROM territorial_units WHERE name_normalized LIKE ? ORDER BY LOCATE(?, name_normalized), name LIMIT 20",
			"%"+q+"%", q,
		)
	}
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	var results []TerritorialUnit
	for rows.Next() {
		var u TerritorialUnit
		if err := rows.Scan(&u.ID, &u.Name, &u.Type, &u.Email, &u.OfficeName, &u.Voivodeship, &u.County); err != nil {
			return nil
		}
		results = append(results, u)
	}

	if results == nil {
		return []TerritorialUnit{}
	}
	return results
}

func (s *Service) GenerateEmail(unit TerritorialUnit) string {
	if unit.Email != "" {
		return unit.Email
	}

	name := strings.ToLower(normalize(unit.Name))
	for _, prefix := range []string{"powiat ", "gmina ", "miasto "} {
		name = strings.TrimPrefix(name, prefix)
	}
	name = strings.ReplaceAll(name, " ", "")

	switch unit.Type {
	case "wojewodztwo":
		return "kontakt@" + name + ".uw.gov.pl"
	case "powiat":
		return "starostwo@" + name + ".pl"
	case "miasto":
		return "urzad@um." + name + ".pl"
	default:
		return "ug@" + name + ".pl"
	}
}

var polishReplacements = map[rune]rune{
	'ą': 'a', 'ć': 'c', 'ę': 'e', 'ł': 'l', 'ń': 'n',
	'ó': 'o', 'ś': 's', 'ź': 'z', 'ż': 'z',
	'Ą': 'a', 'Ć': 'c', 'Ę': 'e', 'Ł': 'l', 'Ń': 'n',
	'Ó': 'o', 'Ś': 's', 'Ź': 'z', 'Ż': 'z',
}

func normalize(s string) string {
	s = norm.NFC.String(s)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if repl, ok := polishReplacements[r]; ok {
			b.WriteRune(repl)
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
