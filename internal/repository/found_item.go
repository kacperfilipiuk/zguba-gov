package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kacperfilipiuk/zguba-gov/internal/model"
)

type FoundItemRepo struct {
	db *sql.DB
}

func NewFoundItemRepo(db *sql.DB) *FoundItemRepo {
	return &FoundItemRepo{db: db}
}

func (r *FoundItemRepo) List(p model.ListParams) ([]model.FoundItem, error) {
	query := "SELECT id, municipality_name, municipality_type, municipality_email, item_name, item_category, item_date, item_location, item_status, item_description, pickup_deadline, pickup_location, pickup_hours, pickup_contact, categories, created_at, updated_at FROM found_items WHERE 1=1"
	var args []any

	if p.Category != "" {
		query += " AND item_category = ?"
		args = append(args, p.Category)
	}
	if p.Municipality != "" {
		query += " AND LOWER(municipality_name) LIKE LOWER(?)"
		args = append(args, "%"+p.Municipality+"%")
	}
	if p.Status != "" {
		query += " AND item_status = ?"
		args = append(args, p.Status)
	}
	if p.Search != "" {
		query += " AND (LOWER(item_name) LIKE LOWER(?) OR LOWER(item_description) LIKE LOWER(?) OR LOWER(item_location) LIKE LOWER(?))"
		s := "%" + p.Search + "%"
		args = append(args, s, s, s)
	}

	query += " ORDER BY created_at DESC"

	if p.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", p.Limit)
	}
	if p.Skip > 0 {
		query += fmt.Sprintf(" OFFSET %d", p.Skip)
	}

	return r.queryItems(query, args...)
}

