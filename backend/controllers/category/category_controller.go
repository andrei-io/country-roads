package category

import (
	"backend/roralis/domain/entity"
	"backend/roralis/domain/repo/category"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	repo category.CategoryRepo
}

func NewCategoryController(c category.CategoryRepo) CategoryController {
	return CategoryController{repo: c}
}

// Gin controller for GET /users
func (r *CategoryController) ReadAll(c *gin.Context) {

	categories, err := r.repo.GetAll()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, entity.NotFoundError)
		return
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, entity.Response{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)

}

// Gin controller for GET /users/:id
func (r *CategoryController) ReadOne(c *gin.Context) {
	id := c.Param("id")

	category, err := r.repo.Get(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, entity.NotFoundError)
		return
	}
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, entity.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
