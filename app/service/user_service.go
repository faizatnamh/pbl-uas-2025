package service

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"pbluas/app/models"
	"pbluas/app/repository"
)

type UserService struct {
	Repo     repository.UserRepository
	PermRepo *repository.PermissionRepository // --- DITAMBAHKAN
}

func NewUserService(repo repository.UserRepository, perm *repository.PermissionRepository) *UserService {
	return &UserService{
		Repo:     repo,
		PermRepo: perm, 
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LOGIN USER
func (s *UserService) Login(c *fiber.Ctx) error {
	var req LoginRequest

	
	// 400 – Invalid Request Body
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":        400,
			"message":     "Bad Request",
			"description": "Invalid request body",
		})
	}


	// 401 – Username not found
	
	user, err := s.Repo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid username or password",
		})
	}


	// 403 – User inactive
	
	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"code":        403,
			"message":     "Forbidden",
			"description": "User is not active",
		})
	}

	
	// 401 – Wrong password

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid username or password",
		})
	}


	//  AMBIL LIST PERMISSIONS SESUAI ROLE
	
	permissions, err := s.PermRepo.GetPermissionsByRole(user.RoleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load permissions",
		})
	}

	
	//  GENERATE ACCESS TOKEN
	
	secret := os.Getenv("JWT_SECRET")
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
	accessToken, _ := token.SignedString([]byte(secret))

	
	// REFRESH TOKEN
	
	refreshClaims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, _ := refresh.SignedString([]byte(secret))

	
	//  RESPONSE 
	authUser := models.AuthUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        user.RoleName,
		Permissions: permissions,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user":         authUser,
		},
	})
}


// PROFILE USER
func (s *UserService) Profile(c *fiber.Ctx) error {
	tokenHeader := c.Get("Authorization")

	if tokenHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Missing token",
		})
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid token format",
		})
	}

	tokenStr := parts[1]
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid or expired token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":   claims["id"],
			"role": claims["role"],
		},
	})
}

// REFRESH TOKEN
func (s *UserService) Refresh(c *fiber.Ctx) error {
	type RefreshReq struct {
		RefreshToken string `json:"refreshToken"`
	}

	var req RefreshReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":        400,
			"message":     "Bad Request",
			"description": "Invalid request body",
		})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{
			"code":        400,
			"message":     "Bad Request",
			"description": "refreshToken is required",
		})
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid or expired refresh token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	user, err := s.Repo.FindByUserID(userId)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":        404,
			"message":     "Not Found",
			"description": "User not found",
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"code":        403,
			"message":     "Forbidden",
			"description": "User is not active",
		})
	}

	permissions, err := s.PermRepo.GetPermissionsByRole(user.RoleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to load permissions"})
	}

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_MINUTES"))
	newClaims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.RoleName,
		"exp":  time.Now().Add(time.Duration(exp) * time.Minute).Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	accessToken, _ := newToken.SignedString([]byte(secret))

	authUser := models.AuthUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        user.RoleName,
		Permissions: permissions,
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        accessToken,
			"refreshToken": req.RefreshToken,
			"user":         authUser,
		},
	})
}

// LOGOUT USER
func (s *UserService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
