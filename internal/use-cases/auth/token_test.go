package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testJWTSecret = "test-secret"

func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWT{
			Secret:          testJWTSecret,
			AccessTokenTTL:  time.Minute,
			RefreshTokenTTL: time.Hour,
		},
	}
}

func signTestToken(
	t *testing.T,
	userID, jti, familyID uuid.UUID,
	version int,
	typ string,
	expiresAt time.Time,
) string {
	t.Helper()

	claims := jwt.MapClaims{
		domain.JWTID:           userID,
		domain.JWTTokenVersion: version,
		domain.JWTType:         typ,
		domain.JWTJTI:          jti.String(),
		domain.JWTFAM:          familyID.String(),
		"exp":                  jwt.NewNumericDate(expiresAt),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}

	return signed
}
