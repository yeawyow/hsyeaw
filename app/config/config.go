package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// โครงสร้าง JSON ส่งไป Frontend
type DbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// โหลดค่า `.env`
func LoadEnv() (*DbConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("ไม่พบไฟล์ .env")
	}

	config := &DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	}
	// ตรวจสอบว่าทุกค่ามีหรือไม่
	if config.Host == "" || config.Port == "" || config.User == "" || config.Password == "" || config.Database == "" {
		return nil, fmt.Errorf("ค่าคอนฟิกฐานข้อมูลไม่สมบูรณ์")
	}

	return config, nil
}
