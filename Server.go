package main

import (
	"RecipeWebApp/Auth"
	"RecipeWebApp/Config"
	"RecipeWebApp/Controller"
	"RecipeWebApp/Service"
	"github.com/gin-gonic/gin"
)

var (
	context, db                                  = Config.InitializingMongoDB()
	redisDB                                      = Config.InitializingRedisCache()
	authService      Auth.AuthService            = Auth.NewAuthHandler(context, db)
	authController   Controller.AuthController   = Controller.NewAuthController(authService)
	recipeService    Service.RecipeService       = Service.NewRecipeService(context, db, redisDB)
	recipeController Controller.RecipeController = Controller.NewRecipeController(recipeService)
)

func main() {
	server := gin.Default()

	server.POST("/signin", func(context *gin.Context) {
		authController.SignIn(context)
	})
	server.POST("/adduser", func(context *gin.Context) {
		authController.AddUser(context)
	})
	server.GET("/recipes", func(context *gin.Context) {
		recipeController.GetAllRecipes(context)
	})
	server.POST("/refresh", func(context *gin.Context) {
		authController.RefreshToken(context)
	})
	authorized := server.Group("/")
	authorized.Use(authController.AuthMiddleware())
	{
		authorized.POST("/recipes", func(context *gin.Context) {
			recipeController.AddRecipe(context)
		})
		authorized.PUT("/recipes/:id", func(context *gin.Context) {
			recipeController.UpdateRecipe(context)
		})
		authorized.DELETE("/recipes/:id", func(context *gin.Context) {
			recipeController.DeleteRecipe(context)
		})
	}

	//server.GET("/recipes/search", func(context *gin.Context) {
	//	context.JSON(200, recipeController.SearchByTag(context))
	//})
	_ = server.Run(":8090")
}
