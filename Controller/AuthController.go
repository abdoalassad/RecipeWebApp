package Controller

import (
	"RecipeWebApp/Auth"
	"RecipeWebApp/Domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController interface {
	SignIn(ctx *gin.Context) *gin.Context
	AddUser(ctx *gin.Context) *gin.Context
	RefreshToken(ctx *gin.Context) *gin.Context
	AuthMiddleware() gin.HandlerFunc
}

type authController struct {
	authService Auth.AuthService
}

func NewAuthController(authService Auth.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (controller *authController) SignIn(ctx *gin.Context) *gin.Context {
	var user Domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return ctx
	} else {
		return controller.authService.SignInHandler(ctx, user)
	}
}

func (controller *authController) AddUser(ctx *gin.Context) *gin.Context {
	var user Domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return ctx
	} else {
		return controller.authService.AddUser(ctx, user)
	}
}

func (controller *authController) RefreshToken(ctx *gin.Context) *gin.Context {
	controller.authService.RefreshToken(ctx)
	return ctx
}

func (controller *authController) AuthMiddleware() gin.HandlerFunc {
	c := controller.authService.AuthMiddleware()
	return c
}
