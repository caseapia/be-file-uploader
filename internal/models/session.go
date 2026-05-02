package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID          string      `bun:"id,pk" json:"id"`
	UserID      int         `bun:"user_id,notnull" json:"user_id"`
	IPAddress   string      `bun:"ip_address,notnull" json:"ip_address"`
	UserAgent   string      `bun:"user_agent,notnull" json:"user_agent"`
	GeoString   string      `bun:"geo_string" json:"-"`
	IsActive    bool        `bun:"is_active,notnull,default:true" json:"is_active"`
	ExpiresAt   time.Time   `bun:"expires_at,notnull" json:"expires_at"`
	CreatedAt   time.Time   `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	RefreshHash string      `bun:"refresh_hash,notnull" json:"refresh_hash"`
	Geolocation Geolocation `bun:"-" json:"geolocation"`
}

type SessionClaims struct {
	UserID    int    `json:"sub"`
	SessionID string `json:"sid"`
	TokenVer  int    `json:"tv"`
	jwt.RegisteredClaims
}
