// File: middleware/jwt.go

package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	config "chat-backend/Config"
	models "chat-backend/Models"
)

// JWTSecret is the secret key used to sign JWT tokens.
// In production, store this in an environment variable or secrets manager.
var JWTSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	JWTSecret = []byte(secret)

}

// Claims defines the structure of JWT claims used in this app.
// It includes the user’s ID and standard registered claims.
type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user ID.
// The token expires after `expiryDuration` (e.g., time.Hour * 24).
func GenerateToken(userID primitive.ObjectID, expiryDuration time.Duration) (string, error) {
	// Set standard claims: issuer, issued at, and expiry.
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chat-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// Authenticate middleware validates the JWT from the "Authorization" header or cookie,
// and attaches the user ID to the request context if valid.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract token from "Authorization: Bearer <token>" or from cookie named "token".
		var tokenString string

		// Check Authorization header first.
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Fallback to cookie "token"
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Missing auth token", http.StatusUnauthorized)
				return
			}
			tokenString = cookie.Value
		}

		// 2. Parse and verify token.
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 3. Optionally, ensure the user still exists in the database.
		usersColl := config.GetCollection(models.CollectionNameUser)
		filter := bson.M{"_id": claims.UserID}
		var user models.User
		err = usersColl.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// 4. Attach user ID to the request context and call next handler.
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
