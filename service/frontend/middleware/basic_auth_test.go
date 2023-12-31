package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicAuthFn(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	fakeAuthToken := AuthToken{
		Token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikpva",
	}
	testCase := []struct {
		name       string
		authHeader string
		authToken  *AuthToken
		httpStatus int
	}{
		{
			name:       "auth header set, auth token set",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikpva",
			authToken:  &fakeAuthToken,
			httpStatus: http.StatusOK,
		},
		{
			name:       "auth header set, auth token unset",
			authHeader: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ikpva",
			authToken:  nil,
			httpStatus: http.StatusUnauthorized,
		},
		{
			name:       "auth header unset, auth token set",
			authHeader: "",
			authToken:  &fakeAuthToken,
			httpStatus: http.StatusUnauthorized,
		},
		{
			name:       "auth header unset, auth token unset",
			authHeader: "",
			authToken:  nil,
			httpStatus: http.StatusUnauthorized,
		},
	}
	// incorrectCreds triggers HTTP 401 Unauthorized upon basic auth
	incorrectCreds := map[string]string{
		"INCORRECT_USERNAME": "INCORRECT_PASSWORD",
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)
			r.Header.Add("Authorization", tc.authHeader)
			w := httptest.NewRecorder()
			authToken = tc.authToken
			BasicAuth(
				"restricted",
				incorrectCreds,
			)(testHandler).ServeHTTP(w, r)
			require.Equal(t, tc.httpStatus, w.Result().StatusCode)
		})
	}
}