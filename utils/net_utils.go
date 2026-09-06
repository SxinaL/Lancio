package utils
import (
	"fmt"
	"net"
)
// pickFreePort 在本地寻找一个空闲端口，用于浮窗子进程独立的 Wails dev server
func pickFreePort() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return ""
	}
	defer ln.Close()
	return fmt.Sprintf("localhost:%d", ln.Addr().(*net.TCPAddr).Port)
}

