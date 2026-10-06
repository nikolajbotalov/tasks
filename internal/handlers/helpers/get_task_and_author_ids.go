package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetTaskAndAuthorIDs(c *gin.Context) (id uuid.UUID, authorID uuid.UUID, err error) {
	id, err = GetTaskID(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	authorID, err = GetAuthorID(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return id, authorID, nil
}
