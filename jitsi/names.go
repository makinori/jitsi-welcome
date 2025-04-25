package jitsi

import (
	_ "embed"
	"encoding/json"
	"math/rand"

	"github.com/charmbracelet/log"
)

// https://github.com/jitsi/js-utils/blob/master/random/roomNameGenerator.js

type namesType struct {
	PluralNoun []string `json:"pluralnoun"`
	Verb       []string `json:"verb"`
	Adverb     []string `json:"adverb"`
	Adjective  []string `json:"adjective"`
}

var (
	//go:embed names.json
	namesJSON []byte

	names namesType = loadNames()
)

func loadNames() namesType {
	var names namesType

	err := json.Unmarshal(namesJSON, &names)
	if err != nil {
		log.Fatal("failed to load names", "err", err)
	}

	return names
}

func GenerateRoomName() string {
	var name string

	categories := [][]string{
		names.Adjective,
		names.PluralNoun,
		names.Verb,
		names.Adverb,
	}

	for _, category := range categories {
		word := category[rand.Intn(len(category))]
		name += word
	}

	return name
}
