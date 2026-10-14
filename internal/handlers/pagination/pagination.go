package pagination

import (
	"TaskFlow/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (page, limit int, err error) {
	pageParam := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageParam)
	if err != nil {
		return 0, 0, domain.ErrInvalidPageParam
	}
	if page < 1 {
		return 0, 0, domain.ErrInvalidPageParam
	}

	limitParam := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitParam)
	if err != nil {
		return 0, 0, domain.ErrInvalidLimitParam
	}
	if limit <= 0 || limit > 100 {
		return 0, 0, domain.ErrInvalidLimitParam
	}

	return page, limit, nil
}
