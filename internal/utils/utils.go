package utils

import "regexp"

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
