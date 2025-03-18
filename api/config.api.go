package api

import (
	"context"
	"fmt"
	"os"

	"hsyncm/app/config"
	"hsyncm/app/db"
	"hsyncm/app/jwt"
	"hsyncm/functions"

	"github.com/gofiber/fiber/v2"
)

func SetupConfigAPI(router fiber.Router) {
	// สร้างกลุ่มเส้นทาง /config
	configAPI := router.Group("/config")
	{
		// เส้นทางสำหรับดึงการตั้งค่า dbcon
		configAPI.Get("/dbcon", getConfigHandler) // เปลี่ยนจาก *gin.Context เป็น *fiber.Ctx
		configAPI.Post("/regist", regist)         // เปลี่ยนจาก *gin.Context เป็น *fiber.Ctx
	}
}

func getConfigHandler(c *fiber.Ctx) error {
	// สมมติว่าคุณโหลดข้อมูลจาก env หรือฐานข้อมูล
	config, err := config.LoadEnv()
	if err != nil {
		// ถ้ามีข้อผิดพลาด, ส่งผลลัพธ์กลับเป็น JSON
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// ส่งผลลัพธ์เป็น JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"config":  config,
	})
}

func regist(c *fiber.Ctx) error {
	// สร้างตัวแปรเพื่อเก็บข้อมูลจาก body
	secretServer := os.Getenv("Secret_Key")
	var data struct {
		Hoscode    string `json:"hoscode"`
		Macaddress string `json:"macaddress"`
		Secret_Key string `json:"secret_key"`
	}

	// อ่านข้อมูลจาก body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	if err := functions.CompareData(secretServer, data.Secret_Key); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "secretkey ไม่ถูกต้อง",
			// "test":    secretServer + data.Secret_Key,
		})
	}
	// เชื่อมต่อกับฐานข้อมูลเพื่อทำการตรวจสอบ
	var exists bool
	provcode := os.Getenv("PROVINCE_CODE")
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM chospital WHERE hoscode=$1 AND provcode=$2)", data.Hoscode, provcode).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Database query error2: %v", err),
		})
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if exists {

		var exists_hoscode bool
		err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM registers WHERE hoscode=$1 )", data.Hoscode).Scan(&exists_hoscode)
		// ถ้ามีข้อมูล
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("Database query error3: %v", err),
			})
		}
		if exists_hoscode {
			var exists_hcmc bool
			err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM registers WHERE hoscode=$1 AND macaddress=$2 )", data.Hoscode, data.Macaddress).Scan(&exists_hcmc)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to insert new record: %v", err),
				})
			}
			if exists_hcmc {
				token, err := jwt.GenerateToken(data.Hoscode)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": data.Macaddress,
						"error":   fmt.Sprintf("Failed to insert new record: %v", err),
					})
				}

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"token": token,
					"error": "",
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": "หน่วยงานของท่านมีเครื่องที่ลงทะเบียนแล้ว",
			})
		} else if !exists_hoscode {
			// สามารถเพิ่มข้อมูลได้ที่นี่ (Optional)
			_, err := db.Pool.Exec(context.Background(), "INSERT INTO registers (hoscode, macaddress) VALUES ($1, $2)", data.Hoscode, data.Macaddress)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error":   fmt.Sprintf("Failed to insert new record: %v", err),
				})
			}
			token, err := jwt.GenerateToken(data.Hoscode)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": data.Macaddress,
					"error":   fmt.Sprintf("Failed to insert new record: %v", err),
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"token": token,
				"error": "",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "",
		})
	} else {
		// ถ้าไม่มีข้อมูลในหน่วยบริการ

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"error":   "ไม่พบหน่วยบริการในจังหวัดนี้",
		})
	}
}
