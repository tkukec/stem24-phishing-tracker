package middleware

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/shared"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Tenant() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if helpers.GetEnvBool(constants.SingleTenantMode, false) {
			ctx.Set(constants.TenantIdentifier, constants.AnyTenant)
			ctx.Next()
			return
		}

		tenantName := ctx.GetHeader(constants.TenantIdentifier)
		if tenantName == "" && helpers.GetEnvBool(constants.SingleTenantMode, false) {
			tenantName = constants.AnyTenant
			return
		}
		if tenantName == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Errorf("missing tenant header"))
			return
		}
		ctx.Set(constants.TenantIdentifier, tenantName)
		ctx.Next()
	}
}
