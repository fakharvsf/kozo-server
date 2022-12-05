package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

func ParseToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeadersArr := r.Header["Authorization"]

		if len((tokenHeadersArr)) == 0 {
			render.JSON(w, r, AppResponse{
				Success: false,
				Data: nil,
				ErrorMessage: "Token not found.",
			})
			return
		}

		tokenFull := tokenHeadersArr[0]
		tokenFull = strings.TrimSpace(tokenFull)

		tokenFullStrs := strings.Split(tokenFull, " ")

		if len(tokenFullStrs) <= 1 {
			render.JSON(w, r, AppResponse{
				Success: false,
				Data: nil,
				ErrorMessage: "No token string found.",
			})
			return
		}

		token := tokenFullStrs[1]

		jwtClaims, error := ParseJwtToken(token)

		if error != nil {
			render.JSON(w, r, AppResponse{
				Success: false,
				Data: nil,
				ErrorMessage: error.Error(),
			})
			return
		}

		ID := jwtClaims.ID

		ctx := context.WithValue(r.Context(), "ID", ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}