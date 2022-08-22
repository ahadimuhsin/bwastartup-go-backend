package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	//map input dari user ke struct RegisterUserInput
	//kemudian teruskan sebagai parameter service
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		//ubah error validasi jadi array of string
		//dengan bantuan helper
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		//response menggunakna helper
		response := helper.APIResponse("Registered account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Registered account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken()
	token := helper.TokenString(64)
	formatter := user.FormatUser(newUser, token)
	//response menggunakna helper
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) LoginUser(c *gin.Context) {
	//user memasukkan input: email dan password
	//input ditangkap handler
	
	//di service mencari dengan bantuan repository user dengan email yg dimasukkan
	//mencocokkan password

	var input user.LoginInput
	err := c.ShouldBindJSON(&input)

	if err != nil{
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		//response menggunakna helper
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return 
	}

	//mapping dari input user ke input struct
	//input struct ke service
	//di service mencari dengan bantuan repository user dengan email yg dimasukkan
	//mencocokkan password
	loggedInUser, err := h.userService.Login(input)

	if err != nil{
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return 
	}

	formatter := user.FormatUser(loggedInUser, helper.TokenString(64))

	response := helper.APIResponse("Succesfully login", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
