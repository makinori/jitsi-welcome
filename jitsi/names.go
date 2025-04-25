package jitsi

import (
	_ "embed"
	"encoding/json"
	"math/rand"

	"github.com/charmbracelet/log"
)

// https://github.com/jitsi/js-utils/blob/master/random/roomNameGenerator.js

type namesTypes struct {
	PluralNoun []string `json:"pluralnoun"`
	Verb       []string `json:"verb"`
	Adverb     []string `json:"adverb"`
	Adjective  []string `json:"adjective"`
}

//go:embed names.json
var namesJSON []byte

var names namesTypes = loadNames()

func loadNames() namesTypes {
	var names namesTypes

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
