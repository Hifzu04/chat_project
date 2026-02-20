// File: middleware/jwt.go

package middleware

import (
	"context"
	"fmt"
	"net/http"
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
var JWTSecret = []byte("uaRxLXE5i6y3/gAM3w/XYh44lj7W0r9hQEF00DkgusQ=")

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

//. It runs before every protected request.(routes that need authentication). like GetMessages =>only possible when there is a valid token.
//This  sits between the User and the protected Routes.
// and attaches the user ID to the request context if valid.

// Middleware Pattern: It takes a next handler (the destination, e.g., GetMessages) and returns a new handler that wraps logic around it.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract token from "Authorization: Bearer <token>" or from cookie named "token".
		var tokenString string
		//Mobile Apps usually send tokens in the Header (Authorization: Bearer <token>).
		//Web Browsers usually send tokens in Cookies (HttpOnly).
		// The Check: If we look in both places and find nothing, we stop immediately (401 Unauthorized). The user cannot pass.
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

		// 2. Parse and verify token
		claims := &Claims{}
		//ParseWithClaims: Attempts to read the token using your JWTSecret.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			//for security: Ensure the signing method is HMAC and specifically HS256.
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return JWTSecret, nil
		})

		//If the secret key doesn't match, OR if the token has expired (past 48 hours), token.Valid will be false. We reject the request.
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 3. Optionally, ensure the user still exists in the database.
		//why Imagine an admin bans/deletes a user at 2:00 PM.
		//The user's token is technically still valid until 2:00 PM tomorrow.
		usersColl := config.GetCollection(models.CollectionNameUser)
		filter := bson.M{"_id": claims.UserID}
		var user models.User
		err = usersColl.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			http.Error(w, "User not found in the database checked by Authenticate", http.StatusUnauthorized)
			return
		}

		// 4. Attach user ID to the request context and call next handler means .
		//context.WithValue(parent, key, value): Input: The old context, a label ("userID"), and the data (claims.UserID).
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		//Output: A NEW context object. It does not change the old one. This is why we assign it to a new variable ctx.
		//r.WithContext(ctx): This method creates a shallow copy of the HTTP request r.
		//It keeps the same URL, Body, and Headers, but it swaps out the context for your new ctx.
		//next.ServeHTTP(w, ...): next: This is the next function in the chain (e.g., SendMessage).
		//ServeHTTP: This is simply the command to "Run that function."
		//We pass it the w (so it can write a response) and the modified request (so it can read the UserID).
		next.ServeHTTP(w, r.WithContext(ctx))
		
         //eg
		//func SendMessage(w http.ResponseWriter, r *http.Request) {
		// This works ONLY because the middleware passed the modified 'r'
		// userID := r.Context().Value("userID")
		//}
		//If you missed that line in the middleware, r.Context().Value("userID") would return nil (nothing), and your app would say "Unauthorized."
	})
}

