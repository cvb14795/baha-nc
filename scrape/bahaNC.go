package scrape

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func Run(baha_id string) (string, string) {
	const TIME_LAYOUT = "2006-01-02 15:04:05"        // 時間格式化常量
	t := []time.Time{}                               // 發佈時間
	t_max, _ := time.Parse(TIME_LAYOUT, TIME_LAYOUT) // 最新時間 初始化為最小
	i_max := 0                                       // 最新時間index
	n := []string{}                                  // 發佈標題
	c := colly.NewCollector()                        // 在colly中使用 Collector 這類物件 來做事情

	c.OnHTML(".ST1", func(e *colly.HTMLElement) { // 找發佈時間
		// 巴哈的分隔符是│不是|
		tmp := strings.Split(e.Text, "│")
		// 有找到(1為自己 大於1才有找到)
		if len(tmp) > 1 {
			// 取出字串的時間部份 轉time格式
			t_f, _ := time.Parse(TIME_LAYOUT, tmp[1])
			t = append(t, t_f)
			// 若這次時間較最大值晚 則更新
			if t_f.After(t_max) {
				t_max = t_f
				i_max = len(t) - 1
			}
			// 列出該時間
			/* fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n",
			t_f.Year(), t_f.Month(), t_f.Day(),
			t_f.Hour(), t_f.Minute(), t_f.Second()) */
		}
	})

	c.OnHTML(".TS1", func(e *colly.HTMLElement) { // 找發佈標題
		n = append(n, e.Text)
	})

	c.OnRequest(func(r *colly.Request) { // iT邦幫忙需要寫這一段 User-Agent才給爬
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	})

	rep := strings.NewReplacer("{baha_id}", baha_id)
	c.Visit(rep.Replace("https://home.gamer.com.tw/creation.php?owner={baha_id}"))

	/* fmt.Printf("最新時間: %d-%02d-%02d %02d:%02d:%02d\n",
		t_max.Year(), t_max.Month(), t_max.Day(),
		t_max.Hour(), t_max.Minute(), t_max.Second())
	fmt.Println("最新標題: ", n[i_max]) */

	return n[i_max], t_max.Format(TIME_LAYOUT)
}
