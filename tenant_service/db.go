package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Init initializes the DB connection
func Init(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

func generateUUID() string {
	return uuid.New().String()
}

// CreateTenant inserts a new tenant into DB
func CreateTenant(name string) (*Tenant, error) {
	id := generateUUID()
	now := time.Now()

	_, err := db.Exec(`INSERT INTO tenants (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		id, name, now, now)
	if err != nil {
		return nil, err
	}

	return &Tenant{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetTenant retrieves a tenant by ID
func GetTenant(id string) (*Tenant, error) {
	row := db.QueryRow(`SELECT id, name, created_at, updated_at FROM tenants WHERE id = $1`, id)

	var tenant Tenant
	err := row.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// UpdateTenant updates tenant name and returns the updated tenant
func UpdateTenant(id, name string) (*Tenant, error) {
	now := time.Now()

	_, err := db.Exec(`UPDATE tenants SET name = $1, updated_at = $2 WHERE id = $3`, name, now, id)
	if err != nil {
		return nil, err
	}

	return GetTenant(id)
}

// DeleteTenant deletes a tenant by ID
func DeleteTenant(id string) (bool, error) {
	res, err := db.Exec(`DELETE FROM tenants WHERE id = $1`, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
