package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

// route api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	//tangkap parameter di handler
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userId)
	//kalau ada error
	if err != nil {
		errorMessage := gin.H{"errors": "Error occured when try to get list of campaigns"}
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//kalau berhasil
	response := helper.APIResponse("List of Campaigns", http.StatusOK, "success", campaigns)
	c.JSON(http.StatusOK, response)
}
