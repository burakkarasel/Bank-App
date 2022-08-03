package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain sets gin's mode to TestMode
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
