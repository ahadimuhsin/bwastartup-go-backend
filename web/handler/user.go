package handler

import (
	"bwastartup/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})
}

func (h *userHandler) New(c *gin.Context){
	c.HTML(http.StatusOK, "user_new.html", nil)
}

func (h *userHandler) Create(c *gin.Context){
	//tampung input
	var input user.FormCreateUserInput

	err := c.ShouldBind(&input)

	if err != nil{
		//menampilkan kembali data inputan ke halaman input
		//jika ada data form yg error
		input.Error = err
		c.HTML(http.StatusOK, "user_new.html", input)
		return;
	}

	registerInput := user.RegisterUserInput{}
	registerInput.Name = input.Name
	registerInput.Email = input.Email
	registerInput.Password = input.Password
	registerInput.Occupation = input.Occupation

	_, err = h.userService.RegisterUser(registerInput)

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return;
	}

	c.Redirect(http.StatusFound, "/users")
}

func (h *userHandler) Edit(c *gin.Context){
	//tangkap parameter
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	//ambil user berdasarkan ID
	registeredUser, err := h.userService.GetUserById(id)

	if err != nil{
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	}

	input := user.FormUpdateUserInput{}
	input.ID = registeredUser.ID
	input.Name = registeredUser.Name
	input.Email = registeredUser.Email
	input.Occupation = registeredUser.Occupation

	c.HTML(http.StatusOK, "user_edit.html", input)
}

func (h *userHandler) Update(c *gin.Context){
	//tangkap parameter
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input user.FormUpdateUserInput

	err := c.ShouldBind(&input)

	if err != nil{
		//menampilkan kembali data inputan ke halaman input
		//jika ada data form yg error
		input.Error = err
		c.HTML(http.StatusOK, "user_edit.html", input)
		return;
	}

	input.ID = id

	_, err = h.userService.UpdateUser(input)

	if err != nil{
		//skip
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/users")
}