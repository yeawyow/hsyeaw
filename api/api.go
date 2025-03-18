package api

import "github.com/gofiber/fiber/v2"

// Setup รวมการตั้งค่า API ทั้งหมด
func Setup(apiGroup fiber.Router) {
	// กำหนด API สำหรับ user
	SetupReceiveAPI(apiGroup)

	// กำหนด API สำหรับ config
	SetupConfigAPI(apiGroup)
}
