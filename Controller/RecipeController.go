package Controller

import (
	"RecipeWebApp/Domain"
	"RecipeWebApp/Service"
	"github.com/gin-gonic/gin"
)

type RecipeController interface {
	AddRecipe(ctx *gin.Context) []Domain.RecipeDomain
	GetAllRecipes() []Domain.RecipeDomain
	UpdateRecipe(ctx *gin.Context) Domain.RecipeDomain
	DeleteRecipe(ctx *gin.Context) string
	SearchByTag(ctx *gin.Context) []Domain.RecipeDomain
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

func (controller *recipeController) UpdateRecipe(ctx *gin.Context) Domain.RecipeDomain {
	var recipe Domain.RecipeDomain
	err := ctx.BindJSON(&recipe)
	if err != nil {
		return recipe
	}
	return controller.recipeService.UpdateRecipe(recipe)
}
func (controller *recipeController) DeleteRecipe(ctx *gin.Context) string {
	return controller.recipeService.DeleteRecipe(ctx.Param("id"))
}
func (controller *recipeController) SearchByTag(ctx *gin.Context) []Domain.RecipeDomain {
	tag := ctx.Query("tag")
	return controller.recipeService.SearchByTag(tag)
}
