package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"path"
//	"time"
)

var (
	commu_port string = ":8237"
	data_port	string = ":8238"
)


func handleFile(file string, c net.Conn, reply chan string) {

	log.Printf("save file %s\n", file)
	var cnt int

	dir,_ := path.Split(file)
	if os.MkdirAll(dir, 644) != nil {
		log.Printf("data error: create dir %s\n", dir)
		reply<- "data error: create dir " + dir
		return
	}
	
	f, err := os.Create(file)
	if err != nil {
		reply<- "data file create error"
		return
	} else {
		reply<- "data ready"
	}
	
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("finised, error = %s\n", err)
			f.Close()
			break;
		} else {
			cnt += n
			buf = buf[:n]
			f.Write(buf)
		}
	}
	
	log.Printf("get file byte %d\n", cnt)	
	reply<- "data ok " + strconv.Itoa(cnt)
}

func handleCmd(ch chan string, l *net.TCPListener, reply chan string) {
	for {
	
		file := <-ch
		log.Printf("handCmd %s\n", file)
		
		reply <- "file ok"
		
		dc, err := l.Accept()
		if err != nil {
			log.Fatalf("net accept error: %s\n", err)
		}
		
		log.Printf("get data connection\n")
		
		go handleFile(file, dc, reply)
		
	}
}


func replyCmd(c net.Conn, reply chan string) {
	for { 
		c.Write([]byte(<-reply))
	}
}

func getCmd(c net.Conn, ch chan string) {
	for { 
		buf := make([]byte, 1024)
		_, err := c.Read(buf)
		switch err {

			case nil:
				log.Printf("get cmd: %s\n", string(buf))
				ch <- string(buf)
			case os.EOF:
				log.Printf("client closed\n")
				return				
			default:
				//the client seems not close correctly
				log.Printf("connet error: %s\n", err)
				return
		}
	}
	
}

func main() {

	log.Printf("send file server begin...\n")

	caddr, err := net.ResolveTCPAddr("tcp", commu_port)
	if err != nil {
		log.Fatalf("address resolve error: %s\n", err)
	}

	l, err := net.ListenTCP("tcp", caddr)
	if err != nil {
		log.Fatalf("address listen error: %s\n", err)
	}

	daddr, err := net.ResolveTCPAddr("tcp", data_port)
	if err != nil {
		log.Fatalf("address resolve error: %s\n", err);
	}			

	dl, err := net.ListenTCP("tcp", daddr)
	if err != nil {
		log.Fatalf("address listen error: %s\n", err)
	}

				
	for {

		cconn, err := l.Accept()
		if err != nil {
			log.Fatalf("net accept error: %s\n", err)
		}
		log.Printf("get connect from %s\n", cconn.RemoteAddr())
		
		ch := make(chan string)
		
		var reply chan string = make(chan string, 10)
		
		go handleCmd(ch, dl, reply)
					
		go replyCmd(cconn, reply)
		
		getCmd(cconn, ch)
	}
}
	