package Auth

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"strings"
	"time"

	"RecipeWebApp/Domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type AuthService interface {
	SignInHandler(ctx *gin.Context, user Domain.User) *gin.Context
	AddUser(ctx *gin.Context, user Domain.User) *gin.Context
	RefreshToken(c *gin.Context) *gin.Context
	AuthMiddleware() func(c *gin.Context)
}

type authService struct {
	collection *mongo.Collection
	ctx        context.Context
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
	Password string `json:"password"`
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(ctx context.Context, database *mongo.Database) AuthService {
	return &authService{
		collection: database.Collection("user"),
		ctx:        ctx,
	}
}

func (auth *authService) SignInHandler(c *gin.Context, user Domain.User) *gin.Context {
	hash := sha512.New()
	hash.Write([]byte(user.Password))

	hashedPassword := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	cur := auth.collection.FindOne(auth.ctx, bson.M{
		"username": user.Username,
		//"password": user.Password,
		"password": hashedPassword,
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return c
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Password: user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte("hiabdo"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return c
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, gin.H{"JWTToken": jwtOutput})
	return c
}

func (auth *authService) AddUser(c *gin.Context, user Domain.User) *gin.Context {
	hash := sha512.New()
	hash.Write([]byte(user.Password))

	hashedPassword := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	us, err := auth.collection.InsertOne(auth.ctx, bson.M{"username": user.Username, "password": hashedPassword})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could Not Save User"})
	} else {
		c.JSON(http.StatusOK, gin.H{"user": us})
	}
	return c
}

func (auth *authService) RefreshToken(c *gin.Context) *gin.Context {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return c
	}
	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return c
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
		return c
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString("hiabdo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return c
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, gin.H{"JWTToken": jwtOutput})
	return c
}

func (auth *authService) AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		tokenValue := c.GetHeader("Authorization")
		claims := &Claims{}
		tokenValue = strings.Trim(tokenValue, "Bearer ")
		parts := strings.Split(tokenValue, ".")
		headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
		token := &jwt.Token{Raw: tokenValue}
		err = json.Unmarshal(headerBytes, &token.Header)

		tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("hiabdo"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
