package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"security-class-practice/message"
	"time"
)

func logerr(err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("read timeout: ", err)
		} else if err == io.EOF {
		} else {
			log.Println("read error: ", err)
		}
		return true
	}
	return false
}

func read(conn net.Conn) {
	// client와 동일한 key
	key := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC}
	iv := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC} // initialize vector

	tmp := make([]byte, 72)

	for {
		_, err := conn.Read(tmp)
		if logerr(err) {
			break
		}

		plainText, _ := message.DesDescryption(key, iv, tmp)

		// convert bytes into Buffer
		tmpbuff := bytes.NewBuffer(plainText)
		tmpstruct := new(message.Message)

		// creates a decoder object
		gobobj := gob.NewDecoder(tmpbuff)
		// decodes buffer and unmarshals it into a Message struct (역직렬화 (byte->객체))
		gobobj.Decode(tmpstruct)

		fmt.Println(tmpstruct)
		return
	}

}

func resp(conn net.Conn) {
	msg := message.Message{ID: "Server", Data: "22.05.06"}
	bin_buf := new(bytes.Buffer)

	gobobje := gob.NewEncoder(bin_buf)
	gobobje.Encode(msg)

	conn.Write(bin_buf.Bytes())
	conn.Close()
}

func handle(conn net.Conn) {
	timeoutDuration := 2 * time.Second
	fmt.Println("Launching server...")
	conn.SetReadDeadline(time.Now().Add(timeoutDuration))

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	read(conn)
	resp(conn)
}

func main() {
	server, _ := net.Listen("tcp", "localhost:8081")
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Connection error: ", err)
			return
		}
		go handle(conn)
	}
}
