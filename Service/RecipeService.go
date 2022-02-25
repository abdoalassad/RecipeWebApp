package Service

import (
	"RecipeWebApp/Domain"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type RecipeService interface {
	AddRecipe(recipes []Domain.RecipeDomain, ctx *gin.Context) *gin.Context
	GetAllRecipes(c *gin.Context) *gin.Context
	UpdateRecipe(recipe Domain.RecipeDomain, id string) string
	DeleteRecipe(recipeId string) string
	//SearchByTag(recipeTag string) []Domain.RecipeDomain
}

type recipeService struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

//func NewRecipeService() RecipeService {
//	return &recipeService{
//		Recipes: []Domain.RecipeDomain{},
//	}
//}

func NewRecipeService(ctx context.Context, database *mongo.Database, redisClient *redis.Client) RecipeService {
	return &recipeService{
		collection:  database.Collection("recipes "),
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (service *recipeService) AddRecipe(recipes []Domain.RecipeDomain, c *gin.Context) *gin.Context {
	for i := 0; i < len(recipes); i++ {
		recipes[i].ID = primitive.NewObjectID()
		recipes[i].PublishedAt = time.Now()
		_, err := service.collection.InsertOne(service.ctx, recipes[i])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return c
		}
	}
	log.Println("Removing Data From Redis")
	service.redisClient.Del("recipes")
	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
	return c
}
func (service *recipeService) GetAllRecipes(c *gin.Context) *gin.Context {
	val, err := service.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		cur, err := service.collection.Find(service.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return c
		}
		defer cur.Close(service.ctx)
		recipes := make([]Domain.RecipeDomain, 0)
		for cur.Next(service.ctx) {
			var recipe Domain.RecipeDomain
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}
		data, _ := json.Marshal(recipes)
		service.redisClient.Set("recipes", string(data), 0)
		c.JSON(http.StatusOK, gin.H{"recipes": recipes, "RedisData": false})
		return c
	} else if err != nil {
		rec := make([]Domain.RecipeDomain, 0)
		c.JSON(http.StatusBadRequest, gin.H{"error": err, "recipes": rec})
		return c
	} else {
		log.Printf("Request to Redis")
		recipes := make([]Domain.RecipeDomain, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, gin.H{"recipes": recipes, "RedisData": true})
		return c
	}
}
func (service *recipeService) UpdateRecipe(recipe Domain.RecipeDomain, id string) string {
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := service.collection.UpdateOne(service.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		return "Can't Update The Recipe"
	}
	return "Recipe Updated Successfully"
}
func (service *recipeService) DeleteRecipe(recipeId string) string {
	objectId, _ := primitive.ObjectIDFromHex(recipeId)
	_, err := service.collection.DeleteOne(service.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		return "An Error Occured While Deleting The Recipe"
	}
	return "Recipe Deleted Successfully"
}

//func (service *recipeService) SearchByTag(recipeTag string) []Domain.RecipeDomain {
//	listOfRecipes := make([]Domain.RecipeDomain, 0)
//	for i := 0; i < len(service.Recipes); i++ {
//		found := false
//		for _, tag := range service.Recipes[i].Tags {
//			if strings.EqualFold(tag, recipeTag) {
//				found = true
//			}
//		}
//		if found {
//			listOfRecipes = append(listOfRecipes, service.Recipes[i])
//		}
//		service.collection.f
//	}
//	return listOfRecipes
//}
