package worker

import (
	"github.com/hetianyi/gox/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

var (
	db         *gorm.DB
	updateLock *sync.Mutex
)

func init() {
	updateLock = new(sync.Mutex)
}

func transformNotFoundErr(err error) error {
	if err == nil {
		return nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	return err
}

func InitMysqlClientConnection(connectionString string) error {
	var _db *gorm.DB
	var err error
	logger.Info("connecting to mysql server...")
	for {
		_db, err = gorm.Open("mysql", connectionString)
		_db.LogMode(false)
		_db.DB().SetMaxOpenConns(50)
		_db.DB().SetMaxIdleConns(10)
		if err != nil {
			continue
			time.Sleep(time.Second * 5)
			logger.Info("try reconnecting to mysql server...")
		}
		break
	}
	db = _db
	//db.LogMode(true)
	return nil
}
