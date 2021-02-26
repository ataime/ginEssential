package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"oceanlearn.teach/ginessential/common"
	"oceanlearn.teach/ginessential/model"
	"oceanlearn.teach/ginessential/util"
)

func Login(ctx *gin.Context)  {
	// 获取参数，验证等
	// 获取参数
	DB := common.GetDB()

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
		name = util.RandomString(10)
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密失败"})
		return
	}

	newUser := model.User{
		Name: name,
		Telephone:telephone,
		Password:string(hasedPassword),
	}
	DB.Create(&newUser)


	var user model.User

	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) ; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg":"系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"data" : gin.H{"token": token},
		"msg":"登陆成功",
	})
}


func Register(ctx *gin.Context) {
	DB := common.GetDB()

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
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 查询手机号
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code":422, "msg": "用户已存在"})
		return
	}

	newUser := model.User{
		Name: name,
		Telephone:telephone,
		Password:password,
	}
	DB.Create(&newUser)


	// 数据验证
	// 判断手机号是否存在
	// 创建用户
	// 返回结果

	ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}


func isTelephoneExist(db *gorm.DB , telephone string) bool {
	var  user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false

}