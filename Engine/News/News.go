package News

import (
	config "../../Config"
	s "../../Struct"
	u "../../Utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var urlPath = "/news"

func Scrap() ([]s.News, error) {
	var (
		Newses []s.News
		R      []s.News
		E      error
	)

	c := colly.NewCollector(
	//colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Minute * 2,
			KeepAlive: time.Minute * 2,
		}).DialContext,
		IdleConnTimeout:       time.Minute * 10,
		TLSHandshakeTimeout:   time.Minute * 1,
		ExpectContinueTimeout: time.Second * 30,
	})

	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	c.OnHTML("article", func(e *colly.HTMLElement) {
		var newsBody string

		visitURL := e.ChildAttr("a[href]", "href")
		secondWorker := c.Clone()

		secondWorker.OnHTML("div[id=konten-artikel]", func(e *colly.HTMLElement) {
			e.ForEach("p", func(_ int, el *colly.HTMLElement) {
				newsBody = fmt.Sprintf("%s\n%s", newsBody, el.Text)
			})
		})

		err := secondWorker.Visit(visitURL)
		if err != nil {
			R = nil
			E = fmt.Errorf("[NEWS SCRAPPER] Second Worker Had Some Error : %v", err.Error())
		}

		Newses = append(Newses, s.News{
			OldCreatedAt: e.ChildText("time[pubdate]"),
			Title:        e.ChildText("h4"),
			ImageURL:     e.ChildAttr("img[src]", "data-original"),
			RealURL:      e.ChildAttr("a[href]", "href"),
			Content:      newsBody,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("[NEWS SCRAPPER] Visiting", r.URL.String())
	})

	err := c.Visit("https://covid19.go.id/p/berita")
	if err != nil {
		return nil, fmt.Errorf("[NEWS SCRAPPER] Error on visit %v \n", err.Error())
	}

	if E != nil {
		return R, E
	}

	return Newses, E
}

func Post(newses *[]s.News) ([]byte, error) {
	apiUrl := u.URLJoin(config.GetApiURL(), urlPath, "batch")

	jsonStr := []byte(u.Stringify(newses))

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	//Defer from response
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

func GetLastNewsTitle() (string, error) {
	var n s.ResponseJSON

	apiUrl := u.URLJoin(config.GetApiURL(), urlPath, "-", "last")
	res, err := http.Get(apiUrl)

	if err != nil {
		return "", err
	}

	data, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(data, &n)
	if err != nil {
		return "", err
	}

	d := n.Data.(map[string]interface{})
	return d["Title"].(string), nil
}

func ScrapAndPost(lastTitle string) ([]byte, error) {
	newses, err := Scrap()

	if err != nil {
		return nil, err
	}

	reversedArray := u.ReverseNewsArray(newses)
	startIndex := 0

	for i, v := range reversedArray {
		if v.Title == lastTitle {
			startIndex = i
		}
	}

	if startIndex == len(reversedArray) {
		return nil, nil
	}

	newArray := reversedArray[startIndex+1:]

	if len(newArray) <= 0 {
		return nil, fmt.Errorf("0 Post")
	}

	resp, err := Post(&newArray)
	if err != nil {
		return nil, fmt.Errorf("Post error : %e", err)
	}

	return resp, nil
}

func Start(interval time.Duration) interface{} {
	log.Printf("[NEWS ENGINE] Waking Up")
	lT, err := GetLastNewsTitle()
	if err != nil {
		log.Printf("[NEWS ENGINE] ERROR : %s | RETRYING", err)
		return Start(interval)
	}

	_, err = ScrapAndPost(lT)
	if err != nil {
		log.Printf("[NEWS ENGINE] ERROR : %s | RETRYING", err)
		if err.Error() == "0 Post" {
			log.Printf("[NEWS ENGINE] Sleeping for %s\n", interval)
			time.Sleep(interval)
		}
		return Start(interval)
	}

	log.Printf("[NEWS ENGINE] Scrapped and Posted Sleeping for %s\n", interval)
	time.Sleep(interval)
	return Start(interval)
}
