package apihandler

import 
(
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"errors"
	c "github.com/FinFreezer/go-pokedex/internal/pokecache"
	"time"
)

const interval = (60 * time.Second)
var pcache *c.Cache

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:name`
		URL string `json:"url"`
	} `json:"results"`
}

func Start(config *Config, direction string) (*Config){
	if pcache == nil {
		pcache = c.NewCache(interval)
	}
	locMap, err := getLocBatch(config, direction, pcache)

	if err != nil {
		fmt.Println("Error getting locations.")
		return locMap
	}

	return locMap
}

func printBatch(locMap *Config) {
	if locMap == nil {
		fmt.Println("No locations found, continuing.")
		return
	}
	for _, loc := range locMap.Results {
		fmt.Println(loc.Name)
	}
	return
}

func getLocBatch(config *Config, direction string, cache *c.Cache) (*Config, error) {
	URL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	if direction == "Init" {
		res, err := http.Get(URL)
		defer res.Body.Close()

		if err != nil {
			return config, err
		}

		data, error := io.ReadAll(res.Body)
		json.Unmarshal(data, &config)
		
		if error != nil {
			return config, fmt.Errorf("error creating request: %w", err)
		}

		cache.Add(URL, data)
		printBatch(config)
		return config, nil
	
	} else if direction == "Next" {
		config, err := nextBatch(config, cache)
		if err != nil {
			println("Error getting next batch.")
			return config, err
		}
		return config, err
	
	} else if direction == "Previous" {
		config, err := previousBatch(config, cache)
		if err != nil {
			println("Error getting previous batch.")
			return config, err
		}
		return config, err
	
	} else {
		return config, errors.New("Unexpected error with *config")
	}
}

func nextBatch(old *Config, cache *c.Cache) (*Config, error) {
	URL := old.Next

	if _, ok := cache.Get(URL); ok {
		fmt.Println("Accessing cached data for Next for URL: " + URL)
		data, _ := cache.Get(URL)
		json.Unmarshal(data, &old)
		printBatch(old)
		return old, nil
	}

	res, err := http.Get(URL)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Error getting next batch from URL.")
		return old, err
	}
	data, error := io.ReadAll(res.Body)
	cache.Add(URL, data)
	json.Unmarshal(data, &old)
	if error != nil {
		return old, fmt.Errorf("error creating request: %w", err)
	}

	printBatch(old)
	return old, nil
}

func previousBatch(old *Config, cache *c.Cache) (*Config, error) {
	if old == nil {
		fmt.Println("No location set, please run 'map' first.")
		return old, nil
	}
	if old.Previous == nil {
		fmt.Println("you're on the first page")
		return old, nil
	
	} else if old.Previous != nil {
    	previousURL, ok := old.Previous.(string)

		if _, ok := cache.Get(previousURL); ok {
			fmt.Println("Accessing cached data for Previous for URL: " + previousURL)
			data, _ := cache.Get(previousURL)
			json.Unmarshal(data, &old)
			printBatch(old)
			return old, nil
		}
		
		if !ok {
			fmt.Println("Error: old.Previous was not a string type")
			return old, nil

		} else {
			res, err := http.Get(previousURL)
			defer res.Body.Close()

			if err != nil {
				return old, err
			}

			data, error := io.ReadAll(res.Body)
			cache.Add(previousURL, data)
			json.Unmarshal(data, &old)

			if error != nil {
				return old, fmt.Errorf("error creating request: %w", err)
			}

			printBatch(old)
			return old, nil
		}
	
	} else {
		return old, errors.New("Error getting previous batch.")
	}
}