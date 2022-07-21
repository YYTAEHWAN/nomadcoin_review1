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

type testData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// rw 유저에게 보내고 싶은 데이터를 적는 곳을
// rw http.ResponseWriter 라고 한다

// 요청이 파일일수도 있고 빅데이터일 수도 있기 때문에
// http.Reqeust를 사용하기 보다 포인터를 사용한다
func home(rw http.ResponseWriter, r *http.Request) {

	data := homeData{"Home", nil}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		//r.ParseForm()
		//data := r.FormValue("blockData") 아마 Tx로 대체하겠죠
		blockchain.AddBlock(blockchain.Blockchain())
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func test(rw http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("explorer/templates/pages/test.gohtml"))
	data := testData{"This is stressoffcoin Page", blockchain.Blocks(blockchain.Blockchain())} // templates에 전달할 데이터를 만들어준 뒤
	tmpl.Execute(rw, data)

}

func Start(ePort int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	handler.HandleFunc("/test", test)
	fmt.Printf("Listening on http://localhost:%d\n", ePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ePort), handler))
}
