package pprof

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func Run(port uint16) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicln(err)
	}
	arr := strings.Split(listen.Addr().String(), ":")
	p := arr[len(arr)-1]
	fmt.Printf("pprof http://localhost:%s/debug/pprof/\n", p)
	log.Println(http.Serve(listen, nil))
}
