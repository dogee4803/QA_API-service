package main

import (
	"log"
	"os/exec"
	"qna-api/internal/config"
)

func main() {
	cfg := config.Load()
	
	dsn := "host=" + cfg.DBHost + 
	       " user=" + cfg.DBUser + 
	       " password=" + cfg.DBPassword + 
	       " dbname=" + cfg.DBName + 
	       " port=" + cfg.DBPort + 
	       " sslmode=disable"
	
	cmd := exec.Command("goose", "-dir", "internal/database/migrations", "postgres", dsn, "up")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Ошибка миграций: %v\nOutput: %s", err, string(output))
	}
	
	log.Printf("Миграции успешно применены: %s", string(output))
}