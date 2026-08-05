package middleware

import (
	"net/http"
	"time"

	"github.com/chenyme/grok2api/backend/internal/application/requeststatus"
	clientkeydomain "github.com/chenyme/grok2api/backend/internal/domain/clientkey"
	"github.com/gin-gonic/gin"
)

var trackedRequestPaths = map[string]struct{}{
	"/v1/responses":         {},
	"/v1/responses/compact": {},
	"/v1/chat/completions":  {},
	"/v1/messages":          {},
}

// TrackRequestStatus records the lifecycle of synchronous text inference requests.
func TrackRequestStatus(registry *requeststatus.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}
		if _, tracked := trackedRequestPaths[c.Request.URL.Path]; !tracked {
			c.Next()
			return
		}
		clientValue, clientExists := c.Get(ClientKey)
		clientKey, clientOK := clientValue.(clientkeydomain.Key)
		requestValue, requestExists := c.Get(RequestIDKey)
		requestID, requestOK := requestValue.(string)
		if !clientExists || !clientOK || !requestExists || !requestOK {
			c.Next()
			return
		}
		registry.Start(clientKey.ID, requestID, time.Now())
		defer func() {
			if recovered := recover(); recovered != nil {
				registry.Finish(clientKey.ID, requestID, http.StatusInternalServerError, time.Now())
				panic(recovered)
			}
			registry.Finish(clientKey.ID, requestID, c.Writer.Status(), time.Now())
		}()
		c.Next()
	}
}
