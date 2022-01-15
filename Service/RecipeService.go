package Service

import (
	"RecipeWebApp/Domain"
	"github.com/rs/xid"
	"time"
)

type RecipeService interface {
	AddRecipe(recipes []Domain.RecipeDomain) []Domain.RecipeDomain
	GetAllRecipes() []Domain.RecipeDomain
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
