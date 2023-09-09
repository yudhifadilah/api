package middleware

import (
	"context"
	"practice-api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	collection *mongo.Collection
)

func RegisterUser(c *fiber.Ctx) error {
	var user models.User

	// Parsing body request ke struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request Body Invalid",
		})
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal Register",
		})
	}

	// Simpan user ke database MongoDB
	user.Password = string(hashedPassword)
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal Register",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil Register!",
	})
}

func LoginUser(c *fiber.Ctx) error {
	var user models.User
	inputPassword := c.FormValue("password")

	// Parsing body request ke struct User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request Body Invalid",
		})
	}

	// Cari user di database berdasarkan username
	var result models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Username atau Password Salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal login",
		})
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(inputPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username atau Password Salah",
		})
	}

	// Buat JWT token
	claims := jwt.MapClaims{
		"userId": result.Username,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key")) // Ganti dengan secret key yang aman
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal Mendapatkan Token",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

func GetMe(c *fiber.Ctx) error {
	// Mendapatkan data user yang sedang login melalui JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["userId"].(string)

	// Cari user di database berdasarkan user ID
	var userData models.User
	err := collection.FindOne(context.TODO(), bson.M{"id": userID}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User Tidak Ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal untuk Mendapatkan Data User",
		})
	}

	return c.JSON(fiber.Map{
		"user": userData,
	})
}

// Middleware untuk otentikasi JWT
func Authenticate(c *fiber.Ctx) error {
	// Mendapatkan token dari header Authorization
	authHeader := c.Get("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Header Otorisasi Salah",
		})
	}

	// Verifikasi token
	claims := new(models.JWTClaims)
	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token Invalid",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Gagal Mengautentikasi Token",
		})
	}

	if !tkn.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token Salah",
		})
	}

	// Menyimpan data user ke local context
	c.Locals("user", tkn)

	return c.Next()
}

func LogoutUser(c *fiber.Ctx) error {
	// Hapus token dari Authorization header
	c.Set("Authorization", "")

	return c.JSON(fiber.Map{
		"message": "Logout berhasil",
	})
}
