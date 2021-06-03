package scrape

import (
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Result struct {
	Time time.Time
	Name string
}

func Run(baha_id string, count int) []string {
	const TIME_LAYOUT = "2006-01-02 15:04:05" // 時間格式化常量
	result := []Result{}                      // 將時間與標題組合以排序
	t := []time.Time{}                        // 發佈時間(time)
	n := []string{}                           // 發佈標題
	n_sort := []string{}                      // 發佈標題(排序後)
	c := colly.NewCollector()

	c.OnHTML(".ST1", func(e *colly.HTMLElement) { // 找發佈時間
		// 巴哈的分隔符是│不是|
		tmp := strings.Split(e.Text, "│")
		// 有找到(1為自己 大於1才有找到)
		if len(tmp) > 1 {
			// 取出字串的時間部份 轉time格式
			t_f, _ := time.Parse(TIME_LAYOUT, tmp[1])
			t = append(t, t_f)
		}
	})

	c.OnHTML(".TS1", func(e *colly.HTMLElement) { // 找發佈標題
		n = append(n, e.Text)
	})

	c.OnRequest(func(r *colly.Request) { // 需要寫這一段 User-Agent才給爬
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	})

	rep := strings.NewReplacer("{baha_id}", baha_id)
	c.Visit(rep.Replace("https://home.gamer.com.tw/creation.php?owner={baha_id}"))

	for i, elem := range t {
		result = append(result, Result{Time: elem, Name: n[i]})
	}
	sort.Slice(result, func(i, j int) bool {
		//由時間新到舊排序
		return result[i].Time.After(result[j].Time)
	})
	for i, elem := range result {
		if i == count {
			break
		}
		n_sort = append(n_sort, elem.Name)
	}
	return n_sort
}
