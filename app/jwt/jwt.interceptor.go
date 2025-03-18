package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

// ฟังก์ชัน init() โหลด Secret Key ตอนเริ่มต้นโปรแกรม
func init() {
	secret := "HISSYNg-9l6-4krmuj8" // ควรเปลี่ยนใน Production

	secretKey = []byte(secret)
}

// ฟังก์ชันสร้าง Token
func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // หมดอายุใน 1 ชั่วโมง
		"iat":     time.Now().Unix(),                    // เวลาสร้าง Token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ฟังก์ชันตรวจสอบ JWT และต่ออายุถ้าหมดอายุ
func JwtVerify(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"result":  "nok",
			"message": "Missing token",
		})
	}

	// แยก "Bearer " ออกจาก tokenString
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"result":  "nok",
			"message": "Invalid token format",
		})
	}
	tokenString = parts[1]

	// Parse Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	// 🔹 ตรวจสอบว่าข้อผิดพลาดเป็น `ErrTokenExpired`
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"result":  "nok",
				"message": "Token expired, please refresh",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"result":  "nok",
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	// ถ้า Token ใช้งานได้
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}

	// Token ไม่ valid
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"result":  "nok",
		"message": "Invalid token",
	})
}
