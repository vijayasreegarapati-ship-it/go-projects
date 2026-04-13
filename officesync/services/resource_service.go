package services

import (
	"errors"
	"strings"

	"officesync/models"

	"gorm.io/gorm"
)

type ResourceService struct {
	DB *gorm.DB
}

func NewResourceService(db *gorm.DB) *ResourceService {
	return &ResourceService{DB: db}
}

func (s *ResourceService) CreateResource(name, resourceType string) (*models.Resource, error) {
	name = strings.TrimSpace(name)
	resourceType = strings.TrimSpace(resourceType)

	if name == "" || resourceType == "" {
		return nil, errors.New("name and resource type are required")
	}

	if resourceType != "desk" && resourceType != "meeting_room" {
		return nil, errors.New("invalid resource type: must be 'desk' or 'meeting_room'")
	}

	resource := &models.Resource{
		Name:         name,
		ResourceType: resourceType,
		Status:       "active",
	}

	if result := s.DB.Create(resource); result.Error != nil {
		return nil, errors.New("failed to create resource in the database")
	}

	return resource, nil
}

func (s *ResourceService) ListResources() ([]models.Resource, error) {
	var resources []models.Resource

	if result := s.DB.Find(&resources); result.Error != nil {
		return nil, errors.New("failed to fetch resources")
	}

	return resources, nil
}

func (s *ResourceService) UpdateResource(id uint, name, resourceType string) (*models.Resource, error) {
	name = strings.TrimSpace(name)
	resourceType = strings.TrimSpace(resourceType)

	if name == "" || resourceType == "" {
		return nil, errors.New("name and resource type are required")
	}

	if resourceType != "desk" && resourceType != "meeting_room" {
		return nil, errors.New("invalid resource type")
	}

	var resource models.Resource
	if result := s.DB.First(&resource, id); result.Error != nil {
		return nil, errors.New("resource not found")
	}

	resource.Name = name
	resource.ResourceType = resourceType

	if result := s.DB.Save(&resource); result.Error != nil {
		return nil, errors.New("failed to update resource")
	}

	return &resource, nil
}
