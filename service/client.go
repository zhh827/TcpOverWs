package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func tcprelayHandler(tcpconn net.Conn, wsStr, tcptarget, passwd string) {
	defer tcpconn.Close()
	// parameter
	para := fmt.Sprintf("token=%s&tcptarget=%s", passwd, tcptarget)
	urlPara := url.QueryEscape(AesEncryptBase64(para, Md5Str(secretKey)))
	c := FormatURL(wsStr)
	c.RawQuery = "info=" + urlPara
	log.Printf("[INFO] Connecting to %s", c.String())

	// websocket
	wsDial := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 15 * time.Second,
	}
	ws, _, err := wsDial.Dial(c.String(), nil)
	if err != nil {
		log.Printf("[ERROR] Server connect target %s failed\n", tcptarget)
		return
	}
	pingTicker := time.NewTicker(pingPeriod)
	// websocket ping
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

	defer pingTicker.Stop()

	doneCh := make(chan bool)
	go WsTcpcopyWorker(tcpconn, ws, doneCh) // ws -> tcp
	go TcpWscopyWorker(ws, tcpconn, doneCh) // tcp -> ws
	<-doneCh
	log.Printf("[INFO] Local: %s , Remote: %s Close!", tcpconn.LocalAddr().String(), tcpconn.RemoteAddr().String())
}

func Client(wsStr, tcplisten, tcptarget, passwd string) {
	// test websocket connect
	c := FormatURL(wsStr)
	var turl string
	if c.Scheme == "ws" {
		turl = "http://" + c.Host + c.Path
	} else if c.Scheme == "wss" {
		turl = "https://" + c.Host + c.Path
	}
	res, err := http.Get(turl)
	if err != nil {
		log.Printf("[ERROR] Test connect to %s error: %s\n", c.String(), err)
		return
	}
	log.Printf("[INFO] Test GET %s OK!\n", turl)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("get data err", err)
		return
	}
	if string(data) != veryStr {
		log.Printf("[ERROR] Verify connect to %s error:%s\n", c.User.String(), string(data))
	}

	// tcp listen server
	listen, err := net.Listen("tcp", tcplisten)
	if err != nil {
		log.Println("[INFO] Listen failed:", err)
		return
	}
	for {
		conn, err := listen.Accept()
		log.Printf("[INFO] New connect Local: %s , Remote: %s ", conn.LocalAddr().String(), conn.RemoteAddr().String())
		if err != nil {
			log.Println("[INFO] Accept failed, err:", err)
			continue
		}
		go tcprelayHandler(conn, wsStr, tcptarget, passwd)
	}
}
