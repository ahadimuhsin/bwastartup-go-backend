package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Registered account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// token := helper.TokenString(64)
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

	if err != nil {
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

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := h.authService.GenerateToken(loggedInUser.ID)

	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatUser(loggedInUser, token)

	response := helper.APIResponse("Succesfully login", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	//tangkap input email
	//mapping input email ke struct input
	var input user.EmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		//ubah error validasi jadi array of string
		//dengan bantuan helper
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		//response menggunakna helper
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//struct input di-passing ke service
	//service akan manggil repository
	//dilakukan pengecekan email sudah terdaftar atau belum
	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	//jika ada error
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//atur message pengembaliannya
	data := gin.H{
		"is_available": IsEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if IsEmailAvailable {
		metaMessage = "Email is Available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	//input dari user
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//buat folder images jika tidak ada
	path := "images/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	//ambil data user dari Context gin, dari auth Middleware
	currentUser := c.MustGet("currentUser").(user.User)
	//set nilai idnya ke variabel
	userId := currentUser.ID

	//!hapus file yg lama dahulu
	//cari current root dulu
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(dir)
	// if currentUser.Avatar != "" {
	// 	location := dir + "/"+ currentUser.Avatar
	// 	oldFile := os.Remove(location)
	// 	if oldFile != nil {
	// 		fmt.Println(oldFile)
    //     	return
	// 	}
	// }

	//simpan dalam format images/id-file_name
	file_path := fmt.Sprintf("%s%d-%s-%s", path, userId, helper.TokenString(10), file.Filename)

	//di service, panggil repo
	//jwt (sementara hardcode, seakan2 user yg login ID=1)
	//repo ambil data user yg ID 1
	//repo update data user simpan lokasi file
	err = c.SaveUploadedFile(file, file_path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userService.SaveAvatar(userId, file_path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar uploaded succesfully", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
	// return
}
