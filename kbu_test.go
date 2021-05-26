package kbu

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

var location *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		panic(err)
	}
	location = loc
}

func TestParseJSONRaw(t *testing.T) {

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Selection
		wantErr bool
	}{
		{
			name: "Valid selection",
			args: args{strings.NewReader(`[
				{
					"url":"https://kbu.logbog.net/Ajax_get2010v2.asp",
					"Valgt":"9. apr 09:00",
					"Lodtr.":" KU",
					"Region":"Hovedst.",
					"Startdato":"1. sep 10",
					"Uddannelsessted":"Frederiksberg Hospital",
					"Afdeling":"Medicinsk afdeling - Kardiologi/endokrinologi",
					"Speciale":"Intern medicin",
					"Uddannelsessted2":"Christen Myrup",
					"Afdeling2":"Almen praksis",
					"Speciale2":"Almen medicin"
				}
			]`)},
			want: []Selection{
				{
					Runde:            Round{Year: 2010, Season: SeasonFall},
					Dato:             time.Date(2010, time.April, 9, 9, 0, 0, 0, location),
					Universitet:      "KU",
					Nummer:           0,
					Region:           "Hovedst.",
					Startdato:        time.Date(2010, time.September, 1, 0, 0, 0, 0, location),
					Uddannelsessted:  "Frederiksberg Hospital",
					Afdeling:         "Medicinsk afdeling - Kardiologi/endokrinologi",
					Speciale:         "Intern medicin",
					Uddannelsessted2: "Christen Myrup",
					Afdeling2:        "Almen praksis",
					Speciale2:        "Almen medicin",
					URL:              "https://kbu.logbog.net/Ajax_get2010v2.asp",
				},
			},
		},
		{
			name: "two valid selections",
			args: args{strings.NewReader(`[
				{
					"url": "https://kbu.logbog.net/Ajax_get2015v1.asp",
					"Valgt": "10. okt 13:20",
					"Lodtr.": "316 KU",
					"Region": "Midt.",
					"Startdato": "1. feb 15",
					"Uddannelsessted": "Hospitaleenheden Vest, Regionshospitalet Holstebro",
					"Afdeling": "Medicinsk afdeling",
					"Speciale": "Intern medicin",
					"Uddannelsessted2": "Almen praksis i Holstebro Kommune",
					"Afdeling2": "Almen praksis",
					"Speciale2": "Almen medicin"
				  },
				  {
					"url": "https://kbu.logbog.net/Ajax_get2015v1.asp",
					"Valgt": "10. okt 13:31",
					"Lodtr.": "317 AU",
					"Region": "Nord.",
					"Startdato": "1. mar 15",
					"Uddannelsessted": "Aalborg Universitetshospital",
					"Afdeling": "Klinik Kirurgi - Kræft, Kirurgi",
					"Speciale": "Kirurgi",
					"Uddannelsessted2": "Almen praksis i område Nordjylland",
					"Afdeling2": "Almen praksis",
					"Speciale2": "Almen medicin"
				  }
			]`)},
			want: []Selection{
				{
					Runde:            Round{Year: 2015, Season: SeasonSpring},
					Dato:             time.Date(2014, time.October, 10, 13, 20, 0, 0, location), // "10. okt 13:20",
					Universitet:      "KU",
					Nummer:           316,
					Region:           "Midt.",
					Startdato:        time.Date(2015, time.February, 1, 0, 0, 0, 0, location), //"1. feb 15",
					Uddannelsessted:  "Hospitaleenheden Vest, Regionshospitalet Holstebro",
					Afdeling:         "Medicinsk afdeling",
					Speciale:         "Intern medicin",
					Uddannelsessted2: "Almen praksis i Holstebro Kommune",
					Afdeling2:        "Almen praksis",
					Speciale2:        "Almen medicin",
					URL:              "https://kbu.logbog.net/Ajax_get2015v1.asp",
				},
				{
					Runde:            Round{Year: 2015, Season: SeasonSpring},
					Dato:             time.Date(2014, time.October, 10, 13, 31, 0, 0, location), //"10. okt 13:31",
					Universitet:      "AU",
					Nummer:           317,
					Region:           "Nord.",
					Startdato:        time.Date(2015, time.March, 1, 0, 0, 0, 0, location), //"1. mar 15",
					Uddannelsessted:  "Aalborg Universitetshospital",
					Afdeling:         "Klinik Kirurgi - Kræft, Kirurgi",
					Speciale:         "Kirurgi",
					Uddannelsessted2: "Almen praksis i område Nordjylland",
					Afdeling2:        "Almen praksis",
					Speciale2:        "Almen medicin",
					URL:              "https://kbu.logbog.net/Ajax_get2015v1.asp",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRawJSON(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
