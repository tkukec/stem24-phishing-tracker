package middleware

import (
	"bytes"
	"fmt"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/runtimebag"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	assecoContext "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"io"
	"time"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"
	ContentTypeXML  = "text/xml; charset=utf-8"
	NotApplicable   = "N/A"
)

var (
	allowedContentTypes = []string{ContentTypeJSON, ContentTypeXML}
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Log(logger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		xCorrelation, _ := ctx.Get(constants.XCorrelationID)
		if xCorrelation == nil {
			xCorrelation = ""
		}

		tenantId, _ := ctx.Get(constants.TenantIdentifier)
		if tenantId == nil {
			tenantId = ""
		}

		var user *assecoContext.User
		token, _ := ctx.Get("user")
		if token != nil {
			user = assecoContext.NewUser(token.(*jwt.Token), tenantId.(string))
		}

		reqContext := assecoContext.NewRequestContext(xCorrelation.(string), tenantId.(string), user, ctx)
		log := reqContext.BuildLog(logger, "middleware.Log")

		ctx.Next()

		jsonBodyData, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			jsonBodyData = []byte{}
		}

		responseBody := NotApplicable
		if runtimebag.GetEnvBool(constants.ResponseBodyLog, false) {
			contentType := ctx.Writer.Header().Get("Content-Type")
			for _, allowedContentType := range allowedContentTypes {
				if contentType == allowedContentType {
					responseBody = w.body.String()
					break
				}
			}
		}

		elapsed := time.Since(start)
		requestTime := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)

		message := fmt.Sprintf(
			"url: %s%s\nmethod: %s\nrequestURI: %s\nrequestBodyData: %s\nresponseStatusCode: %d\nresponseBodyData: %s\nrequestTime: %s\nrequestMemoryUsage: %s",
			ctx.Request.Host, ctx.Request.URL.Path, ctx.Request.Method, ctx.Request.URL.RequestURI(), string(jsonBodyData),
			ctx.Writer.Status(), responseBody, requestTime, NotApplicable)
		log.Info().Msg(message)
	}
}
