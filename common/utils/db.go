package commonUtils

import (
	"flag"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func ConnectDB() *gorm.DB {
	dsnparams := flag.String("DSN", "", "")
	dsn := ""
	flag.Parse()

	if *dsnparams == "" {
		dsn = "root:120590111@tcp(192.168.3.25:3306)/zimuku?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true"
	} else {
		dsn = *dsnparams
	}
	fmt.Printf("%s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		println(err)
		panic(err)
	}
	return db
}
