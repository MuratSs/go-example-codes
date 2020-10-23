package main
import (
	"fmt"
	"io"
	"log"
	"net"
	//"net/http"
)

var (
	counter int
	// TODO configgurable
	listenAddr = "localhost:3003"

	// TODO configurable
	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)
func main() {
   listener, err := net.Listen("tcp" ,listenAddr)
   if err != nil {
   	log.Fatal("failed listener", err)
   }
   defer listener.Close()
   for {
	   conn, err := listener.Accept()
	   if err != nil{
		   log.Printf("failed to accept onnection:%s",err)
	   }

	   backend := choseBackend()
	   fmt.Println("counter=%d backend=%s\n",counter,backend)
	   go func() {
		   err := proxy(backend,conn)
		   if err != nil {
              log.Printf("warning: proxy failed: %v", err)
		   }
	   }()
   }

}

func proxy(backend string, c net.Conn) error {
	bc, err := net.Dial("tcp",backend)
	if err != nil {
		return fmt.Errorf("failed to connect to backend %s: %v", backend, err)
	}

	// c -> bc
	go io.Copy(bc,c)
	// bc -> c
	go io.Copy(c ,bc)
    return nil
}
func choseBackend() string{
	//TODO chose randomly
	s := server[counter%len(server)]
	counter++
	return s
}