package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch ���� URL �� body ���ݣ����ҽ������ҳ�����ҵ��� URL �ŵ�һ�� slice �С�
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl ʹ�� fetcher ��ĳ�� URL ��ʼ�ݹ����ȡҳ�棬ֱ���ﵽ�����ȡ�
func Crawl(url string, depth int, fetcher Fetcher, crawled Crawled, out chan string, end chan bool) {
	if depth <= 0 {
		end <- true
		return
	}

	crawled.mux.Lock()
	if _, ok := crawled.crawled[url]; ok {
		crawled.mux.Unlock()
		end <- true
		return
	}

	crawled.crawled[url] = 1
	crawled.mux.Unlock()

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		end <- true
		return
	}

	out <- url
	//fmt.Println("found: ", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, crawled, out, end)
	}

	for i := 0; i < len(urls); i++ {
		<-end
	}

	end <- true
	return
}

type Crawled struct {
	crawled map[string]int
	mux     sync.Mutex
}

func main() {
	crawled := Crawled{make(map[string]int), sync.Mutex{}}
	out := make(chan string)
	end := make(chan bool)
	go Crawl("http://golang.org/", 4, fetcher, crawled, out, end)

	for {
		select {
		case url := <-out:
			fmt.Println("found: ", url)
		case <-end:
			return
		}
	}
}

// fakeFetcher �Ƿ������ɽ���� Fetcher��
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher ������� fakeFetcher��
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
