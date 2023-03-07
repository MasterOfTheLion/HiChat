package HiChat

import (
	"HiChat/initialize"
	"HiChat/router"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	router := router.Router()
	router.Run(":8000")
}
