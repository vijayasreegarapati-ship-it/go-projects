package handlers

import (
	"net/http"

	"officesync/services"

	"github.com/gin-gonic/gin"

	"strconv"
)

type CreateResourceRequest struct {
	Name         string `json:"name" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"`
}

type ResourceHandler struct {
	resourceService *services.ResourceService
}

func NewResourceHandler(resourceService *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService}
}

func (h *ResourceHandler) HandleCreateResource(c *gin.Context) {
	var req CreateResourceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	resource, err := h.resourceService.CreateResource(req.Name, req.ResourceType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resource)
}

func (h *ResourceHandler) HandleListResources(c *gin.Context) {
	resources, err := h.resourceService.ListResources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resources)
}

func (h *ResourceHandler) HandleUpdateResource(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	var req CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	resource, err := h.resourceService.UpdateResource(uint(id), req.Name, req.ResourceType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}
