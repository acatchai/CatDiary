package repository

import (
	"log"

	"github.com/acatchai/catdiary/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接并自动迁移表结构
func InitDB(dsn string) {
	var err error
	// 连接Mysql
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境打印SQL语句方便调试
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 执行自动迁移
	// GORM 会自动检查 "user" 和 "diary" 表是否存在，
	// 如果不存在则创建；如果存在但缺字段，会自动追加新字段。
	err = DB.AutoMigrate(
		&model.User{},
		&model.Diary{},
	)
	if err != nil {
		log.Fatalf("自动迁移表结构失败: %v", err)
	}
	log.Println("数据库连接成功，表结构迁移完成！🚀")
}
