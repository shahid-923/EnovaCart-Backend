package helper

import (
	"ecommerce-app/internal/domain"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {

	if len(p) < 6 {
		return "", errors.New("password length shall be of atleast 6 characters long")
	}

	//hashing
	hashedPswd, err := bcrypt.GenerateFromPassword([]byte(p), 10)

	if err != nil {
		// log actual error and report to logging tool
		return "", errors.New("password hashing failed")

	}
	return string(hashedPswd), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}

	// these are user info and is map, key is always string
	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generates and signs the JWT using the secret key.
	// If the token is modified, signature verification will fail. the response will be like knefijw.eirjijweri.j234234nks
	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("unable to sign token")
	}
	return tokenStr, nil
}

func (a Auth) VerifyPassword(plain_password string, hashed_password string) error {

	if len(plain_password) < 6 {
		return errors.New("password length shall be of atleast 6 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(plain_password))

	if err != nil {
		return errors.New("password doesn't match")
	}
	return nil
}

func (a Auth) VerifyToken(t string) (domain.User, error) {

	// Bearer t3523525jhjkk
	tokenArr := strings.Split(t, " ")

	// 1. Check if the header is completely empty
	if t == "" {
		return domain.User{}, errors.New("authorization header is missing")
	}
	// check the token length after spiliting
	if len(tokenArr) != 2 {
		return domain.User{}, errors.New("invalid authorization header format")
	}

	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("Invalid token")
	}

	tokenStr := tokenArr[1] //ie jjj.i888.8yhi
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method %v", token.Header) // checks if the token was signed using hmac
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return domain.User{}, errors.New("invalid signing method")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return domain.User{}, errors.New("token is expired")
		}

		user := domain.User{}
		user.ID = uint(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.UserType = claims["role"].(string)
		return user, nil
	}
	return domain.User{}, errors.New("token verification failed")
}

func (a Auth) Authorize() fiber.Handler {
	return func(ctx fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(401).JSON(fiber.Map{
				"message": "authorization failed",
				"reason":  "authorization header is missing",
			})
		}

		user, err := a.VerifyToken(authHeader)
		if err != nil || user.ID == 0 {
			return ctx.Status(401).JSON(fiber.Map{
				"message": "authorization failed",
				"reason":  err.Error(),
			})
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}

func (a Auth) GetCurrentUser(ctx fiber.Ctx) (domain.User, error) {

	user := ctx.Locals("user")
	if user == nil {
		return domain.User{}, errors.New("user not found in context")
	}

	u, ok := user.(domain.User)
	if !ok {
		return domain.User{}, errors.New("invalid user type")
	}

	return u, nil
}

func (a Auth) GenerateCode() (int, error) {
	return GenerateOTP()
}

func (a Auth) AuthorizeSeller() fiber.Handler {
	return func(ctx fiber.Ctx) error {

		authHeader := ctx.Get("Authorization")
		user, err := a.VerifyToken(authHeader)

		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"message": "authorization failed",
				"reason":  err,
			})
		} else if user.ID > 0 && user.UserType == domain.SELLER {
			ctx.Locals("user", user)
			return ctx.Next()
		} else {
			return ctx.Status(401).JSON(&fiber.Map{
				"message": "authorization failed",
				"reason":  errors.New("please join seller program to manage products"),
			})
		}
	}
}
