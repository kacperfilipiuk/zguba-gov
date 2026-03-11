package municipality

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

//go:embed territorial-units.json
var dataFS embed.FS

type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*f = FlexString(fmt.Sprintf("%07s", n.String()))
	return nil
}

type TerritorialUnit struct {
	ID          FlexString `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Email       string     `json:"email,omitempty"`
	OfficeName  string     `json:"officeName,omitempty"`
	Voivodeship string     `json:"voivodeship,omitempty"`
	County      string     `json:"county,omitempty"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) (*Service, error) {
	svc := &Service{db: db}
	if err := svc.seed(); err != nil {
		return nil, fmt.Errorf("seed territorial units: %w", err)
	}
	return svc, nil
}

func (s *Service) seed() error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM territorial_units").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data, err := dataFS.ReadFile("territorial-units.json")
	if err != nil {
		return err
	}

	var units []TerritorialUnit
	if err := json.Unmarshal(data, &units); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.Prepare("INSERT INTO territorial_units (id, name, type, email, office_name, voivodeship, county, name_normalized) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, u := range units {
		_, err := stmt.Exec(string(u.ID), u.Name, u.Type, u.Email, u.OfficeName, u.Voivodeship, u.County, normalize(u.Name))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
		var id string
		if err := rows.Scan(&id, &u.Name, &u.Type, &u.Email, &u.OfficeName, &u.Voivodeship, &u.County); err != nil {
			return nil
		}
		u.ID = FlexString(id)
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
