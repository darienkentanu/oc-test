package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	r := Response{}

	res, err := http.Get("https://gist.githubusercontent.com/nubors/eecf5b8dc838d4e6cc9de9f7b5db236f/raw/d34e1823906d3ab36ccc2e687fcafedf3eacfac9/jne-awb.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	temp := make([]string, 0)

	doc.Find("table.table_style.tracking").Each(func(i int, sel *goquery.Selection) {
		if i > 1 {
			sel.Find("td").Each(func(i int, sel *goquery.Selection) {
				if i > 0 {
					temp = append(temp, sel.Text())
				}
			})
		}
	})

	var hs = []Histories{}

	var h = Histories{}
	for i := 0; i < len(temp); i += 2 {
		for j := i + 1; i < len(temp); j += 2 {
			h.Description = temp[j]
			h.CreatedAt = parsingTimeFromString(temp[i])
			h.Formatted.CreatedAt = timeFromTime(h.CreatedAt)
			hs = append(hs, h)
			break
		}
	}

	r.Status.Code = "060101"
	r.Status.Message = "Delivery tracking detail fetched successfully"
	r.Data.Histories = hs
	receiv := hs[len(hs)-1]
	received := receiv.Description
	receivedSLC := strings.SplitN(received, "DELIVERED TO [", 2)
	received2SLC := strings.SplitN(receivedSLC[1], "  |", 2)
	r.Data.ReceivedBy = received2SLC[0]

	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

}

type Response struct {
	Status struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Data Data `json:"data"`
}

type Data struct {
	ReceivedBy string      `json:"receivedBy"`
	Histories  []Histories `json:"histories"`
}

type Histories struct {
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Formatted   struct {
		CreatedAt string `json:"createdAt"`
	} `json:"formatted"`
}

func timeFromTime(t time.Time) string {
	var date = t.Format("02 January 2006, 15:04 MST")
	return date
}

func parsingTimeFromString(val string) time.Time {
	layoutFormat := "02-01-2006T15:04 MST"
	val = strings.Replace(val, " ", "T", -1)
	val = fmt.Sprint(val, " WIB")
	date, err := time.Parse(layoutFormat, val)
	if err != nil {
		fmt.Println(err)
	}
	return date
}
