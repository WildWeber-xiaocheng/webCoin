package msdb

import "gorm.io/gorm"

type DbConn interface {
	Begin()
	Rollback()
	Commit()
}

type MsDB struct {
	*gorm.DB
}
