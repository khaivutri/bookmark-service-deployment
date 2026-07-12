package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler/v1/dto"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/internal/service"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// ShortenURL defines the interface for handling URL shortening operations.
type ShortenURL interface {
	CreateShortenLink(ctx *gin.Context) 
	Redirect( ctx *gin.Context)
}

// shortenURL implements the ShortenURL interface.
type shortenURL struct {
	service service.ShortenURL
}


// NewShortenURL creates and returns a new instance of ShortenURL handler.
func NewShortenURL(service service.ShortenURL) ShortenURL {
	return &shortenURL{service: service}
}


// CreateShortenLink generates a shortened code for the provided URL.
//
// @Summary      Create Shorten Link
// @Description  Generates a short code for the provided URL.
// @Tags         ShortenURL
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ShortenURLReq  true  "Shorten URL request"
// @Success      200      {object}  dto.ShortenURLRes
// @Failure      400      {object}  response.Message
// @Failure      500      {object}  response.Message
// @Router       /v1/links/shorten [post]
func (s *shortenURL) CreateShortenLink(ctx *gin.Context) {
	var req dto.ShortenURLReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	code, err := s.service.CreateCodeFromLink(ctx, req.URL, req.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenURL.CreateShortenLink").Msg("failed to code from link")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	
	response := dto.ShortenURLRes{	Code: 		code,
									Message: 	"Shorten URL generated successfully!",
								}
	ctx.JSON(http.StatusOK, response)
}


// Redirect retrieves the original URL associated with the given code
// and redirects the client to that URL.
//
// @Summary      Redirect to original URL
// @Description  Redirects the client to the original URL using the provided short code.
// @Tags         ShortenURL
// @Accept       json
// @Produce      json
// @Param        code  path      string  true  "Short code"
// @Success      302   "Redirect to original URL"
// @Failure      400   {object}  response.Message
// @Failure      404   {object}  response.Message
// @Failure      500   {object}  response.Message
// @Router       /v1/links/redirect/{code} [get]
func (s *shortenURL) Redirect( ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := s.service.GetLinkFromCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrCodeNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, response.CodeNotFoundResponse)
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenURL.Redirect").Msg("failed to get link from code")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	
	ctx.Redirect(http.StatusFound, url)
}