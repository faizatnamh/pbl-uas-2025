package service

import (
	"os"
	"time"
	"strconv"
	"strings" 

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pbluas/app/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserService) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request"})
	}

	user, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "invalid username or password"})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{"message": "user not active"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"message": "invalid username or password"})
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret123"
	}

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_MINUTES"))
	if exp == 0 {
		exp = 60
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.RoleName,
		"exp":  time.Now().Add(time.Duration(exp) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": tokenString,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"fullname": user.FullName,
				"email":    user.Email,
				"role":     user.RoleName,
			},
		},
	})
}

func (s *UserService) Profile(c *fiber.Ctx) error {
	// Ambil token dari header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"message": "missing token"})
	}

	// Format: "Bearer <token>" → ambil token saja
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"message": "invalid token"})
	}

	tokenString = parts[1]
	secret := os.Getenv("JWT_SECRET")

	// Parse token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "invalid or expired token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":       claims["id"],
			"role":     claims["role"],
		},
	})
}
