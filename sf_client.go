package main

import (
	"log"
	"net"
	"os"
	//"time"
	"flag"
)

var (
	commu_port string = ":8237"
	data_port	string = ":8238"
	server_ip	string = "127.0.0.1"
)

func init() {
	flag.StringVar(&server_ip, "s", "127.0.0.1", "the server ip")
}

func main() {
	flag.Parse()
	
	nfiles := flag.NArg()
	if nfiles < 1 {
		log.Fatalf("Error: need arguements for files\n")
	}
		
	log.Printf("send file client begin...\n")
	
	addr := server_ip + commu_port
	caddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("address resolve error: %s\n", err)
	}
	
	cc, err := net.DialTCP("tcp", nil, caddr)
	if err != nil {
		log.Fatalf("connect to address error: %s\n", err)
	}
	
	addr = server_ip + data_port
	daddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("address resolve error: %s\n", err)
	}
	
	dest_dir := ""
	file := flag.Arg(nfiles-1)
	if fi, err := os.Lstat(file); !(err == nil && fi.IsRegular()) {
		//the last arg is the dest dir
		nfiles -= 1
		dest_dir = file + "/"
	}
	
	log.Printf("nfiles = %d\n", nfiles)	
	for i:=0; i<nfiles; i++ {
		file := flag.Arg(i)
		log.Printf("arg[i] = %s\n", file)
				
		//open file
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("open file %s error: %s.\n", file, err)
		}
		
		//write cmd
		log.Printf("dest path: %s\n", dest_dir + file)
		cc.Write([]byte(dest_dir + file))
	
		//get reply
		b := make([]byte, 256)
		cc.Read(b)
		log.Printf("get cmd reply: %s\n", string(b))
		
		//check reply here
		
		dc, err := net.DialTCP("tcp", nil, daddr)
		if err != nil {
			log.Fatalf("connect to address error: %s\n", err)
		}	

		//data ready
		cc.Read(b)
		log.Printf("get cmd reply: %s\n", string(b)) //data ready?
		//check data ready here
		
		var cnt int
		data := make([]byte, 1024)
		for {
			n, err := f.Read(data)
			if err == nil {
				//log.Printf("write out %d bytes\n", n)
				data = data[:n]
				dc.Write(data)
				cnt += n
			} else {
				log.Printf("finished, status:%s\n", err)
				f.Close()
				dc.Close()
				break
			}
		}
	
		log.Printf("write out %d bytes file\n", cnt)
		
		cc.Read(b)
		log.Printf("get cmd reply: %s\n", string(b))
	
	}
	
	cc.Close()
}