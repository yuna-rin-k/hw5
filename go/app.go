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

	// Goだとstringの[x]で「x番目の値」を取ると、「x番目の文字」ではな
	// くて、「x番目のbyte」を取るので、日本語の文字場合に１文字が複数
	// のbyteで表しているので、「x番目のbyte」は特に意味がありません。
	// 「x番目の文字」を取るためにまずstringから[]rune （runeのスライ
	// ス）に変換した方がいいです。ちなみに「rune」は英語では「文字」
	// みたいな意味です。
	//
	// （参考のために：もう一つの文字１個ずつ取る方法として、
	// strings.split(content.A, "") でstringから[]stringに変換すること
	// が出来ます。）
	aMoji := []rune(content.A)
	bMoji := []rune(content.B)
	// とりあえずPataを簡単な操作で設定しますけど、すこし工夫をすれば
	// パタトクカシーーができます。
	str := ""
	for i:=0; i<len(aMoji); i++ {
		str := str + aMoji[i]
		str := str + bMoji[i];
	}
	pata := append(aMoji)
	//pata := append(aMoji, bMoji...)
	// []runeからstringに戻して、テンプレートで使うcontent.Pataの変数
	// に入れておきます。
	content.Pata = string(pata)

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
