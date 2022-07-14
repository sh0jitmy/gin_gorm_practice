package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"net/http"
	"github.com/mattn/go-gimei"
)

type httpCode struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
    Id int
    Name string
    Age int
}


var dsn =  "host=localhost port=45432 user=db_user dbname=db_name password=mypassword sslmode=disable"

func main() {
	// router init
	r := gin.Default()
	
	
	// random user name gen
	v2db, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		db,_ := v2db.DB()
		db.Close()
	}
	v2db.Migrator().CreateTable(&User{})
	r.POST("/user", func(c *gin.Context) {
		// random user name gen
		name := gimei.NewName()
		user := User{Name:name.String(),Age:22}	
		
		// create entry 
		result := v2db.Create(&user)
		if result.Error != nil {
			resCode(c, http.StatusOK, "create error")
		} else {
			resCode(c, http.StatusOK, user)
		}
	})
	r.GET("/user", func(c *gin.Context) {
		var users []User	
		
		// find entry 
		result := v2db.Find(&users)
		if result.Error != nil {
			resCode(c, http.StatusOK, "create error")
		} else {
			resCode(c, http.StatusOK, users)
		}
	})
	r.GET("/ping", func(c *gin.Context) {
		// rest check 
		resCode(c, http.StatusOK, "pong")
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func resCode(ctx *gin.Context, code int, data ...interface{}) {
	er := httpCode{
		Status:  code,
		Message: http.StatusText(code),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, er)
}
