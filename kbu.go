package kbu

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const latestUrl string = "https://kbu.logbog.net/AJAX_Timelines.asp"

func ParseRawJSON(r io.Reader) ([]Selection, error) {
	var raw []RawFormat
	dec := json.NewDecoder(r)
	dec.Decode(&raw)

	rounds := make(map[string]int)

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
		// if the selection wasn't taken, it should not add to number of doctors that year
		if number > -1 {
			rounds[fmt.Sprintf("%d-%s", round.Year, round.Season)]++
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
			Round:      round,
			University: uni,
			Number:     number,
			Date:       dato,
			Region:     region,
			Start:      startdato,
			Positions: []Position{
				{
					Location:   strings.TrimSpace(r.Uddannelsessted),
					Department: strings.TrimSpace(r.Afdeling),
					Specialty:  strings.TrimSpace(r.Speciale),
				},
				{
					Location:   strings.TrimSpace(r.Uddannelsessted2),
					Department: strings.TrimSpace(r.Afdeling2),
					Specialty:  strings.TrimSpace(r.Speciale2),
				},
			},
		}
		s.ID = s.GenerateID()

		data[i] = s
	}

	for i := range data {
		data[i].RelNumber = 1.0 - float64(data[i].Number)/float64(rounds[fmt.Sprintf("%d-%s", data[i].Round.Year, data[i].Round.Season)])
	}

	return data, nil
}

func FlattenSelections(selections []Selection) []SelectionFlat {
	flattened := make([]SelectionFlat, 0, len(selections))
	for _, sel := range selections {
		flattened = append(flattened, sel.Flatten())
	}
	return flattened
}
