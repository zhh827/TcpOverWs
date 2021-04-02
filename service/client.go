package service

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func tcprelayHandler(conn net.Conn, wsStr, tcptarget, passwd string) {
	para := fmt.Sprintf("token=%s&tcptarget=%s", passwd, tcptarget)

	u := url.URL{Scheme: "ws",
		Host:     wsStr,
		Path:     "/",
		RawQuery: "info=" + AesEncryptBase64(para, Md5Str(secretKey))}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()

	go func() {
		for {
			<-pingTicker.C
			ws.SetWriteDeadline(time.Now().Add(pongWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				// log.Println(err)
				return
			}
			// log.Println("ping is ok!")
		}
	}()

	if err != nil {
		log.Fatal(err)
	}
	doneCh := make(chan bool)
	go WsTcpcopyWorker(conn, ws, doneCh) // ws -> tcp
	go TcpWscopyWorker(ws, conn, doneCh) // tcp -> ws
	<-doneCh
	conn.Close()
	ws.Close()

	log.Printf("Local: %s , Remote: %s Close!", conn.LocalAddr().String(), conn.RemoteAddr().String())
}

func Client(wsStr, tcplisten, tcptarget, passwd string) {
	listen, err := net.Listen("tcp", tcplisten)
	if err != nil {
		log.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		log.Printf("New connect Local: %s , Remote: %s ", conn.LocalAddr().String(), conn.RemoteAddr().String())
		if err != nil {
			log.Println("accept failed, err:", err)
			continue
		}
		go tcprelayHandler(conn, wsStr, tcptarget, passwd)
	}
}
