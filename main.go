package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func fetchWeather(client *http.Client, city string, ch chan<- string, wg *sync.WaitGroup) {
	var data struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	defer wg.Done()

	const apikey = ""
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apikey)
	resp, err := client.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching weather for %s: %v", city, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- fmt.Sprintf("Error fetching weather for %s: %v", city, err)
	}

	ch <- fmt.Sprintf(" the city tempis %s %+v", city, data.Main.Temp-273)

}

func main() {
	times := time.Now()
	cities := []string{"chennai", "vellore", "delhi"}

	ch := make(chan string)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var wg sync.WaitGroup
	for _, city := range cities {
		wg.Add(1)
		go fetchWeather(client, city, ch, &wg)

	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
	fmt.Println("The operation took with:", time.Since(times))

}
