package main

import (
	"RecipeWebApp/Controller"
	"RecipeWebApp/Service"
	"github.com/gin-gonic/gin"
)

var (
	recipeService    Service.RecipeService       = Service.NewRecipeService()
	recipeController Controller.RecipeController = Controller.NewRecipeController(recipeService)
)

func main() {
	server := gin.Default()
	server.POST("/recipes", func(context *gin.Context) {
		context.JSON(200, recipeController.AddRecipe(context))
	})
	server.GET("/recipes", func(context *gin.Context) {
		context.JSON(200, recipeController.GetAllRecipes())
	})
	server.PUT("/recipes", func(context *gin.Context) {
		context.JSON(200, recipeController.UpdateRecipe(context))
	})
	server.DELETE("/recipes/:id", func(context *gin.Context) {
		context.JSON(200, recipeController.DeleteRecipe(context))
	})
	server.GET("/recipes/search", func(context *gin.Context) {
		context.JSON(200, recipeController.SearchByTag(context))
	})
	_ = server.Run(":8090")
}
