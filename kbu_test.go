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
					Round:       Round{Year: 2010, Season: SeasonFall},
					Date:        time.Date(2010, time.April, 9, 9, 0, 0, 0, location),
					University:  "KU",
					Number:      0,
					RelNumber:   1.0,
					Region:      RegionH,
					Start:       time.Date(2010, time.September, 1, 0, 0, 0, 0, location),
					Place1:      "Frederiksberg Hospital",
					Department1: "Medicinsk afdeling - Kardiologi/endokrinologi",
					Specialty1:  "Intern medicin",
					Place2:      "Christen Myrup",
					Department2: "Almen praksis",
					Specialty2:  "Almen medicin",
					URL:         "https://kbu.logbog.net/Ajax_get2010v2.asp",
				},
			},
		},
		{
			name: "two valid selections",
			args: args{strings.NewReader(`[
				{
					"url": "https://kbu.logbog.net/Ajax_get2015v1.asp",
					"Valgt": "10. okt 13:20",
					"Lodtr.": "1 KU",
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
					"Lodtr.": "2 AU",
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
					Round:       Round{Year: 2015, Season: SeasonSpring},
					Date:        time.Date(2014, time.October, 10, 13, 20, 0, 0, location), // "10. okt 13:20",
					University:  "KU",
					Number:      1,
					RelNumber:   0.5,
					Region:      RegionMidt,
					Start:       time.Date(2015, time.February, 1, 0, 0, 0, 0, location), //"1. feb 15",
					Place1:      "Hospitaleenheden Vest, Regionshospitalet Holstebro",
					Department1: "Medicinsk afdeling",
					Specialty1:  "Intern medicin",
					Place2:      "Almen praksis i Holstebro Kommune",
					Department2: "Almen praksis",
					Specialty2:  "Almen medicin",
					URL:         "https://kbu.logbog.net/Ajax_get2015v1.asp",
				},
				{
					Round:       Round{Year: 2015, Season: SeasonSpring},
					Date:        time.Date(2014, time.October, 10, 13, 31, 0, 0, location), //"10. okt 13:31",
					University:  "AU",
					Number:      2,
					RelNumber:   0,
					Region:      RegionNord,
					Start:       time.Date(2015, time.March, 1, 0, 0, 0, 0, location), //"1. mar 15",
					Place1:      "Aalborg Universitetshospital",
					Department1: "Klinik Kirurgi - Kræft, Kirurgi",
					Specialty1:  "Kirurgi",
					Place2:      "Almen praksis i område Nordjylland",
					Department2: "Almen praksis",
					Specialty2:  "Almen medicin",
					URL:         "https://kbu.logbog.net/Ajax_get2015v1.asp",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRawJSON(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONRaw() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}

			// override the calculated ID so tests dont fail on them
			for i := 0; i < len(tt.want); i++ {
				got[i].ID = ""
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSON() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
