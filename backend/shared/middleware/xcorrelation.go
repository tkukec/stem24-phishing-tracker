package middleware

import (
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func XCorrelate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		xCorrelationValue := ctx.GetHeader(constants.XCorrelationID)
		if xCorrelationValue == "" {
			xCorrelationValue = uuid.NewString()
		}
		ctx.Set(constants.XCorrelationID, xCorrelationValue)
		ctx.Next()
	}
}
