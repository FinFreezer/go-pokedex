package apihandler

import 
(
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"errors"
)

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
	locMap, err := getLocBatch(config, direction)

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

func getLocBatch(config *Config, direction string) (*Config, error) {
	if direction == "Init" {
		res, err := http.Get("https://pokeapi.co/api/v2/location-area")

		if err != nil {
			return config, err
		}

		data, error := io.ReadAll(res.Body)
		json.Unmarshal(data, &config)
		
		if error != nil {
			return config, fmt.Errorf("error creating request: %w", err)
		}

		defer res.Body.Close()
		printBatch(config)
		return config, nil
	
	} else if direction == "Next" {
		config, err := nextBatch(config)
		if err != nil {
			println("Error getting next batch.")
			return config, err
		}
		return config, err
	
	} else if direction == "Previous" {
		config, err := PreviousBatch(config)
		if err != nil {
			println("Error getting previous batch.")
			return config, err
		}
		return config, err
	
	} else {
		return config, errors.New("Unexpected error with *config")
	}
}

func nextBatch(old *Config) (*Config, error) {
	var batch2 *Config
	URL := old.Next
	res, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error getting next batch from URL.")
		return old, err
	}
	data, error := io.ReadAll(res.Body)
	json.Unmarshal(data, &batch2)
	if error != nil {
		return old, fmt.Errorf("error creating request: %w", err)
	}

	defer res.Body.Close()
	printBatch(batch2)
	return batch2, nil
}

func PreviousBatch(old *Config) (*Config, error) {
	var batch2 *Config
	if old == nil {
		fmt.Println("No location set, please run 'map' first.")
		return old, nil
	}
	if old.Previous == nil {
		fmt.Println("you're on the first page")
		return old, nil
	
	} else if old.Previous != nil {
    	previousURL, ok := old.Previous.(string)
		
		if !ok {
			fmt.Println("Error: old.Previous was not a string type")
			return old, nil

		} else {
			res, err := http.Get(previousURL)

			if err != nil {
				return batch2, err
			}

			data, error := io.ReadAll(res.Body)
			json.Unmarshal(data, &batch2)

			if error != nil {
				return batch2, fmt.Errorf("error creating request: %w", err)
			}

			defer res.Body.Close()
			printBatch(batch2)
			return batch2, nil
		}
	
	} else {
		return old, errors.New("Error getting previous batch.")
	}
}
