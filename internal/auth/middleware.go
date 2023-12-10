package auth

import (
	"net/http"

	"github.com/gofrs/uuid"
)

// AuthCookieMiddleware(next http.Handler) http.Handler — middleware-для входящих HTTP-запросов.
// Проверяет у пользователя подписанную куку, содержащую уникальный идентификатор пользователя,
// и записывает ID пользователя в контекст
// если такой куки не существует или она не проходит проверку подлинности возвращает статус Unauthorized.
func AuthCookieMiddleware(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			UserID      uuid.UUID
			tokenString string
			ctx         = r.Context()
			token, err  = r.Cookie(NameCookie)
		)

		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		tokenString = token.Value
		UserID = GetUserID(tokenString)

		if UserID == uuid.Nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx = WithUserID(ctx, UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
