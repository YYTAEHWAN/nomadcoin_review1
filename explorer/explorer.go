package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nomadcoders_review/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template // template는 import 되고 있기 때문

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block // blockchain 패키지에 있는 Block의 포인터를 가져온거구나
}

// rw 유저에게 보내고 싶은 데이터를 적는 곳을
// rw http.ResponseWriter 라고 한다

// 요청이 파일일수도 있고 빅데이터일 수도 있기 때문에
// http.Reqeust를 사용하기 보다 포인터를 사용한다
func home(rw http.ResponseWriter, r *http.Request) {

	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.FormValue("blockData")
		blockchain.Blockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Start(ePort int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost:%d\n", ePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ePort), handler))
}
