package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"  // 手动添加
)


type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique)"`
	Password string `gorm:"size:255;not null"`
}


func main() {

	log.Println("!!!!!!!!!!")

	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		if len(telephone) != 11 {
			//ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code":422,"msg":"手机号必须11位"})
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"手机号必须11位"})
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422,"msg":"密码不能少于6位"})
		}

		if len(name) == 0{
			name = RandomString(10)
		}

		log.Println(name, telephone, password)

		// 查询手机号
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422, "msg": "用户已存在"})
			return
		}

		newUser := User{
			Name: name,
			Telephone:telephone,
			Password:password,
		}
		db.Create(&newUser)


		// 数据验证
		// 判断手机号是否存在
		// 创建用户
		// 返回结果

		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	//r.Run() // listen and serve on 0.0.0.0:8080
	panic(r.Run())
}

func isTelephoneExist(db *gorm.DB , telephone string) bool {
	var  user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
	
}


func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfdghjklxzxcvbvnbvmQWERREYTUSDFGFJXVCNFGJRTUI")
	result := make([]byte,n)

	rand.Seed(time.Now().Unix())
	for i:= range result{
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "192.168.25.110"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil{
		panic("failed to connect database, err: "+err.Error())
	}

	db.AutoMigrate(&User{})  // 自动创建数据表
	return  db
}