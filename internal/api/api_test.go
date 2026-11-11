package api

import (
	"net"
	"testing"

	redisPkg "bookmark-management/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Start blocks until the server stops, so it is tested with a port that is already taken:
// the server cannot bind it and Start returns the error right away.
func TestStart_PortInUse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	app := NewEngine(&Config{AppPort: port}, redisPkg.InitMockRedis(t))

	assert.Error(t, app.Start())
}
