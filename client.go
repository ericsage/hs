package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type metadataJSON struct {
	Sets     []assocJSON
	Types    []assocJSON
	Rarities []assocJSON
	Classes  []assocJSON
}

type assocJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type cardsJSON struct {
	Cards []cardJSON `json:"cards"`
}

type cardJSON struct {
	ID         int    `json:"id"`
	ClassID    int    `json:"classId"`
	CardTypeID int    `json:"cardTypeId"`
	CardSetID  int    `json:"cardSetId"`
	RarityID   int    `json:"rarityId"`
	Name       string `json:"name"`
	ImageURL   string `json:"image"`
}

type hearthstoneCard struct {
	ID       int
	Class    string
	Type     string
	Set      string
	Rarity   string
	Name     string
	ImageURL string
}

type hearthstoneClient struct {
	token string
	c     *http.Client
	m     hearthstoneMetadata
}

type hearthstoneMetadata struct {
	sets     map[int]string
	types    map[int]string
	rarities map[int]string
	classes  map[int]string
}

func newHearthstoneClient() (*hearthstoneClient, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("Could not fetch API access token: %w", err)
	}

	hsc := &hearthstoneClient{token: token, c: http.DefaultClient}
	metadata, err := hsc.fetchMetadata()
	if err != nil {
		return nil, fmt.Errorf("Could not fetch hearthstone metadata: %w", err)
	}

	hsc.m = hearthstoneMetadata{
		sets:     assocsToMap(metadata.Sets),
		types:    assocsToMap(metadata.Types),
		rarities: assocsToMap(metadata.Rarities),
		classes:  assocsToMap(metadata.Classes),
	}

	return hsc, nil
}

func (hsc *hearthstoneClient) search(classes, rarities, manaCost string) ([]hearthstoneCard, error) {
	cardsJSON, err := hsc.fetchCards(map[string]string{
		"class":    classes,
		"rarity":   rarities,
		"manaCost": manaCost,
	})
	if err != nil {
		return nil, err
	}

	var cards []hearthstoneCard
	for _, cj := range cardsJSON {
		cards = append(cards, hearthstoneCard{
			ID:       cj.ID,
			Set:      hsc.m.sets[cj.CardSetID],
			Type:     hsc.m.types[cj.CardTypeID],
			Rarity:   hsc.m.rarities[cj.RarityID],
			Class:    hsc.m.classes[cj.ClassID],
			Name:     cj.Name,
			ImageURL: cj.ImageURL,
		})
	}
	return cards, nil
}

func (hsc *hearthstoneClient) fetchCards(params map[string]string) ([]cardJSON, error) {
	req, err := hsc.buildRequest("/cards/", params)

	resp, err := hsc.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	var cards cardsJSON
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&cards)
	if err != nil {
		return nil, err
	}

	return cards.Cards, nil
}

func (hsc *hearthstoneClient) fetchMetadata() (*metadataJSON, error) {
	req, err := hsc.buildRequest("/metadata/", map[string]string{})

	resp, err := hsc.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}

	var metadata metadataJSON
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (hsc *hearthstoneClient) buildRequest(path string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("GET", hearthstoneURL+path, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("locale", "en_US")
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hsc.token))
	return req, nil
}

func assocsToMap(assocs []assocJSON) map[int]string {
	m := make(map[int]string)
	for _, assoc := range assocs {
		m[assoc.ID] = assoc.Name
	}
	return m
}
