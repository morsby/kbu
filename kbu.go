package kbu

import (
	"encoding/json"
	"io"
	"strings"
	"time"
)

type RawFormat struct {
	URL              string `json:"url"`
	Valgt            string `json:"Valgt"`
	Lodtr            string `json:"Lodtr."`
	Region           string `json:"Region"`
	Startdato        string `json:"Startdato"`
	Uddannelsessted  string `json:"Uddannelsessted"`
	Afdeling         string `json:"Afdeling"`
	Speciale         string `json:"Speciale"`
	Uddannelsessted2 string `json:"Uddannelsessted2"`
	Afdeling2        string `json:"Afdeling2"`
	Speciale2        string `json:"Speciale2"`
}

type Season int

const (
	SeasonFall Season = iota + 1
	SeasonSpring
)

type Round struct {
	Season Season
	Year   int
}

type Selection struct {
	Runde            Round     `json:"runde"`
	Dato             time.Time `json:"date"`
	Universitet      string    `json:"universitet"`
	Nummer           int       `json:"no"`
	Region           string    `json:"region"`
	Startdato        time.Time `json:"startdato"`
	Uddannelsessted  string    `json:"uddannelsessted"`
	Afdeling         string    `json:"afdeling"`
	Speciale         string    `json:"speciale"`
	Uddannelsessted2 string    `json:"uddannelsessted2"`
	Afdeling2        string    `json:"afdeling2"`
	Speciale2        string    `json:"speciale2"`
	URL              string    `json:"url"`
}

const latestUrl string = "https://kbu.logbog.net/AJAX_Timelines.asp"

func ParseRawJSON(r io.Reader) ([]Selection, error) {
	var raw []RawFormat
	dec := json.NewDecoder(r)
	dec.Decode(&raw)

	data := make([]Selection, len(raw))
	for i, r := range raw {
		round, err := calculateRound(r.URL)
		if err != nil {
			return nil, err
		}

		number, uni, err := calculateNumberUni(r.Lodtr)
		if err != nil {
			return nil, err
		}

		startdato, err := calculateDate(r.Startdato, 0)
		if err != nil {
			return nil, err
		}
		year := startdato.Year()
		if startdato.Month() < 7 {
			year--
		}
		dato, err := calculateDate(r.Valgt, year)

		if err != nil {
			return nil, err
		}

		data[i] = Selection{
			Runde:            round,
			URL:              strings.TrimSpace(r.URL),
			Universitet:      uni,
			Nummer:           number,
			Dato:             dato,
			Region:           strings.TrimSpace(r.Region),
			Startdato:        startdato,
			Uddannelsessted:  strings.TrimSpace(r.Uddannelsessted),
			Afdeling:         strings.TrimSpace(r.Afdeling),
			Speciale:         strings.TrimSpace(r.Speciale),
			Uddannelsessted2: strings.TrimSpace(r.Uddannelsessted2),
			Afdeling2:        strings.TrimSpace(r.Afdeling2),
			Speciale2:        strings.TrimSpace(r.Speciale2),
		}
	}
	return data, nil
}
