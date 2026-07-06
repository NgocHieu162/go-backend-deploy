package socket

import (
	"fmt"
	"go-backend/internal/handler"

	"github.com/gin-gonic/gin"
	server "github.com/zishang520/socket.io/servers/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

type socket struct {
	chatHandler *handler.ChatHandler
}

func NewSocket(chatHandler *handler.ChatHandler) *socket {
	return &socket{
		chatHandler: chatHandler,
	}
}

func (s *socket) Start(ginEngine *gin.Engine, AllowOrigins []string) {
	option := server.DefaultServerOptions()

	option.SetCors(&types.Cors{
		Origin: AllowOrigins,
	})
	io := server.NewServer(nil, nil)

	io.On("connection", func(args ...any) {
		socket := args[0].(*server.Socket)
		fmt.Printf("connected: %s\n", socket.Id())

		socket.On("CREATE_ROOM", func(args ...any) {
			// handler
			s.chatHandler.CreateGroup(args...)

		})

		socket.On("JOIN_ROOM", func(args ...any) {
			// handler
			s.chatHandler.JoinGroup(socket, args...)
		})

		socket.On("SEND_MESSAGE", func(args ...any) {
			// handler
			s.chatHandler.SendMessage(io, args...)
		})

		socket.On("disconnect", func(args ...any) {
			fmt.Printf("disconnected: %s\n", socket.Id())
		})
	})

	ginEngine.Any("/socket.io/", gin.WrapH(io.ServeHandler(option)))
}
