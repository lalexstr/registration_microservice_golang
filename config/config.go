package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port       string
	JWTSecret  []byte
	JWTTTLMin  int
	SQLitePath string
)

func init() {
	execDir, err := os.Executable()
	if err != nil {
		log.Fatalf("❌ Cannot get executable path: %v", err)
	}
	envPath := filepath.Join(filepath.Dir(execDir), ".env")

	// Загружаем .env
	if err := godotenv.Load(envPath); err != nil {
		_ = godotenv.Load(".env")
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(JWTSecret) == 0 {
		log.Fatal("❌ Missing JWT_SECRET in environment")
	}

	ttlStr := os.Getenv("JWT_TTL_MINUTES")
	if ttlStr == "" {
		ttlStr = "60"
	}

	ttlInt, err := strconv.Atoi(ttlStr)
	if err != nil {
		log.Fatalf("❌ Invalid JWT_TTL_MINUTES value: %v", err)
	}
	JWTTTLMin = ttlInt

	SQLitePath = os.Getenv("SQLITE_PATH")
	if SQLitePath == "" {
		SQLitePath = "auth.db"
	}

	log.Printf("✅ Config loaded: PORT=%s | TTL=%d min | DB=%s", Port, JWTTTLMin, SQLitePath)
}
