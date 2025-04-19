package anime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/makinori/jitsi-welcome/common"

	"github.com/charmbracelet/log"
)

// 50 is max per page

const graphqlQuery = `
fragment FullName on CharacterConnection {
  nodes {
    name {
      full
    }
  }
}

query($userName: String, $page: Int)  {
  Page(page: $page, perPage: 50) {
    mediaList(
      userName: $userName, 
      status_in: [CURRENT, COMPLETED], 
      type: ANIME
    ) {
      media { 
        title {
          english
		  romaji
        }
        main: characters(role: MAIN) {
          ...FullName
        }
        supporting: characters(role: SUPPORTING) {
          ...FullName
        }
      }
    }
  }
}
`

type GraphqlFullName struct {
	Nodes []struct {
		Name struct {
			Full string `json:"full"`
		} `json:"name"`
	} `json:"nodes"`
}

type GraphqlResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	} `json:"errors,omitempty"`
	Data struct {
		Page struct {
			MediaList []struct {
				Media struct {
					Title struct {
						English string `json:"english"`
						Romaji  string `json:"romaji"`
					} `json:"title"`
					Main       GraphqlFullName `json:"main"`
					Supporting GraphqlFullName `json:"supporting"`
				} `json:"media"`
			} `json:"mediaList"`
		} `json:"Page"`
	} `json:"data"`
}

type AnimeNames struct {
	Names  []string `json:"names"`
	Titles []string `json:"titles"`
}

func getAnimeNamesPage(username string, page int) (AnimeNames, error) {
	inputData := struct {
		Query     string            `json:"query"`
		Variables map[string]string `json:"variables"`
	}{
		Query: strings.TrimSpace(graphqlQuery),
		Variables: map[string]string{
			"userName": username,
			"page":     strconv.Itoa(page),
		},
	}

	inputJson, err := json.Marshal(inputData)
	if err != nil {
		return AnimeNames{}, err
	}

	req, err := http.NewRequest(
		"POST", "https://graphql.anilist.co",
		bytes.NewBuffer(inputJson),
	)

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return AnimeNames{}, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return AnimeNames{}, err
	}
	defer res.Body.Close()

	outputBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return AnimeNames{}, err
	}

	var graphqlRes GraphqlResponse
	err = json.Unmarshal(outputBytes, &graphqlRes)
	if err != nil {
		return AnimeNames{}, err
	}

	if len(graphqlRes.Errors) > 0 {
		return AnimeNames{}, fmt.Errorf("%s", graphqlRes.Errors)
	}

	var output AnimeNames

	for _, anime := range graphqlRes.Data.Page.MediaList {
		if anime.Media.Title.English != "" {
			output.Titles = append(output.Titles, anime.Media.Title.English)
		} else {
			output.Titles = append(output.Titles, anime.Media.Title.Romaji)
		}

		for _, character := range anime.Media.Main.Nodes {
			if slices.Index(output.Names, character.Name.Full) == -1 {
				output.Names = append(output.Names, character.Name.Full)
			}
		}
		for _, character := range anime.Media.Supporting.Nodes {
			if slices.Index(output.Names, character.Name.Full) == -1 {
				output.Names = append(output.Names, character.Name.Full)
			}
		}
	}

	return output, nil
}

func getAnimeNames(username string) (AnimeNames, error) {
	page := 1

	var output AnimeNames

	for {
		log.Infof("fetching page %d", page)

		found, err := getAnimeNamesPage(username, page)
		if err != nil {
			return AnimeNames{}, err
		}

		if len(found.Names) == 0 {
			break
		}

		output.Titles = append(output.Titles, found.Titles...)

		for _, name := range found.Names {
			if slices.Index(output.Names, name) == -1 {
				output.Names = append(output.Names, name)
			}
		}

		page++
	}

	log.Infof(
		"found %d names from %d animes",
		len(output.Names), len(output.Titles),
	)

	return output, nil
}

type AnimeNamesCache struct {
	Expire time.Time  `json:"expire"`
	Data   AnimeNames `json:"data"`
}

func setCachedAnimeNames(data AnimeNames) error {
	cacheData := AnimeNamesCache{
		Expire: time.Now().Add(time.Hour * 24 * 7),
		Data:   data,
	}

	cacheJSON, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	err = os.WriteFile(common.ConfigCacheJSONPath, cacheJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getCachedAnimeNames() (AnimeNames, error) {
	cacheJSON, err := os.ReadFile(common.ConfigCacheJSONPath)
	if err != nil {
		return AnimeNames{}, err
	}

	var cacheData AnimeNamesCache
	err = json.Unmarshal(cacheJSON, &cacheData)
	if err != nil {
		return AnimeNames{}, err
	}

	if time.Now().After(cacheData.Expire) {
		return AnimeNames{}, errors.New("cache expired")
	}

	return cacheData.Data, nil
}

func GetRandomAnimeName(username string) (string, error) {
	var data AnimeNames

	data, err := getCachedAnimeNames()
	if err != nil {
		log.Warn("failed to get cache", "err", err)

		data, err = getAnimeNames(username)
		if err != nil {
			return "", err
		}

		err = setCachedAnimeNames(data)
		if err != nil {
			log.Error("failed to set cache", "err", err)
		}
	}

	i := rand.Intn(len(data.Names))

	return data.Names[i], nil
}
