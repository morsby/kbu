package kbu

import "time"

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

type Season string

const (
	SeasonSpring Season = "forår"
	SeasonFall   Season = "efterår"
)

type Region string

const (
	RegionMidt Region = "Midt"
	RegionNord Region = "Nord"
	RegionSj   Region = "Sj"
	RegionH    Region = "H"
	RegionSyd  Region = "Syd"
)

type University string

const (
	UniversityAU  University = "AU"
	UniversityKU  University = "KU"
	UniversitySDU University = "SDU"
	UniversityAAU University = "AAU"
	UniversityNA  University = "NA"
)

type Round struct {
	Season Season `json:"season"`
	Year   int    `json:"year"`
}

type Selection struct {
	ID          string     `json:"id"`
	Round       Round      `json:"round"`
	Date        time.Time  `json:"date"`
	University  University `json:"university"`
	Number      int        `json:"no"`
	RelNumber   float64    `json:"relNumber"`
	Region      Region     `json:"region"`
	Start       time.Time  `json:"start"`
	Place1      string     `json:"place1"`
	Department1 string     `json:"department1"`
	Specialty1  string     `json:"specialty1"`
	Place2      string     `json:"place2"`
	Department2 string     `json:"department2"`
	Specialty2  string     `json:"specialty2"`
	URL         string     `json:"url"`
}
