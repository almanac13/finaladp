package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string `json:"userId"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	FavoriteLeague string `json:"favoriteLeague,omitempty"`
	FavoriteTeam   string `json:"favoriteTeam,omitempty"`
	jwt.RegisteredClaims
}

// Sign creates token
func Sign(secret []byte, userID, email, role, favLeague, favTeam string) (string, error) {
	claims := Claims{
		UserID:         userID,
		Email:          email,
		Role:           role,
		FavoriteLeague: favLeague,
		FavoriteTeam:   favTeam,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// Verify parses token and returns claims
func Verify(secret []byte, tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		// ensure HS256
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
