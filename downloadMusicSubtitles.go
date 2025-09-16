package main

import (
	"encoding/json"
	"fmt"
	"go.senan.xyz/taglib"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// this is going to go through every file in this directory and subdirectories,
// check if subtitles exist for it's music name metadata and then find subtitiles
// then download them
func main() {
	root := "."

	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() == false {
			if strings.HasSuffix(d.Name(), ".flac") ||
				strings.HasSuffix(d.Name(), ".wav") {

				// gets subtitles from musixmatch with formatted name
				getSubtitles(getTitle(d.Name(), path))
			}
		}
		return nil
	})
}

func getTitle(song string, path string) (string, string) {

	// gets metadata title, fallsback to song name if none
	songData, err := taglib.ReadTags(path) // reading metadata

	if err != nil || songData[taglib.Title] == nil {
		// captures just the song title using regex
		re := regexp.MustCompile(`^\d+\s*\.?-?\s*(.*?)(?:\.[^.]+)?$`)
		title := re.FindStringSubmatch(song)
		if len(title) > 1 {
			return title[1], "Unknown"
		}
	}
	title := songData[taglib.Title][0]
	artist := songData[taglib.Artist][0]
	return title, artist
}

func getSubtitles(title string, artist string) {
	// get the song id
	if artist == "Unknown" {
		fmt.Println("no artist:", title)
		return
	}
	getSubtitlesURL := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artist, title)
	// getSubtitlesURL := "https://api.lyrics.ovh/v1/Radiohead/Let Down"
	fmt.Println("Attempting:", getSubtitlesURL)

	requ, _ := http.NewRequest("GET", getSubtitlesURL, nil)

	resp, err := http.DefaultClient.Do(requ)

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err, title)
		fmt.Println(resp)
	}

	body, _ := io.ReadAll(resp.Body) // returns in format []byte

	// fmt.Println(string(body))
	fmt.Println(convertJSONtoText(body))
}

type lyricsStruct struct {
	// when lyricsStruct.Lyrics is used: it will get the json
	// key "Lyrics" to then save its corresponding value
	Lyrics string `json:"lyrics"`
	// `` means its a tag
}

func convertJSONtoText(body []byte) string {

	// defining variable with struct to be able to access lyrics key of json
	var response lyricsStruct

	// takes in json mapped to []byte -- go's json library works with byte slices -> []byte
	// stores json encoded data to pointer response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Sprintf("Failed: %s\n%s", err, string(body))
	}

	formattedBody := response.Lyrics // gets the 'lyrics' key in json

	fmt.Println(formattedBody)

	return formattedBody
}
