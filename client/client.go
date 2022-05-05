package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"security-class-practice/message"
)

func send(conn net.Conn) {
	key := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC}
	iv := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC} // initialize vector

	msg := message.Message{ID: "YoungLee Choi", Data: "Hello"}
	bin_buf := new(bytes.Buffer)

	// create a encoder object
	gobobj := gob.NewEncoder(bin_buf)
	// encode buffer and marshal it into a gob object (직렬화, (객체->byte))
	gobobj.Encode(msg)

	// conn.Write(bin_buf.Bytes())
	// 암호화를 한 후, TCP 소켓(전송계층)에 내려보냄
	cryptoText, _ := message.DesEncryption(key, iv, bin_buf.Bytes())
	conn.Write(cryptoText)
	// conn.Close()
}

func recv(conn net.Conn) {
	tmp := make([]byte, 500)
	conn.Read(tmp)

	tmpbuff := bytes.NewBuffer(tmp)
	tmpstruct := new(message.Message)

	gobobjdec := gob.NewDecoder(tmpbuff)
	gobobjdec.Decode(tmpstruct)

	fmt.Println("ID: ", tmpstruct.ID)
	fmt.Println("DATA: ", tmpstruct.Data)
}

func main() {
	conn, _ := net.Dial("tcp", "localhost:8081")

	send(conn)
	recv(conn)
}
