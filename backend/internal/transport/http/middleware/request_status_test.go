package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chenyme/grok2api/backend/internal/application/requeststatus"
	clientkeydomain "github.com/chenyme/grok2api/backend/internal/domain/clientkey"
	"github.com/gin-gonic/gin"
)

func TestTrackRequestStatusUsesFinalHTTPStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	registry := requeststatus.NewRegistry(time.Minute)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ClientKey, clientkeydomain.Key{ID: 11})
		c.Set(RequestIDKey, "tracked-request")
		c.Next()
	})
	router.Use(TrackRequestStatus(registry))
	router.POST("/v1/responses", func(c *gin.Context) {
		status, ok := registry.Get(11, "tracked-request", time.Now())
		if !ok || status.State != requeststatus.StateRunning {
			t.Fatalf("status during handler = %#v, ok=%v", status, ok)
		}
		c.Status(http.StatusBadGateway)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/responses", nil))
	status, ok := registry.Get(11, "tracked-request", time.Now())
	if !ok || status.State != requeststatus.StateFailed || status.StatusCode != http.StatusBadGateway {
		t.Fatalf("final status = %#v, ok=%v", status, ok)
	}
}

func TestTrackRequestStatusRecordsPanicAsFailedBeforeRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	registry := requeststatus.NewRegistry(time.Minute)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Set(ClientKey, clientkeydomain.Key{ID: 12})
		c.Set(RequestIDKey, "panicked-request")
		c.Next()
	})
	router.Use(TrackRequestStatus(registry))
	router.POST("/v1/responses", func(c *gin.Context) {
		panic("inference panic")
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/v1/responses", nil))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("response status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	status, ok := registry.Get(12, "panicked-request", time.Now())
	if !ok || status.State != requeststatus.StateFailed || status.StatusCode != http.StatusInternalServerError {
		t.Fatalf("panic status = %#v, ok=%v", status, ok)
	}
}
