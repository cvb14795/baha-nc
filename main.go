package main

import (
	"bahaNC/scrape"
	"fmt"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/gorilla/mux"
)

func main() {
	port := "80"
	if v := os.Getenv("PORT"); len(v) > 0 {
		port = v
	}
	r := mux.NewRouter()
	r.Path("/get").
		//接收?後參數
		Queries("baha_id", "{baha_id}").
		HandlerFunc(draw)

	http.ListenAndServe(":"+port, r)
}

func draw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	vars := mux.Vars(r)
	// 預設為自己ID
	baha_id := "cvb14795"
	if len(vars) > 0 {
		baha_id = vars["baha_id"]
	}

	// 爬巴哈創作
	n, t := scrape.Run(baha_id)

	width := 800
	height := 50
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Text(0, height/2, "最新創作： "+n, "font-size:15px;font-family:微軟正黑體")
	canvas.Text(0, height-5, "發佈時間： "+t, "font-size:15px;font-family:微軟正黑體")
	canvas.End()

	fmt.Println("巴哈ID: " + baha_id)
	fmt.Println("最新創作: " + n)
	fmt.Println("發佈時間: " + t)

}
