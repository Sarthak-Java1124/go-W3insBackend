package routes

import (
	"e-CommerceBackend/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hey")
}

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/sign-up", controllers.SignUp)
	incomingRoutes.POST("/users/sign-in", controllers.Login)
	incomingRoutes.POST("/users/form-submit", controllers.FormSubmit)
	incomingRoutes.GET("/users/get-form-data/:address", controllers.GetFormData)

}
