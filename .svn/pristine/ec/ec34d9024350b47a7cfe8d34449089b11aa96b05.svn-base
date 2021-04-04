package service

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

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

func FormatURL(urlStr string) url.URL {
	// url
	var u, porotocal string
	if strings.HasPrefix(urlStr, "http://") {
		u = strings.Replace(urlStr, "http://", "", 1)
		porotocal = "ws"
	} else if strings.HasPrefix(urlStr, "https://") {
		u = strings.Replace(urlStr, "https://", "", 1)
		porotocal = "wss"
	} else if strings.HasPrefix(urlStr, "ws://") {
		u = strings.Replace(urlStr, "ws://", "", 1)
		porotocal = "ws"
	} else if strings.HasPrefix(urlStr, "wss://") {
		u = strings.Replace(urlStr, "wss://", "", 1)
		porotocal = "wss"
	} else {
		porotocal = "ws"
		u = urlStr
	}
	q, err := url.Parse(fmt.Sprintf("%s%s", porotocal+"://", u))
	if err != nil {
		log.Println(u, err.Error())
		return url.URL{}
	}
	if !strings.HasSuffix(q.Path, "/") {
		q.Path = q.Path + "/"
	}
	t := url.URL{Scheme: porotocal,
		Host: q.Host,
		Path: q.Path,
	} //RawQuery: "info=" + urlPara
	return t

}
