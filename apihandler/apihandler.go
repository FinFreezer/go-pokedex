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

const interval = (20 * time.Second)
var pcache *c.Cache
var lastArea *Area
var pdex *Pokedex
var targetData *Pokemon

type Config struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:name`
		URL string `json:"url"`
	} `json:"results"`
}

type Area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func Start(config *Config, direction string, addParams []string) (*Config){
	if pcache == nil {
		pcache = c.NewCache(interval)
	}
	if pdex == nil {
		pdex = &Pokedex{Data: make(map[string]Pokemon)}
	}
	locMap, err := processAPIcall(config, direction, pcache, addParams)

	if err != nil {
		fmt.Println("Something went wrong with API call.")
		return locMap
	}

	return locMap
}

func processAPIcall(config *Config, direction string, cache *c.Cache, addParams []string) (*Config, error) {
	URL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	switch direction {
		case "Init":
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
		
		case "Next":
			config, err := nextBatch(config, cache)
			if err != nil {
				println("Error getting next batch.")
				return config, err
			}
			return config, err

		case "Previous":
			config, err := previousBatch(config, cache)
			if err != nil {
				println("Error getting previous batch.")
				return config, err
			}
			return config, err
		
		case "Explore":
			err := getLocDetails(config, cache, addParams)
			return config, err
		
		case "Catch":
			getPokeDetails(cache, addParams)
			catchPokemon(addParams[0])
			return config, nil
		case "Inspect":
			inspectPokemon(addParams[0])
			return config, nil
		case "Pokedex":
			displayPokedex()
			return config, nil
		default:
			return config, errors.New("Unexpected error with *config")
	}
}

func getPokeDetails(cache *c.Cache, addParams []string) error {
	target := addParams[0]
	URL := "https://pokeapi.co/api/v2/pokemon/" + target

	if _, ok := cache.Get(URL); ok {
		fmt.Println("Accessing cached data for Pokemon for URL: " + URL)
		data, _ := cache.Get(URL)
		json.Unmarshal(data, &targetData)
		return nil
	}

	res, err := http.Get(URL)
	defer res.Body.Close()

	if err != nil {
		fmt.Println("Error getting API call for target.")
		return err
	}

	data, error := io.ReadAll(res.Body)
	json.Unmarshal(data, &targetData)
	
	if error != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	cache.Add(URL, data)
	return nil
}

func getLocDetails(config *Config, cache *c.Cache, addParams []string) (error) {
	URL := "https://pokeapi.co/api/v2/location-area/" + addParams[0]

	if _, ok := cache.Get(URL); ok {
		fmt.Println("Accessing cached data for Location for URL: " + URL)
		data, _ := cache.Get(URL)
		json.Unmarshal(data, &lastArea)
		printLocDetails()
		return nil
	}

	res, err := http.Get(URL)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	data, error := io.ReadAll(res.Body)
	json.Unmarshal(data, &lastArea)
	
	if error != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	printLocDetails()
	cache.Add(URL, data)
	return nil
}

func printLocDetails() {
	if lastArea == nil {
		fmt.Println("No locations found, continuing.")
		return
	}
	fmt.Println("Exploring " + lastArea.Name + "...")
	fmt.Println("Found Pokemon: ")
	for _, encounters := range lastArea.PokemonEncounters {
		fmt.Println("- " + encounters.Pokemon.Name)
	}

	return
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