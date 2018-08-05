package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	Domain     = "http://www.piano-e-competition.com"
	MIDIFiles  = "/ecompetition/midifiles"
	MIDIFiles1 = "midifiles"
	MIDIFiles2 = "/midifiles"
)

func fetch(page string) {
	response, err := http.Get(page)
	if err != nil {
		panic(err)
	}
	tokenizer := html.NewTokenizer(response.Body)
	defer response.Body.Close()
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			fmt.Println(tokenizer.Err())
			return
		case html.StartTagToken:
			t := tokenizer.Token()
			for _, attribute := range t.Attr {
				if attribute.Key == "href" {
					href := strings.ToLower(attribute.Val)
					if strings.HasPrefix(href, MIDIFiles1) {
						href = "/ecompetition/" + href
					} else if strings.HasPrefix(href, MIDIFiles2) {
						href = "/ecompetition" + href
					}
					if strings.HasPrefix(href, MIDIFiles) && strings.HasSuffix(href, ".mid") {
						fmt.Println(href)
						_, err := os.Stat("." + href)
						if err != nil {
							i := strings.LastIndex(href, "/")
							dir := "." + href[:i]
							err := os.MkdirAll(dir, 0777)
							if err != nil {
								panic(err)
							}
							path := attribute.Val
							if !strings.HasPrefix(path, "/") {
								path = "/" + path
							}
							response, err := http.Get(Domain + path)
							if err != nil {
								panic(err)
							}
							file, err := os.Create("." + href)
							if err != nil {
								panic(err)
							}
							_, err = io.Copy(file, response.Body)
							if err != nil {
								panic(err)
							}
							file.Close()
							response.Body.Close()
						}
					}
				}
			}
		}
	}
}

var pages = []string{
	"http://www.piano-e-competition.com/midi_2002.asp",
	"http://www.piano-e-competition.com/midi_2004.asp",
	"http://www.piano-e-competition.com/midi_2006.asp",
	"http://www.piano-e-competition.com/midi_2008.asp",
	"http://www.piano-e-competition.com/midi_2009.asp",
	"http://www.piano-e-competition.com/midi_2011.asp",
}

func main() {
	for _, page := range pages {
		fetch(page)
		time.Sleep(3 * time.Second)
	}
}
