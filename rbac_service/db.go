package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// rolePermissions maps roles to allowed resources/actions
var rolePermissions = map[string][]string{
	"admin":  {"create_user", "delete_user", "create_tenant", "view_dashboard"},
	"editor": {"create_user", "view_dashboard"},
	"viewer": {"view_dashboard"},
}

// InitDB initializes the database connection.
func InitDB(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	log.Println("Connected to RBAC DB")
	return nil
}

// InsertUserRole assigns a role to a user.
func InsertUserRole(userID, role string) error {
	_, err := db.Exec(
		"INSERT INTO user_roles (user_id, role) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID, role)
	if err != nil {
		log.Printf("Error inserting role: %v", err)
	}
	return err
}

// RemoveUserRole revokes a role from a user.
func RemoveUserRole(userID, role string) error {
	_, err := db.Exec("DELETE FROM user_roles WHERE user_id=$1 AND role=$2", userID, role)
	if err != nil {
		log.Printf("Error removing role: %v", err)
	}
	return err
}

// CheckUserAccess checks whether the user has a role that permits access to the requested resource.
func CheckUserAccess(userID, resource string) (bool, error) {
	rows, err := db.Query("SELECT role FROM user_roles WHERE user_id=$1", userID)
	if err != nil {
		log.Printf("DB query error in CheckUserAccess: %v", err)
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			log.Printf("Error scanning role: %v", err)
			continue
		}

		allowedResources, exists := rolePermissions[role]
		if !exists {
			continue
		}

		for _, allowed := range allowedResources {
			if allowed == resource {
				return true, nil
			}
		}
	}

	return false, nil
}
