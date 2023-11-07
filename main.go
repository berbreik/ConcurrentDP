package Concurrent_DP

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	count := 0
	// Fetch html form 5 sites
	input := []string{"https://www.google.com", "https://www.facebook.com", "https://www.twitter.com", "https://www.youtube.com", "https://www.instagram.com"}
	for _, inp := range input {
		wg.Add(1)
		go func(inp string) {
			// Fetch html
			defer wg.Done()
			fetch, err := FetchHTML(inp)
			if err != nil {
				return
			}
			ch <- fetch
		}(inp)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Count words
	//for i := range ch {
	//	count += i
	//}

	select {
	case i := <-ch:
		count += i
	case <-time.After(2 * time.Second):
		fmt.Println("Time out")
		close(ch)
	}

	fmt.Println(count)
}

func FetchHTML(url string) (int, error) {
	inp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(inp.Body)
	//Read REsponse body

	res, err := io.ReadAll(inp.Body)

	str := string(res)
	count := len(strings.Fields(str))
	//count words in response body

	return count, nil

}
