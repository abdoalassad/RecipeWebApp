package Controller

import (
	"RecipeWebApp/Domain"
	"RecipeWebApp/Service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RecipeController interface {
	AddRecipe(ctx *gin.Context) *gin.Context
	GetAllRecipes(ctx *gin.Context) *gin.Context
	UpdateRecipe(ctx *gin.Context) string
	DeleteRecipe(ctx *gin.Context) string
	//SearchByTag(ctx *gin.Context) []Domain.RecipeDomain
}

type recipeController struct {
	recipeService Service.RecipeService
}

func NewRecipeController(service Service.RecipeService) RecipeController {
	return &recipeController{
		recipeService: service,
	}
}

func (controller *recipeController) AddRecipe(ctx *gin.Context) *gin.Context {
	var recipes []Domain.RecipeDomain
	err := ctx.BindJSON(&recipes)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return ctx
	}
	return controller.recipeService.AddRecipe(recipes, ctx)
}

func (controller *recipeController) GetAllRecipes(c *gin.Context) *gin.Context {
	return controller.recipeService.GetAllRecipes(c)
}

func (controller *recipeController) UpdateRecipe(ctx *gin.Context) string {
	id := ctx.Param("id")
	var recipe Domain.RecipeDomain
	err := ctx.BindJSON(&recipe)
	if err != nil {
		return "There Was An Error While Binding To JSON"
	}
	return controller.recipeService.UpdateRecipe(recipe, id)
}
func (controller *recipeController) DeleteRecipe(ctx *gin.Context) string {
	return controller.recipeService.DeleteRecipe(ctx.Param("id"))
}

//func (controller *recipeController) SearchByTag(ctx *gin.Context) []Domain.RecipeDomain {
//	tag := ctx.Query("tag")
//	return controller.recipeService.SearchByTag(tag)
//}
