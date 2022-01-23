package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	tokenHelpers "github.com/NAziz21/TheSystemOfOnlineCourses/pkg/token"
	_ "github.com/golang-jwt/jwt/v4"
)


var ErrNoAuthentication = errors.New("no authentication")



// A variable that will be the key by which the value will be added.
var authenticationContextKey = &contextKey{"authentication context"}


// Non-exportable type.
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}




func Auth() func(http.Handler) http.Handler{
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			
			tokenString := request.Header.Get("Authorization")
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			tokenString = strings.TrimSpace(tokenString)
		
		
			claims, err := tokenHelpers.ValidateToken(tokenString) 
			if err != nil {
				log.Println(err, "Not Authorized")
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if err != nil {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			id, _ := strconv.ParseInt(claims.Id, 10, 0)
			ctx := context.WithValue(request.Context(), authenticationContextKey, id)
			request = request.WithContext(ctx)
		
			handler.ServeHTTP(writer, request)
		}) 
	}
}



// Athuntecation - helper function, to extract value from context.
func AuthenticationJWT(ctx context.Context) (int64, error) {

	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}

	return 0, ErrNoAuthentication
}






















