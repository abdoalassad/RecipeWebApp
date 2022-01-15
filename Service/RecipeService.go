package Service

import (
	"RecipeWebApp/Domain"
	"github.com/rs/xid"
	"strings"
	"time"
)

type RecipeService interface {
	AddRecipe(recipes []Domain.RecipeDomain) []Domain.RecipeDomain
	GetAllRecipes() []Domain.RecipeDomain
	UpdateRecipe(recipe Domain.RecipeDomain) Domain.RecipeDomain
	DeleteRecipe(recipeId string) string
	SearchByTag(recipeTag string) []Domain.RecipeDomain
}

type recipeService struct {
	Recipes []Domain.RecipeDomain
}

func NewRecipeService() RecipeService {
	return &recipeService{
		Recipes: []Domain.RecipeDomain{},
	}
}

func (service *recipeService) AddRecipe(recipes []Domain.RecipeDomain) []Domain.RecipeDomain {
	for i := 0; i < len(recipes); i++ {
		recipes[i].ID = xid.New().String()
		recipes[i].PublishedAt = time.Now()
		service.Recipes = append(service.Recipes, recipes[i])
	}
	return service.Recipes
}
func (service *recipeService) GetAllRecipes() []Domain.RecipeDomain {
	return service.Recipes
}
func (service *recipeService) UpdateRecipe(recipe Domain.RecipeDomain) Domain.RecipeDomain {
	found := false
	for i := 0; i < len(service.Recipes); i++ {
		if recipe.ID == service.Recipes[i].ID {
			service.Recipes[i] = recipe
			found = true
		}
	}
	if found == false {
		service.Recipes = append(service.Recipes, recipe)
	}
	return recipe
}
func (service *recipeService) DeleteRecipe(recipeId string) string {
	found := false
	for i := 0; i < len(service.Recipes); i++ {
		if recipeId == service.Recipes[i].ID {
			service.Recipes = append(service.Recipes[:i], service.Recipes[i+1:]...)
			found = true
		}
	}
	if found == false {
		return "Recipe Not Found"
	}
	return "Recipe Deleted Successfully"
}
func (service *recipeService) SearchByTag(recipeTag string) []Domain.RecipeDomain {
	listOfRecipes := make([]Domain.RecipeDomain, 0)
	for i := 0; i < len(service.Recipes); i++ {
		found := false
		for _, tag := range service.Recipes[i].Tags {
			if strings.EqualFold(tag, recipeTag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, service.Recipes[i])
		}
	}
	return listOfRecipes
}
