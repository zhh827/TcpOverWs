package service

import (
	// _ "embed" // need go 1.16+
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed static/index.html
// var htmlfile string

var veryStr string = "hello word!"

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 6 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = pongWait - 2
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

func wsrelayHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	info := values.Get("info")
	if info == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, veryStr)
		return
	}
	r.URL.RawQuery = AesDecryptBase64(info, Md5Str(secretKey))
	// log.Println(AesDecryptBase64(info, Md5Str(secretKey)))
	passwd := r.URL.Query().Get("token")
	tcptarget := r.URL.Query().Get("tcptarget")
	if passwd == token {
		log.Printf("[INFO] login sucess: %s", passwd)
	} else {
		log.Printf("[INFO] login failed: %s", passwd)
		return
	}

	conn, err := net.Dial("tcp", tcptarget)
	log.Printf("[INFO] %s From %s Start New Connect to %s ", passwd, r.Host, tcptarget)
	if err != nil {
		log.Printf("[ERROR] %v \n", err)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] %v \n", err)
		return
	}
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	doneCh := make(chan bool)
	go WsTcpcopyWorker(conn, ws, doneCh) // ws -> tcp
	go TcpWscopyWorker(ws, conn, doneCh) // tcp -> ws
	<-doneCh
	conn.Close()
	ws.Close()
	log.Printf("[INFO] %s From %s To Connect to %s  Close!", passwd, ws.RemoteAddr().String(), tcptarget)

}

var token string

func Server(wsStr, passwd string) {
	var certFile string
	var keyFile string
	token = passwd
	log.Printf("[INFO] Websocket Listening on %s\n", wsStr)
	log.Printf("[INFO] Access password: %s\n", passwd)
	http.HandleFunc("/", wsrelayHandler)

	var err error
	if certFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(wsStr, certFile, keyFile, nil)
	} else {
		err = http.ListenAndServe(wsStr, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
