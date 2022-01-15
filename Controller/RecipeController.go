package Controller

import (
	"RecipeWebApp/Domain"
	"RecipeWebApp/Service"
	"github.com/gin-gonic/gin"
)

type RecipeController interface {
	AddRecipe(ctx *gin.Context) []Domain.RecipeDomain
	GetAllRecipes() []Domain.RecipeDomain
}

type recipeController struct {
	recipeService Service.RecipeService
}

func NewRecipeController(service Service.RecipeService) RecipeController {
	return &recipeController{
		recipeService: service,
	}
}

func (controller *recipeController) AddRecipe(ctx *gin.Context) []Domain.RecipeDomain {
	var recipes []Domain.RecipeDomain
	err := ctx.BindJSON(&recipes)
	if err != nil {
		return recipes
	}
	controller.recipeService.AddRecipe(recipes)
	return recipes
}

func (controller *recipeController) GetAllRecipes() []Domain.RecipeDomain {
	return controller.recipeService.GetAllRecipes()
}
