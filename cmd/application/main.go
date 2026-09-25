package main

import (
	"fmt"

	"github.com/CakeForKit/rsoi-lab1/internal/common/config"
	"github.com/CakeForKit/rsoi-lab1/internal/common/server"
	"github.com/CakeForKit/rsoi-lab1/internal/common/utils"
	"github.com/CakeForKit/rsoi-lab1/internal/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitLogger()

	utils.CheckedError(config.Load(&config.CoreConfig))

	fmt.Printf("config: %d\n\n", config.CoreConfig.Port)

	utils.CheckedError(server.NewHttpServer(func(router *gin.Engine) {
		utils.RegisterController(router, controller.GetPersonController)
	}).Run())
}
