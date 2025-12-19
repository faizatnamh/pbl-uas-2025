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
	"pbluas/config"
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

// Login godoc
// @Summary Login user
// @Description Autentikasi user dan generate JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
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


// Profile godoc
// @Summary Get user profile
// @Description Ambil data user dari JWT
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/profile [get]
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

// Refresh godoc
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Refresh token"
// @Router /auth/refresh [post]
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

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Router /auth/logout [post]
func (s *UserService) Logout(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "Missing token",
		})
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid token format",
		})
	}

	token := parts[1]

	// blacklist token
	config.BlacklistToken(token)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Router /users [get]
func (s *UserService) GetAllUsers(c *fiber.Ctx) error {
    users, err := s.Repo.GetAllUsers()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to retrieve users",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   users,
    })
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [get]
func (s *UserService) GetUserByID(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := s.Repo.FindByUserID(id)

    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "User not found",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "data":   user,
    })
}

// CreateUser godoc
// @Summary Create new user
// @Description Admin creates a new user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CreateUserRequest true "Create user payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users [post]
func (s *UserService) CreateUser(c *fiber.Ctx) error {
    var req models.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Invalid request body",
        })
    }

    // Hash password
    hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    newUser := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashed),
        FullName:     req.FullName,
        RoleID:       req.RoleID,
        IsActive:     true,
    }

    if err := s.Repo.CreateUser(newUser); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to create user",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "User created successfully",
    })
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user data by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body models.UpdateUserRequest true "Update user payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [put]
func (s *UserService) UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")

    user, err := s.Repo.FindByUserID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "User not found",
        })
    }

    var req models.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Invalid request body",
        })
    }

    // Update fields
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.FullName != "" {
        user.FullName = req.FullName
    }
    if req.RoleID != "" {
        user.RoleID = req.RoleID
    }
    if req.IsActive != nil {
        user.IsActive = *req.IsActive
    }
    if req.Password != "" {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        user.PasswordHash = string(hashed)
    }

    if err := s.Repo.UpdateUser(user); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to update user",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "User updated successfully",
    })
}

/// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (s *UserService) DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")

    err := s.Repo.DeleteUser(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to delete user",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "User deleted successfully",
    })
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update role of a user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param body body models.UpdateUserRoleRequest true "Update role payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id}/role [put]
func (s *UserService) UpdateUserRole(c *fiber.Ctx) error {
    id := c.Params("id")

    var req models.UpdateUserRoleRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Invalid request body",
        })
    }

    if err := s.Repo.UpdateUserRole(id, req.RoleID); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to update role",
        })
    }

    return c.JSON(fiber.Map{
        "status": "success",
        "message": "Role updated successfully",
    })
}

