package kbu

import (
	"crypto/md5"
	"fmt"
	"time"
)

// RawFormat is the data format one gets when parsing the tables obtained at
// e.g. https://kbu.logbog.net/AJAX_Timelines.asp.
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

// Season identifies a season (fall or spring)
type Season string

const (
	SeasonSpring Season = "forår"
	SeasonFall   Season = "efterår"
)

// Region identifies valid regions
type Region string

const (
	RegionMidt Region = "Midt"
	RegionNord Region = "Nord"
	RegionSj   Region = "Sj"
	RegionH    Region = "H"
	RegionSyd  Region = "Syd"
)

// University and its constants identify the universities in the system
type University string

const (
	UniversityAU  University = "AU"
	UniversityKU  University = "KU"
	UniversitySDU University = "SDU"
	UniversityAAU University = "AAU"
	UniversityNA  University = "NA"
)

// Round contains information on a round
type Round struct {
	Season Season `json:"season"`
	Year   int    `json:"year"`
	URL    string `json:"url"`
}

// Position contains information on a Position
type Position struct {
	Location   string `json:"location"`
	Department string `json:"department"`
	Specialty  string `json:"specialty"`
}

// Selection contains information on a selection
type Selection struct {
	ID         string     `json:"id"`
	Round      Round      `json:"round"`
	Date       time.Time  `json:"date"`
	University University `json:"university"`
	Number     int        `json:"no"`
	RelNumber  float64    `json:"relNumber"`
	Region     Region     `json:"region"`
	Start      time.Time  `json:"start"`
	Positions  []Position `json:"positions"`
}

var ids map[string]bool = make(map[string]bool)

// GenerateID creates an ID for a Selection by calculating the
// md5 checksum of its fmt.Sprintf("%v", selection).
// IDs are unique, 'a's are appended if two selections were identical.
func (s *Selection) GenerateID() string {
	id := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", s))))

	for _, ok := ids[id]; ok; _, ok = ids[id] {
		id += "a"
	}
	ids[id] = true

	return id
}

// SelectionFloat contains information about af selection without any
// nested fields
type SelectionFlat struct {
	ID          string     `json:"id"`
	RoundYear   int        `json:"roundYear"`
	RoundSeason Season     `json:"roundSeason"`
	RoundURL    string     `json:"roundUrl"`
	Date        time.Time  `json:"date"`
	University  University `json:"university"`
	Number      int        `json:"no"`
	RelNumber   float64    `json:"relNumber"`
	Region      Region     `json:"region"`
	Start       time.Time  `json:"start"`
	Location1   string     `json:"location1"`
	Department1 string     `json:"department1"`
	Specialty1  string     `json:"specialty1"`
	Location2   string     `json:"location2"`
	Department2 string     `json:"department2"`
	Specialty2  string     `json:"specialty2"`
}

// Flatten flattens a Selection (i.e. flattens the nested struct Round and slice Positions)
func (s Selection) Flatten() SelectionFlat {
	flat := SelectionFlat{}
	flat.ID = s.ID
	flat.RoundYear = s.Round.Year
	flat.RoundSeason = s.Round.Season
	flat.RoundURL = s.Round.URL
	flat.Date = s.Date
	flat.University = s.University
	flat.Number = s.Number
	flat.RelNumber = s.RelNumber
	flat.Region = s.Region
	flat.Start = s.Start
	flat.Location1 = s.Positions[0].Location
	flat.Department1 = s.Positions[0].Department
	flat.Specialty1 = s.Positions[0].Specialty
	flat.Location2 = s.Positions[1].Location
	flat.Department2 = s.Positions[1].Department
	flat.Specialty2 = s.Positions[1].Specialty

	return flat
}
