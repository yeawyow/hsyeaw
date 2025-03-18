package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hsyncm/app/db"
	"hsyncm/server"

	"github.com/joho/godotenv"
)

func main() {
	// สร้าง context สำหรับจัดการการยกเลิกการทำงาน
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop() // ทำการหยุดเมื่อได้รับสัญญาณ

	// เริ่มต้นการเชื่อมต่อฐานข้อมูล
	if err := db.InitDB("hissync", "Sync6530", "192.168.0.90", "5432", "hissync"); err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB() // ปิดการเชื่อมต่อเมื่อทำงานเสร็จ

	// ทดสอบการเชื่อมต่อฐานข้อมูล
	if err := db.PingDB(); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	// แสดงข้อความเมื่อเชื่อมต่อฐานข้อมูลสำเร็จ
	fmt.Println("Successfully connected to the database!")

	// เรียกใช้เซิร์ฟเวอร์เพื่อเริ่มทำงาน
	server.Server(ctx)
}
