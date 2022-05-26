package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rampl1ng/go-todoList/controller"
)

func SetUpRouters() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/", controller.GetAllTodos)
		v1.POST("/add", controller.AddOneToDo)
		// TODO: why not use delete operation?
		v1.GET("/delete/:id", controller.DeleteOneToDo)
		v1.GET("/update/:id", controller.UpdateOneToDo)
	}

	return router
}
