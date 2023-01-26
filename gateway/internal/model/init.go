package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"powernotes-server/gateway/internal/config"
)

var DB *gorm.DB

func InitDB(config config.DBConfig) {
	var err error
	var dialector gorm.Dialector
	switch config.Driver {
	case "sqlite":
		dialector = sqlite.Open(config.DSN)
	case "mysql":
		dialector = mysql.Open(config.DSN)
	default:
		panic("unsupported db driver")
	}
	DB, err = gorm.Open(dialector)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&Note{}, &Flow{}, &FlowNoteRelation{})
	if err != nil {
		panic(err)
	}
}
