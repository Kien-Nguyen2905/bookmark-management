package handler

import (
	"net/http"
	"time"

	"bookmark-management/internal/service"

	"github.com/gin-gonic/gin"
)

// ShortenLink is the HTTP handler of the link shortening endpoint.
type ShortenLink interface {
	CreateShortenLink(c *gin.Context)
}

type shortenLink struct {
	svc service.ShortenUrl
}

func NewShortenLink(svc service.ShortenUrl) ShortenLink {
	return &shortenLink{svc: svc}
}

type shortenLinkInput struct {
	// Exp is the lifetime of the link in seconds, at most 604800 (7 days).
	Exp int    `json:"exp" binding:"required,gt=0,lte=604800" example:"604800"`
	URL string `json:"url" binding:"required,url" example:"https://example.com"`
}

// CreateShortenLink godoc
// @Summary Shorten a link
// @Description Generate a 7-character code for the url and store it in Redis for exp seconds
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body shortenLinkInput true "Input required"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/links/shorten [post]
func (h *shortenLink) CreateShortenLink(c *gin.Context) {
	input := &shortenLinkInput{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	exp := time.Duration(input.Exp) * time.Second
	code, err := h.svc.CreateShortenLink(c, input.URL, exp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code, "message": "Shorten URL generated successfully!"})
}
