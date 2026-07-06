package gorm_client

import (
	"go-backend/internal/common/env"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(env *env.Env) *gorm.DB {
	db, err := gorm.Open(mysql.Open(env.DatabaseUrl), &gorm.Config{})

	if err != nil {
		log.Fatalf("[GORM] failed opening connection to mysql: %v", err)
	}

	err = db.Raw("SELECT 1 + 1").Error

	if err != nil {
		log.Fatalf("[GORM] failed opening connection to mysql: %v", err)
	}

	// fmt.Println("[GORM] Connection to my SQL Successfully")

	return db
}
