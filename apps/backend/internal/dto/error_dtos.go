package dto

import "github.com/gin-gonic/gin"

func BadRequestError(issue, details string) gin.H {
	return gin.H{
		"issue":   issue,
		"details": details,
	}
}
