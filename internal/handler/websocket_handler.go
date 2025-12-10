package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/omidnikrah/duckparty-backend/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	broadcaster *ws.SocketBroadcaster
}

func NewWebSocketHandler(broadcaster *ws.SocketBroadcaster) *WebSocketHandler {
	return &WebSocketHandler{
		broadcaster: broadcaster,
	}
}

// HandleWebSocket handles websocket requests from clients
// @Summary      WebSocket connection for real-time duck notifications
// @Description  Establishes a WebSocket connection to receive real-time notifications when new ducks are added
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Success      101  "Switching Protocols"
// @Failure      400  {object}  map[string]string  "Error message"
// @Router       /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade WebSocket connection", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	h.broadcaster.Add(conn)
}