func (r *FoundItemRepo) GetByID(id string) (*model.FoundItem, error) {
	items, err := r.queryItems(
		"SELECT id, municipality_name, municipality_type, municipality_email, item_name, item_category, item_date, item_location, item_status, item_description, pickup_deadline, pickup_location, pickup_hours, pickup_contact, categories, created_at, updated_at FROM found_items WHERE id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &items[0], nil
}

func (r *FoundItemRepo) Create(c model.FoundItemCreate) (*model.FoundItem, error) {
	id := uuid.New().String()
	now := time.Now().UTC()

	catsJSON := "[]"
	if len(c.Categories) > 0 {
		b, _ := json.Marshal(c.Categories)
		catsJSON = string(b)
	}

	status := c.Item.Status
	if status == "" {
		status = "available"
	}

	_, err := r.db.Exec(`
		INSERT INTO found_items (id, municipality_name, municipality_type, municipality_email, item_name, item_category, item_date, item_location, item_status, item_description, pickup_deadline, pickup_location, pickup_hours, pickup_contact, categories, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		c.Municipality.Name, c.Municipality.Type, c.Municipality.ContactEmail,
		c.Item.Name, c.Item.Category, c.Item.Date, c.Item.Location, status, nullStr(c.Item.Description),
		c.Pickup.Deadline, c.Pickup.Location, nullStr(c.Pickup.Hours), nullStr(c.Pickup.Contact),
		catsJSON,
		now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}

	return r.GetByID(id)
}

func (r *FoundItemRepo) Update(id string, u model.FoundItemUpdate) (*model.FoundItem, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}

	var sets []string
	var args []any

	if u.Municipality != nil {
		sets = append(sets, "municipality_name = ?", "municipality_type = ?", "municipality_email = ?")
		args = append(args, u.Municipality.Name, u.Municipality.Type, u.Municipality.ContactEmail)
	}
	if u.Item != nil {
		sets = append(sets, "item_name = ?", "item_category = ?", "item_date = ?", "item_location = ?", "item_status = ?", "item_description = ?")
		args = append(args, u.Item.Name, u.Item.Category, u.Item.Date, u.Item.Location, u.Item.Status, nullStr(u.Item.Description))
	}
	if u.Pickup != nil {
		sets = append(sets, "pickup_deadline = ?", "pickup_location = ?", "pickup_hours = ?", "pickup_contact = ?")
		args = append(args, u.Pickup.Deadline, u.Pickup.Location, nullStr(u.Pickup.Hours), nullStr(u.Pickup.Contact))
	}
	if u.Categories != nil {
		b, _ := json.Marshal(*u.Categories)
		sets = append(sets, "categories = ?")
		args = append(args, string(b))
	}

	if len(sets) == 0 {
		return existing, nil
	}

	sets = append(sets, "updated_at = ?")
	args = append(args, time.Now().UTC())
	args = append(args, id)

	query := fmt.Sprintf("UPDATE found_items SET %s WHERE id = ?", strings.Join(sets, ", "))
	if _, err := r.db.Exec(query, args...); err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return r.GetByID(id)
}

func (r *FoundItemRepo) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM found_items WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *FoundItemRepo) Categories() ([]map[string]string, error) {
	rows, err := r.db.Query("SELECT DISTINCT item_category FROM found_items ORDER BY item_category")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []map[string]string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		label := cat
		if len(label) > 0 {
			label = strings.ToUpper(label[:1]) + label[1:]
		}
		result = append(result, map[string]string{"value": cat, "label": label})
	}
	return result, rows.Err()
}

func (r *FoundItemRepo) Stats() (*model.StatsResponse, error) {
	stats := &model.StatsResponse{}

	if err := r.db.QueryRow("SELECT COUNT(*) FROM found_items").Scan(&stats.FoundItems.Total); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow("SELECT COUNT(*) FROM found_items WHERE item_status = 'available'").Scan(&stats.FoundItems.Available); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow("SELECT COUNT(*) FROM found_items WHERE item_status = 'claimed'").Scan(&stats.FoundItems.Claimed); err != nil {
		return nil, err
	}

	catRows, err := r.db.Query("SELECT item_category, COUNT(*) as cnt FROM found_items GROUP BY item_category ORDER BY cnt DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer func() { _ = catRows.Close() }()
	for catRows.Next() {
		var c model.CategoryCount
		if err := catRows.Scan(&c.Category, &c.Count); err != nil {
			return nil, err
		}
		stats.TopCategories = append(stats.TopCategories, c)
	}

	munRows, err := r.db.Query("SELECT municipality_name, COUNT(*) as cnt FROM found_items GROUP BY municipality_name ORDER BY cnt DESC LIMIT 10")
	if err != nil {
		return nil, err
	}
	defer func() { _ = munRows.Close() }()
	for munRows.Next() {
		var m model.MunicipalityCount
		if err := munRows.Scan(&m.Name, &m.Count); err != nil {
			return nil, err
		}
		stats.TopMunicipalities = append(stats.TopMunicipalities, m)
	}

	if stats.TopCategories == nil {
		stats.TopCategories = []model.CategoryCount{}
	}
	if stats.TopMunicipalities == nil {
		stats.TopMunicipalities = []model.MunicipalityCount{}
	}

	return stats, nil
}

func (r *FoundItemRepo) Count(where string, args ...any) (int, error) {
	query := "SELECT COUNT(*) FROM found_items"
	if where != "" {
		query += " WHERE " + where
	}
	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *FoundItemRepo) QueryRaw(query string, args ...any) ([]model.FoundItem, error) {
	return r.queryItems(query, args...)
}

func (r *FoundItemRepo) queryItems(query string, args ...any) ([]model.FoundItem, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var items []model.FoundItem
	for rows.Next() {
		var fi model.FoundItem
		if err := rows.Scan(
			&fi.ID, &fi.MunicipalityName, &fi.MunicipalityType, &fi.MunicipalityEmail,
			&fi.ItemName, &fi.ItemCategory, &fi.ItemDate, &fi.ItemLocation,
			&fi.ItemStatus, &fi.ItemDescription,
			&fi.PickupDeadline, &fi.PickupLocation, &fi.PickupHours, &fi.PickupContact,
			&fi.Categories,
			&fi.CreatedAt, &fi.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, fi)
	}
	return items, rows.Err()
}

func nullStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
