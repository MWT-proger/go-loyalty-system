package auth

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword(password string) (string, error) - принимает пароль в виде строки,
// возвращает хеш bcrypt пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash(password, hash string) bool - сравнивает хешированный пароль bcrypt
// с его возможным эквивалентом в виде открытого текста.
// Возвращает ноль в случае успеха или ошибку в случае неудачи.
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
