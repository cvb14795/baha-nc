package main

import (
	"bahaNC/scrape"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
		Queries("baha_id", "{baha_id}", "count", "{count}").
		HandlerFunc(draw)
	http.ListenAndServe(":"+port, r)
}

func draw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	vars := mux.Vars(r)
	// 要抓幾筆資料
	count, err := strconv.Atoi(vars["count"])
	if err != nil {
		fmt.Println(err)
	}
	// 爬巴哈創作
	n := scrape.Run(vars["baha_id"], count)
	width := 800
	const LINE_SPACE = 30
	y := 0
	canvas := svg.New(w)
	canvas.Start(width, LINE_SPACE*len(n))
	for i, elem := range n {
		if i == len(n)-1 { // 最後一行預留3px避免文字被切掉
			y += LINE_SPACE - 3
		} else {
			y += LINE_SPACE
		}
		canvas.Text(0, y, elem, "font-size:15px;font-family:微軟正黑體")
	}
	canvas.End()

	fmt.Println("巴哈ID: " + vars["baha_id"])
	for i, elem := range n {
		fmt.Println("創作標題" + strconv.Itoa(i) + ": " + elem)
	}
}
