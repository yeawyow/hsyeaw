package api

import (
	"context"
	"fmt"
	"sync"

	"hsyncm/app/db"  // ใส่ package สำหรับ db ให้ถูกต้อง
	"hsyncm/app/jwt" // ใส่ package สำหรับ jwt ให้ถูกต้อง

	"github.com/gofiber/fiber/v2"
)

func SetupReceiveAPI(router fiber.Router) {
	// สร้างกลุ่มเส้นทาง /receive
	receiveAPI := router.Group("/receive")
	{
		// เส้นทางสำหรับรับข้อมูล
		receiveAPI.Post("/", jwt.JwtVerify, ReceiverHandler)
	}
}

// Struct สำหรับข้อมูลที่ต้องการส่งไปประมวลผล
type RequestData struct {
	Text1 string `json:"text1"` // แก้ไข field เป็น Text1 แทน text1
}

func worker(id int, jobs <-chan RequestData, wg *sync.WaitGroup) {
	defer wg.Done() // แจ้งว่า worker นี้ทำงานเสร็จแล้ว
	for data := range jobs {
		// บันทึกข้อมูลลงใน PostgreSQL
		if err := saveToDatabase(data); err != nil {
			fmt.Printf("Error inserting data: %v\n", err)
		} else {
			fmt.Printf("Worker %d: Inserted text1=%s\n", id, data.Text1)
		}
	}
}

// ฟังก์ชันบันทึกข้อมูลลงใน PostgreSQL
func saveToDatabase(data RequestData) error {
	// สร้าง query เพื่อ insert ข้อมูลลงในฐานข้อมูล
	query := "INSERT INTO test (text1) VALUES ($1)"
	_, err := db.Pool.Exec(context.Background(), query, data.Text1)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}
	return nil
}

// Handler สำหรับรับคำขอจาก client
func ReceiverHandler(c *fiber.Ctx) error {
	// รับข้อมูลจาก request body
	var input RequestData
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// สร้าง channel สำหรับส่งข้อมูลไปยัง worker
	jobs := make(chan RequestData, 100) // buffer 100 jobs
	var wg sync.WaitGroup

	// เริ่มงาน worker
	for i := 1; i <= 5; i++ { // สร้าง worker 5 ตัว
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// ส่งข้อมูลไปยัง worker
	jobs <- input

	// ปิดช่องทาง jobs หลังจากส่งข้อมูลแล้ว
	close(jobs)

	// รอให้ worker ประมวลผลเสร็จ
	wg.Wait()

	// ส่งการตอบกลับให้ client ว่ากระบวนการเริ่มต้นแล้ว
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Worker started processing data",
	})
}
