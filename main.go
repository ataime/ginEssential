package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // 手动添加
	"oceanlearn.teach/ginessential/common"
	"os"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	r = CollectRoute(r)
	//r.Run() // listen and serve on 0.0.0.0:8080
	panic(r.Run())
}

func InitConfig() {
	workDir,_ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println( "viper:" + viper.GetString("server.port"))  //  读取配置
}
