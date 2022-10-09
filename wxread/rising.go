package wxread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type RisingWxReader struct {
	BookName      string  `json:"bookname"`
	Author        string  `json:"author"`
	ReadingCount  int     `json:"readingCount"`
	AdvicePercent float64 `json:"advicePercent"`
	Description   string  `json:"description"`
	BookImage     string  `json:"bookImage"`
}

func GetWxReadRing() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	repos := []*RisingWxReader{}

	c.OnHTML(".wr_bookList_item", func(e *colly.HTMLElement) {
		repo := &RisingWxReader{}
		repo.BookName = e.ChildText(".wr_bookList_item_title")

		repo.Author = e.ChildText(".wr_bookList_item_author")

		readingCount, err := strconv.Atoi(e.ChildText(".wr_bookList_item_reading_number"))
		if err != nil {
			log.Println("Could not get id")
		}
		repo.ReadingCount = readingCount
		advicePercentStr := e.ChildText(".wr_bookList_item_reading_percent")
		// 去掉 %
		advicePercentStr = strings.Trim(advicePercentStr, "%")
		advicePercent, _ := strconv.ParseFloat(advicePercentStr, 64)
		repo.AdvicePercent = advicePercent

		repo.Description = e.ChildText(".wr_bookList_item_desc")

		repo.BookImage = e.ChildAttr(".wr_bookList_item_cover > img", "src")

		repos = append(repos, repo)
	})

	c.Visit("https://weread.qq.com/web/category/rising")

	fmt.Printf("%d repositories\n", len(repos))
	repos_json, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("rhinofacts.json", repos_json, 0644)
	for _, repo := range repos {
		fmt.Println("Name:", repo.BookName)
		fmt.Println("Author:", repo.Author)
	}
	fmt.Println(string(repos_json))
}
