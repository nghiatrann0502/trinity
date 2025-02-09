package ginhttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/domain"
	"github.com/nghiatrann0502/trinity/internal/ranking/core/ports"
	"github.com/nghiatrann0502/trinity/pkg/common"
	"github.com/nghiatrann0502/trinity/pkg/logger"
)

type VideoRankingUpdateBody struct {
	Action string `json:"action" example:"view"` // view, like, comment, share, watch
	Value  int    `json:"value" example:"1"`
} // @name VideoRankingUpdateBody

type GinHandler interface {
	Ping() gin.HandlerFunc
	Health() gin.HandlerFunc
	UpdateVideoScore() gin.HandlerFunc
	GetTopRanked() gin.HandlerFunc
}

type ginHandler struct {
	logger  logger.Logger
	service ports.Service
}

var _ GinHandler = (*ginHandler)(nil)

func NewGinHandler(service ports.Service, log logger.Logger) GinHandler {
	return &ginHandler{
		service: service,
		logger:  log.With(map[string]interface{}{"component": "gin_handler"}),
	}
}

func (h *ginHandler) handleError(c *gin.Context, err error) {
	var status int
	var response common.HTTPError

	if appErr, ok := err.(common.AppError); ok {
		switch appErr.Type {
		case common.NotFound:
			status = http.StatusNotFound
		case common.ValidationErr:
			status = http.StatusBadRequest
		case common.DatabaseErr:
			status = http.StatusInternalServerError
		default:
			h.logger.Error("internal server error", appErr.Err, nil)
			status = http.StatusInternalServerError
		}

		response = common.HTTPError{
			Error: appErr.Message,
			Type:  string(appErr.Type),
		}

	} else {
		h.logger.Error("internal server error", err, nil)
		status = http.StatusInternalServerError
		response = common.HTTPError{
			Error: "Internal server error",
			Type:  string(common.InternalErr),
		}
	}

	c.JSON(status, response)
}

func (h *ginHandler) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, common.SuccessResponse("pong"))
	}
}

func (h *ginHandler) Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, common.SuccessResponse("healthy"))
	}
}

// UpdateVideoScore updates the score of a video
// @Summary Update video score
// @Description Update the score of a video based on the given ID and action
// @Tags videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Param body body VideoRankingUpdateBody true "Video Ranking Update"
// @Success 201
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/videos/{id}/score [post]
func (h *ginHandler) UpdateVideoScore() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Param("id")
		id, err := strconv.Atoi(sid)
		if err != nil {
			h.logger.Error("invalid video id", err, map[string]interface{}{"id": sid})
			h.handleError(c, common.NewValidationError("invalid video id"))

			return
		}

		h.logger.Debug("video id", map[string]interface{}{"id": id})

		var body VideoRankingUpdateBody

		if err := c.BindJSON(&body); err != nil {
			h.logger.Error("invalid request body", err, map[string]interface{}{"body": c.Request.Body})
			h.handleError(c, common.NewValidationError("invalid request body"))
			return
		}

		dto := domain.VideoRankingUpdate{
			VideoID: id,
			Action:  body.Action,
			Value:   body.Value,
		}

		if err := h.service.UpdateRanking(c, dto); err != nil {
			h.logger.Error("cannot update video score", err, map[string]interface{}{"video_id": id})
			h.handleError(c, err)
			return
		}

		c.Status(http.StatusCreated)
	}
}

// GetTopRanked get top ranked videos
// @Summary Get top ranked videos
// @Description Get top ranked videos
// @Tags videos
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit number"
// @Success 200 {object} Response
// @Router /v1/videos/ranked [get]
func (h *ginHandler) GetTopRanked() gin.HandlerFunc {
	return func(c *gin.Context) {
		paging := common.Paging{}
		if err := c.ShouldBindQuery(&paging); err != nil {
			h.logger.Error("invalid request query", err, map[string]interface{}{"query": c.Request.URL.RawQuery})
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request query",
			})
		}

		h.logger.Debug("paging", map[string]interface{}{"page": paging.Page, "limit": paging.Limit})

		videos, err := h.service.GetTopRanked(c.Request.Context(), paging.Page, paging.Limit)
		if err != nil {
			h.logger.Error("cannot get top ranked", err, nil)
			h.handleError(c, err)
			return
		}

		c.JSON(200, common.SuccessResponseWithPaging(videos, paging))
	}
}
