package ui

import (
	"code.google.com/p/go.net/websocket"
	"github.com/CRVV/p2pFileSystem/filesystem"
	"github.com/CRVV/p2pFileSystem/logger"
	"github.com/CRVV/p2pFileSystem/settings"
	"net/http"
    "strings"
    "github.com/CRVV/p2pFileSystem/data"
    "bytes"
    "time"
    "os/exec"
    "runtime"
)

func StartHttpServer() {
	http.Handle("/", httpHandler{})
	http.Handle("/ws", websocket.Handler(socketHandler))
    go func() {
        err := http.ListenAndServe(settings.HttpPort, nil)
        logger.Error(err)
    }()
    startBrowser("http://localhost" + settings.HttpPort)
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
		noticeChan := make(chan int)
		//TODO: this will go a goroutine every time
		go fileListChangeListener(ws, noticeChan)
		ws.Write(filesystem.GetFileListJson(noticeChan))
	}
	a := make(chan int)
	<-a
}
func fileListChangeListener(ws *websocket.Conn, noticeChan chan int) {
	for {
		<-noticeChan
		ws.Write(filesystem.GetFileListJson(noticeChan))
	}
}

type httpHandler struct {}
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    for strings.HasPrefix(path, "/") {
        path = path[1:]
    }
    if path == "" {
        path = "index.html"
    }
    logger.Info(path)
    data, err := data.Asset(path)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    http.ServeContent(w, r, path, time.Now(), bytes.NewReader(data))
}

func startBrowser(url string) bool {
    var args []string
    switch runtime.GOOS {
        case "darwin":
        args = []string{"open"}
        case "windows":
        args = []string{"cmd", "/c", "start"}
        default:
        args = []string{"xdg-open"}
    }
    cmd := exec.Command(args[0], append(args[1:], url)...)
    return cmd.Start() == nil
}
