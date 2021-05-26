package kbu

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var ids map[string]bool = make(map[string]bool)

func (s *Selection) GenerateID() string {
	id := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", s))))

	for _, ok := ids[id]; ok; _, ok = ids[id] {
		id += "a"
	}
	ids[id] = true

	return id
}

const latestUrl string = "https://kbu.logbog.net/AJAX_Timelines.asp"

func ParseRawJSON(r io.Reader) ([]Selection, error) {
	var raw []RawFormat
	dec := json.NewDecoder(r)
	dec.Decode(&raw)

	data := make([]Selection, len(raw))
	for i, r := range raw {
		round, err := getRound(r.URL)
		if err != nil {
			return nil, err
		}

		number, uni, err := getNumberUni(r.Lodtr)
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

		region, err := getRegion(r.Region)
		if err != nil {
			return nil, err
		}

		s := Selection{
			Round:       round,
			URL:         strings.TrimSpace(r.URL),
			University:  uni,
			Number:      number,
			Date:        dato,
			Region:      region,
			Start:       startdato,
			Place1:      strings.TrimSpace(r.Uddannelsessted),
			Department1: strings.TrimSpace(r.Afdeling),
			Specialty1:  strings.TrimSpace(r.Speciale),
			Place2:      strings.TrimSpace(r.Uddannelsessted2),
			Department2: strings.TrimSpace(r.Afdeling2),
			Specialty2:  strings.TrimSpace(r.Speciale2),
		}
		s.ID = s.GenerateID()

		data[i] = s
	}
	return data, nil
}
