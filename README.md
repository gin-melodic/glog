# gin-logger

This toolkit provides rotation logging capabilities.

`GLog` contain some useful middleware for:

- [gingonic/gin](https://github.com/gin-gonic/gin)
- [go-gorm/gorm](https://github.com/go-gorm/gorm)

# Best Practices

```go
package main

import (
	"github.com/gin-melodic/glog"
	"github.com/gin-gonic/gin"
	gingonicLogger "github.com/gin-melodic/glog/middleware/gingonic"
	gormLogger "github.com/gin-melodic/glog/middleware/gorm"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// logger init
	loggerConfig := &glog.LoggerOptions{
		MinAllowLevel:   logrus.InfoLevel,
		HighPerformance: false,
		OutputDir:       "./logs",
		FilePrefix:      "demo-project",
		SaveDay:         30,
		ExtLoggerWriter: []io.Writer{os.Stdout},
	}
	// you can reset config here, 
	// e.g. Disable os.Stdout logger writer in prod.
	//if env == 'production' {
	//    loggerConfig.ExtLoggerWriter = nil
	//}
	if err := glog.InitGlobalLogger(loggerConfig); err != nil {
		panic(err)
	}

	// add gorm middleware
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: gormLogger.New(gormLogger.Options{}),
	})
	
	// add gin middleware
	router := gin.New()
	router.Use(gingonicLogger.InjectLogger(&gingonicLogger.Options{
		BodyMaxSize:          500,
	}))
}
```

# License

Apache-2.0 License

