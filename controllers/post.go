package controllers

import (
	"net/http"

	"apipress/database"
	"apipress/models"
	"apipress/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostInput struct {
	Title       string `json:"title" binding:"required"`
	ContentMD   string `json:"content_md" binding:"required"`
	ContentHTML string `json:"content_html"`
	Published   bool   `json:"published"`
}

type UpdatePostInput struct {
	Title       string `json:"title"`
	ContentMD   string `json:"content_md"`
	ContentHTML string `json:"content_html"`
	Published   *bool  `json:"published"`
}

func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	// Check if slug exists
	slugExists := func(s string) bool {
		var count int64
		database.DB.Model(&models.Post{}).Where("slug = ?", s).Count(&count)
		return count > 0
	}

	slug := utils.GenerateUniqueSlug(input.Title, slugExists)

	post := models.Post{
		Title:       input.Title,
		Slug:        slug,
		ContentMD:   input.ContentMD,
		ContentHTML: input.ContentHTML,
		Published:   input.Published,
		AuthorID:    userID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	// Only return published posts for public view, unless implemented otherwise.
	// For this extensive endpoint, let's assume public timeline
	if err := database.DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Where("published = ?", true).Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post

	if err := database.DB.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Where("slug = ?", slug).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Ownership check
	if post.AuthorID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to edit this post"})
		return
	}

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update slug if title changed significantly? Usually slugs are permanent, but let's allow basic updates
	// For now, let's keep slug persistent unless logic demands otherwise.

	if input.Title != "" {
		post.Title = input.Title
	}
	if input.ContentMD != "" {
		post.ContentMD = input.ContentMD
	}
	if input.ContentHTML != "" {
		post.ContentHTML = input.ContentHTML
	}
	if input.Published != nil {
		post.Published = *input.Published
	}

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Ownership check
	if post.AuthorID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
