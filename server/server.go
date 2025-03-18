package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"hsyncm/api" // ทำให้แน่ใจว่า import แพ็กเกจ api ให้ถูกต้อง

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Server ฟังก์ชันเริ่มเซิร์ฟเวอร์และจัดการ Shutdown
func Server(ctx context.Context) {
	app := fiber.New()

	// Middleware
	app.Use(logger.New()) // Log request

	// สร้าง api group สำหรับ backend routes
	apiGroup := app.Group("/api") // ใช้ /api เป็น prefix สำหรับ backend
	api.Setup(apiGroup)           // เรียกฟังก์ชัน Setup

	// เสิร์ฟ Static Files สำหรับ frontend
	app.Static("/assets", "./static/assets") // เสิร์ฟไฟล์จาก ./static/assets

	// เสิร์ฟ React HTML ไฟล์ สำหรับ frontend
	// เส้นทางอื่น ๆ ที่ไม่ใช่ /api/ จะถือว่าเป็น frontend
	app.Use(func(c *fiber.Ctx) error {
		if strings.Contains(c.Path(), "/api") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "API route not found"})
		}
		return c.SendFile("./static/index.html")
	})

	// รันเซิร์ฟเวอร์ใน Goroutine
	go func() {
		fmt.Println("Server running on http://localhost:8081")
		if err := app.Listen(":8081"); err != nil {
			log.Fatalf("Server error: %s\n", err)
		}
	}()
	// สร้าง channel สำหรับการจับสัญญาณจากระบบปฏิบัติการ (เช่น CTRL+C หรือการหยุดโปรเซส)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// รอรับสัญญาณ Shutdown
	<-ctx.Done()
	fmt.Println("\nShutting down server...")

	// ปิดเซิร์ฟเวอร์อย่างปลอดภัยภายใน 5 วินาที
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server exited")
}
