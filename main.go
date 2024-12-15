package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	arr := []string{
		"https://httpbin.org/get?proglib=coll",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/5"}

	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		result := make(chan string, 3)

		for _, URL := range arr {
			wg.Add(1)

			go func(URL string) {
				defer wg.Done()
				response, _ := client.Get(URL)
				if response != nil {
					defer response.Body.Close()
					body, _ := io.ReadAll(response.Body)
					bs := string(body)
					result <- bs
					return
				}
			}(URL)
		}
		wg.Wait()
		close((result))
		fmt.Println(time.Now())

		for response := range result {
			fmt.Println(response)
		}
	}
}
