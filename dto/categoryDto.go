package dto

import "movie/models"

type CategoryDto struct {
	Name string `form:"name" `
}

type CategoryResult struct {
	Category   []models.Category
	MovieCount int
}
