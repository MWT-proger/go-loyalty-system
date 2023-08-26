package auth

import (
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/MWT-proger/go-loyalty-system/configs"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// HashPassword(password string) (string, error) - принимает пароль в виде строки,
// возвращает хеш bcrypt пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash(password, hash string) bool - сравнивает хешированный пароль bcrypt
// с его возможным эквивалентом в виде открытого текста.
// Возвращает true в случае успеха или false в случае неудачи.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword(password string) bool - принимает пароль в виде строки,
// и проверяет, соответствует ли строка следующим требованиям к паролю:
// - Длина пароля должна быть не менее 8 символов.
// - Пароль должен содержать хотя бы одну заглавную букву.
// - Пароль должен содержать хотя бы одну строчную букву.
// - Пароль должен содержать хотя бы одну цифру.
// возвращает true или false
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	if !hasUppercase(password) {
		return false
	}

	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString
	if !hasLowercase(password) {
		return false
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString

	return hasNumber(password)

}

// BuildJWTString(UserID uuid.UUID) (string, error) принимает UserID
// создаёт токен для пользователя
// и в случае успеха возвращает его в виде строки
func BuildJWTString(UserID uuid.UUID) (string, error) {
	conf := configs.GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           UserID,
	})

	tokenString, err := token.SignedString([]byte(conf.Auth.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetAuthTokenToCookie(w http.ResponseWriter, token string) {
	newCookie := http.Cookie{Name: "token"}
	newCookie.Value = token
	http.SetCookie(w, &newCookie)

}
