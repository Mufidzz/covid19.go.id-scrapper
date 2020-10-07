package Protocol

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

var urlPath = "/protocol"

func Scrap() ([]s.Protocol, error) {
	var (
		Protocols []s.Protocol
		R         []s.Protocol
		E         error
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
		var protocol s.Protocol

		visitUrl := e.ChildAttr("a[href]", "href")

		protocol.OldCreatedAt = e.ChildText("time[pubdate]")
		protocol.Title = e.ChildText("h4")
		protocol.ImageURL = e.ChildAttr("img[src]", "data-original")
		protocol.RealURL = e.ChildAttr("a[href]", "href")

		secondWorker := c.Clone()
		secondWorker.OnHTML("div[id=konten-artikel]", func(el *colly.HTMLElement) {
			el.ForEach("p", func(_ int, ele *colly.HTMLElement) {
				if ele.Text == "Detail:" || ele.Text == "Unduh Materi" {
					protocol.DownloadURL = ele.ChildAttr("a[href]", "href")
				} else {
					protocol.Content = fmt.Sprintf("%s\n%s", protocol.Content, ele.Text)
				}
			})
		})

		err := secondWorker.Visit(visitUrl)
		if err != nil {
			E = err
		}

		Protocols = append(Protocols, protocol)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("[PROTOCOL SCRAPPER] Visiting", r.URL.String())
	})

	err := c.Visit("https://covid19.go.id/p/protokol")
	if err != nil {
		return nil, fmt.Errorf("[PROTOCOL SCRAPPER] Error on visit %v \n", err.Error())
	}

	if E != nil {
		return R, E
	}

	return Protocols, E
}

func Post(protocol *[]s.Protocol) ([]byte, error) {
	jwt, err := u.GenerateJWT()
	if err != nil {
		return nil, err
	}

	apiUrl := u.URLJoin(config.GetApiURL(), urlPath, "batch")

	jsonStr := []byte(u.Stringify(protocol))

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", jwt)

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

func GetLastProtocolTitle() (string, error) {
	var n s.ResponseJSON
	jwt, err := u.GenerateJWT()

	apiUrl := u.URLJoin(config.GetApiURL(), urlPath, "-", "last")

	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Token", jwt)

	client := &http.Client{}
	res, err := client.Do(req)

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
	protocol, err := Scrap()

	if err != nil {
		return nil, err
	}

	reversedArray := u.ReverseProtocolArray(protocol)
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
	log.Printf("[PROTOCOL ENGINE] Waking Up")
	lT, err := GetLastProtocolTitle()
	if err != nil {
		log.Printf("[PROTOCOL ENGINE] ERROR : %s | RETRYING", err)
		return Start(interval)
	}

	_, err = ScrapAndPost(lT)
	if err != nil {
		log.Printf("[PROTOCOL ENGINE] ERROR : %s | RETRYING", err)
		if err.Error() == "0 Post" {
			log.Printf("[PROTOCOL ENGINE] Sleeping for %s\n", interval)
			time.Sleep(interval)
		}
		return Start(interval)
	}

	log.Printf("[PROTOCOL ENGINE] Scrapped and Posted Sleeping for %s\n", interval)
	time.Sleep(interval)
	return Start(interval)
}
