package kbu

import (
	"reflect"
	"testing"
	"time"
)

func TestSelection_GenerateID(t *testing.T) {
	type fields struct {
		Md5              string
		Runde            Round
		Dato             time.Time
		Universitet      University
		Nummer           int
		Region           Region
		Startdato        time.Time
		Uddannelsessted  string
		Afdeling         string
		Speciale         string
		Uddannelsessted2 string
		Afdeling2        string
		Speciale2        string
	}
	tests := []struct {
		name       string
		selections []fields
		want       []string
	}{
		{
			name: "one selection",
			selections: []fields{
				{
					Runde:  Round{Year: 2020, Season: SeasonFall},
					Dato:   time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
					Nummer: 123,
				},
			},
			want: []string{"5904c63ef91b4b1b141f3e53e491f329"},
		},
		{
			name: "two selections, different",
			selections: []fields{
				{
					Runde:     Round{Year: 2020, Season: SeasonFall},
					Dato:      time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
					Nummer:    999,
					Startdato: time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
				},
				{
					Runde:     Round{Year: 2020, Season: SeasonFall},
					Dato:      time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
					Nummer:    123,
					Startdato: time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
				},
			},
			want: []string{"2dbaf085f830790d1a1aefb34bfb7df0", "bae4df2f765246b7aeb2fb1f9e1b60fc"},
		},
		{
			name: "two selections, the same; IDs should be unique",
			selections: []fields{
				{
					Runde:     Round{Year: 2020, Season: SeasonFall},
					Nummer:    143,
					Startdato: time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
				},
				{
					Runde:     Round{Year: 2020, Season: SeasonFall},
					Nummer:    143,
					Startdato: time.Date(2020, time.January, 11, 0, 0, 0, 0, location),
				},
			},
			want: []string{"d076d36c26a1b9d2b1d33743eadec41b", "d076d36c26a1b9d2b1d33743eadec41ba"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, sel := range tt.selections {
				s := &Selection{
					Md5:        sel.Md5,
					Round:      sel.Runde,
					Date:       sel.Dato,
					University: sel.Universitet,
					Number:     sel.Nummer,
					Region:     sel.Region,
					Start:      sel.Startdato,
					Positions: []Position{
						{Location: sel.Uddannelsessted,
							Department: sel.Afdeling,
							Specialty:  sel.Speciale},
						{Location: sel.Uddannelsessted2,
							Department: sel.Afdeling2,
							Specialty:  sel.Speciale2},
					},
				}
				if got := s.GenerateID(); got != tt.want[i] {
					t.Errorf("Selection.GenerateID() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func TestSelection_Flatten(t *testing.T) {
	type fields struct {
		Md5        string
		Round      Round
		Date       time.Time
		University University
		Number     int
		RelNumber  float64
		Region     Region
		Start      time.Time
		Positions  []Position
	}
	tests := []struct {
		name   string
		fields fields
		want   SelectionFlat
	}{
		{
			name: "Flattening a single Selection",
			fields: fields{
				Md5:        "asd",
				Round:      Round{Year: 2020, Season: SeasonFall},
				Date:       time.Date(2020, 03, 9, 0, 0, 0, 0, location),
				University: UniversityAU,
				Number:     14,
				RelNumber:  0.12,
				Region:     RegionMidt,
				Start:      time.Date(2020, 8, 1, 0, 0, 0, 0, location),
				Positions: []Position{
					{
						Location:   "Herning",
						Department: "Akutafdelingen",
						Specialty:  "Akutmedicin",
					},
					{
						Location:   "Vildbjerg",
						Department: "Almen praksis",
						Specialty:  "Almen medicin",
					},
				},
			},
			want: SelectionFlat{
				Md5:         "asd",
				RoundYear:   2020,
				RoundSeason: SeasonFall,
				Date:        time.Date(2020, 03, 9, 0, 0, 0, 0, location),
				University:  UniversityAU,
				Number:      14,
				RelNumber:   0.12,
				Region:      RegionMidt,
				Start:       time.Date(2020, 8, 1, 0, 0, 0, 0, location),
				Location1:   "Herning",
				Department1: "Akutafdelingen",
				Specialty1:  "Akutmedicin",
				Location2:   "Vildbjerg",
				Department2: "Almen praksis",
				Specialty2:  "Almen medicin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Selection{
				Md5:        tt.fields.Md5,
				Round:      tt.fields.Round,
				Date:       tt.fields.Date,
				University: tt.fields.University,
				Number:     tt.fields.Number,
				RelNumber:  tt.fields.RelNumber,
				Region:     tt.fields.Region,
				Start:      tt.fields.Start,
				Positions:  tt.fields.Positions,
			}
			if got := s.Flatten(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Selection.Flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}
