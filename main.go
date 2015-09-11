package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/tent/http-link-go"
)

type gist struct {
	URL    string `json:"url"`
	Public bool   `json:"public"`
}

func main() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatal("Missing GH_TOKEN env var. Create API token with Gist access here: https://github.com/settings/tokens")
	}
	auth := "token " + token

	dryrun := flag.Bool("dryrun", true, "don't actually delete anything")

	flag.Parse()

	// List of gists to delete
	todel := make([]string, 0)

	// Page through all gists and build up list of non-public gists to delete
	pageurl := "https://api.github.com/gists"
	next := true
	for next {
		req, err := http.NewRequest("GET", pageurl, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", auth)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 200 {
			log.Fatalf("non-200 response: %d", resp.StatusCode)
		}

		page := make([]gist, 0)
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			log.Fatalf("error decoding response body: %s", err)
		}

		// Append non-public urls to the delete list
		for _, g := range page {
			if !g.Public {
				todel = append(todel, g.URL)
			}
		}

		// Get next page from link headers
		next = false

		linkheader := resp.Header.Get("Link")
		if linkheader == "" {
			continue
		}
		links, err := link.Parse(linkheader)
		if err != nil {
			log.Fatalf("error parse link header %q: %s", linkheader, err)
		}

		for _, link := range links {
			if link.Rel == "next" {
				pageurl = link.URI
				next = true
				break
			}
		}
	}

	if *dryrun {
		log.Printf("Dryrun, not deleting %d private gists", len(todel))
		log.Printf("Run with -dryrun=false to delete them")
		return
	}

	// Now go through and delete all private gists
	for _, url := range todel {
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Fatalf("unable to create request to delete %s: %v", url, err)
		}

		req.Header.Set("Authorization", auth)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 204 {
			log.Fatalf("non-204 status code when deleting %s: %d", url, resp.StatusCode)
		}
	}

	log.Printf("Deleted %d private gists", len(todel))
}
