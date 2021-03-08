package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
)

const (
	serverAddr       = "0.0.0.0:8080"
	blizzardTokenURL = "https://us.battle.net/oauth/token"
	hearthstoneURL   = "https://us.api.blizzard.com/hearthstone"
	clientID         = ""
	clientSecret     = ""
)

type server struct {
	c *hearthstoneClient
}

func main() {
	hsc, err := newHearthstoneClient()
	if err != nil {
		log.Fatalf("Creating hearthstone client failed: %v", err)
	}

	s := server{hsc}

	log.Printf("Starting HTTP server on %s\n", serverAddr)
	http.HandleFunc("/", s.tableHandler)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func (s *server) tableHandler(w http.ResponseWriter, r *http.Request) {
	cards, err := s.c.search("druid,warlock", "legendary", "7,8,9,10")
	if err != nil {
		log.Printf("Hearthstone card search failed: %v", err)
	}

	log.Printf("Fetched %d cards for table request\n", len(cards))

	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].ID < cards[j].ID
	})

	t, _ := template.ParseFiles("table.html")
	t.Execute(w, struct {
		Cards []hearthstoneCard
	}{Cards: cards})
}
