package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword recibe la contraseña en claro y devuelve su hash bcrypt
func HashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword compara el hash bcrypt con la contraseña en claro.
// Devuelve nil si coinciden o un error en caso contrario.
func CheckPassword(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
