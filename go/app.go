package app

import (
	"encoding/json"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/pata", handlePata)
	http.HandleFunc("/norikae", handleNorikae)
	http.HandleFunc("/", handleRoot)
}

// このディレクトリーに入っているすべての「.html」終わるファイルをtemplateとして読み込む。
var tmpl = template.Must(template.ParseGlob("*.html"))

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// hello.htmlというtemplateを埋めて、出力する。
	tmpl.ExecuteTemplate(w, "hello.html", nil)
}

// Templateに渡す内容を分かりやすくするためのtypeを定義しておきます。
// （「Page」という名前などは重要ではありません）。
type Page struct {
	A    string
	B    string
	Pata string
}

func handlePata(w http.ResponseWriter, r *http.Request) {
	// Appengineの「Context」を通してAppengineのAPIを利用する。
	ctx := appengine.NewContext(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// templateに埋める内容をrequestのFormValueから用意する。
	content := Page{
		A: r.FormValue("a"),
		B: r.FormValue("b"),
	}

	// とりあえずPataを簡単な操作で設定しますけど、すこし工夫をすれば
	// パタトクカシーーができます。
	content.Pata = content.A + content.B

	// example.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	err := tmpl.ExecuteTemplate(w, "example.html", content)
	if err != nil {
		// もしテンプレートに問題があったらこのエラーが出ます。
		log.Errorf(ctx, "rendering template example.html failed: %v", err)
	}
}

// LineはJSONに入ってくる線路の情報をtypeとして定義している。このJSON
// にこの名前にこういうtypeのデータが入ってくるということを表している。
type Line struct {
	Name     string
	Stations []string
}

// TransitNetworkは http://fantasy-transit.appspot.com/net?format=json
// の一番外側のリストのことを表しています。
type TransitNetwork []Line

func handleNorikae(w http.ResponseWriter, r *http.Request) {
	// Appengineの「Context」を通してAppengineのAPIを利用する。
	ctx := appengine.NewContext(r)

	// clientはAppengine用のHTTPクライエントで、他のウェブページを読み込
	// むことができる。
	client := urlfetch.Client(ctx)

	// JSONとしての路線グラフ内容を読み込む
	resp, err := client.Get("http://fantasy-transit.appspot.com/net?format=json")
	if err != nil {
		panic(err)
	}

	// 読み込んだJSONをパースするJSONのDecoderを作る。
	decoder := json.NewDecoder(resp.Body)

	// JSONをパースして、「network」に保存する。
	var network TransitNetwork
	if err := decoder.Decode(&network); err != nil {
		panic(err)
	}

	// handleExampleと同じようにtemplateにテンプレートを埋めて、出力する。
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.ExecuteTemplate(w, "norikae.html", network)
	if err != nil {
		// もしテンプレートに問題があったらこのエラーが出ます。
		log.Errorf(ctx, "rendering template norikae.html failed: %v", err)
	}
}
