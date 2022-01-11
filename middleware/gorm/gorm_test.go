package gorm

import (
	"github.com/gin-melodic/glog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"os"
	"testing"
)

const tDir = "./gorm-log"

func setUp(resetLog bool, opt Options) (db *gorm.DB, err error) {
	if resetLog {
		err = glog.InitGlobalLogger(&glog.LoggerOptions{
			MinAllowLevel:   logrus.DebugLevel,
			OutputDir:       tDir,
			FilePrefix:      "gorm-test",
			SaveDay:         1,
			ExtLoggerWriter: []io.Writer{os.Stdout},
		})
		if err != nil {
			return
		}
	}
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: New(opt),
	})
	return
}

func setDown() {
	_ = os.RemoveAll("test.db")
	_ = os.RemoveAll(tDir)
}

const createTableSql = `
CREATE TABLE IF NOT EXISTS  COMPANY(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL,
   AGE            INT     NOT NULL,
   ADDRESS        CHAR(50),
   SALARY         REAL
);`

func TestDBLogger(t *testing.T) {
	db, err := setUp(true, Options{})
	assert.Nil(t, err)
	db.Exec(createTableSql)
	db.Exec("SELECT * FROM COMPANY;")
	assert.FileExists(t, tDir+"/latest-combine-gorm-test-log")
	//db.Exec("ERROR SQL!!!")
	//assert.FileExists(t, tDir+"/latest-error-gorm-test-log")
	setDown()
}
