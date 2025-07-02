package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims define la info que incluís en el token
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"` // opcional: admin, user, etc.
	jwt.RegisteredClaims
}

// GenerateToken crea un JWT firmado con HS256.
// - userID: identificador de usuario.
// - role: opcional, puede quedar vacío.
// - secret: clave secreta (e.g. cfg.JWTSecret).
// - expiresIn: duración del token (e.g. time.Hour*24).
func GenerateToken(userID, role, secret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken valida el JWT y devuelve las Claims.
// Retorna error si el token está vencido, mal formado o la firma no coincide.
func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Asegurarse de que sea HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
