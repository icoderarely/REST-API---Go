package models

import "database/sql"

type Exec struct {
	ID                   int            `json:"id,omitempty"`
	FirstName            string         `json:"first_name,omitempty"`
	LastName             string         `json:"last_name,omitempty"`
	Email                string         `json:"email,omitempty"`
	Username             string         `json:"username,omitempty"`
	Password             string         `json:"password,omitempty"`
	PasswordChangedAt    sql.NullString `json:"password_changed_at"`
	UserCreatedAt        sql.NullString `json:"user_created_at"`
	PasswordResetToken   sql.NullString `json:"password_reset_token"`
	PasswordTokenExpires sql.NullString `json:"password_token_expires"`
	InactiveStatus       sql.NullBool   `json:"inactive_status"`
	Role                 string         `json:"role,omitempty"`
}
