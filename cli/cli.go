package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/nomadcoders_review/explorer"
	"github.com/nomadcoders_review/rest"
)

func usage() {
	fmt.Printf("	Welcome to 노마드 코인\n\n")
	fmt.Printf("	Please use the following commands\n\n")
	fmt.Printf("	explorer: Start the HTML Explorer\n\n")
	fmt.Printf("	rest: Start the REST API (recommanded)\n\n")
	// os.Exit(0) 을 쓰면 에러코드가 안나오고
	//os.Exit(1) //을 쓰면 에러코드가 나오고
	runtime.Goexit() // 을 쓰면 모든 함수가 종료됨 그 이후는 defer만 동작
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

	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	fmt.Printf("\t변수 : %s\n", os.Args)
	if *port == 4000 {
		if strings.Compare(*mode, "rest") == 0 {
			fmt.Printf("\t(run defalut value)  ")
		}
	}
	fmt.Printf("mode = '%s', port = '%d' 로 웹서버를 구동합니다\n", *mode, *port)
	modeAndPort(port, mode)
}
