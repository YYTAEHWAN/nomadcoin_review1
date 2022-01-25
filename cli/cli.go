package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/nomadcoders_review/explorer"
	"github.com/nomadcoders_review/rest"
)

func usage() {
	fmt.Printf("	Welcome to 노마드 코인\n\n")
	fmt.Printf("	Please use the following commands\n\n")
	fmt.Printf("	explorer: Start the HTML Explorer\n\n")
	fmt.Printf("	rest: Start the REST API (recommanded)\n\n")
	// os.Exit(0) 을 쓰면 에러코드가 안나오고
	os.Exit(1) //을 쓰면 에러코드가 나오고
}

func modeAndPort(port *int, mode *string) {
	switch *mode {
	case "rest":
		//fmt.Println(*mode, *port)
		rest.Strat(*port)
	case "html":
		//fmt.Println(*mode, *port)
		explorer.Start(*port)
	default:
		usage()
	}
}

func Start() {

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()
	fmt.Println(os.Args)
	fmt.Println(*port, *mode)
	fmt.Println("출력 완료")
	modeAndPort(port, mode)
}
