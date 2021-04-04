package service

import (
	"io"

	"github.com/gorilla/websocket"
)

var secretKey string

func SetSecretKey(sec string) {
	if sec == "" {
		secretKey = Md5Str("default secret")
	} else {
		secretKey = Md5Str(sec)
	}
}

func WsTcpcopyWorker(dst io.Writer, src *websocket.Conn, doneCh chan<- bool) {
	for {
		_, b, err := src.ReadMessage()
		if err != nil {
			// log.Panicln(err)
			break
		}

		_, err = dst.Write(AesDecrypt(b, secretKey))
		if err != nil {
			// log.Panicln(err)
			break
		}
	}
	doneCh <- true
}

func TcpWscopyWorker(dst *websocket.Conn, src io.Reader, doneCh chan<- bool) {
	p := make([]byte, 512)
	for {
		n, err := src.Read(p)
		if err != nil {
			// log.Panicln(err)
			break
		}
		err = dst.WriteMessage(websocket.BinaryMessage,
			AesEncrypt(p[:n], secretKey))
		if err != nil {
			// log.Panicln(err)
			break
		}
	}
	doneCh <- true
}
