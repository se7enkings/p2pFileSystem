package ui

import (
	"code.google.com/p/go.net/websocket"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net/http"
)

func StartHttpServer() {
	http.Handle("/", http.FileServer(http.Dir("ui/")))
	http.Handle("/ws", websocket.Handler(socketHandler))
	err := http.ListenAndServe(settings.HttpPort, nil)
	logger.Error(err)
}
func socketHandler(ws *websocket.Conn) {
	buff := make([]byte, settings.MessageBufferSize)
	messageSize, err := ws.Read(buff)
	logger.Warning(err)
	message := string(buff[:messageSize])
    logger.Info(message)
	switch message {
	case settings.FileListRequestProtocol:
        logger.Info("requested file list by web browser")
        fileListJson := filesystem.GetFileListJson()
        logger.Info(string(fileListJson))
		ws.Write(fileListJson)

	}
}
