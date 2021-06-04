package kbu

import (
	"reflect"
	"testing"
	"time"
)

func TestSelection_GenerateID(t *testing.T) {
	type fields struct {
		ID               string
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
			want: []string{"3306b9a58533f1f952d65b5485467381"},
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
			want: []string{"67bfdf2e27b92f66dd14bd660575b10c", "99f4fc8891832d8ad4a010d5e052a198"},
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
			want: []string{"030af0a9c671f134df0a92d0d18604de", "030af0a9c671f134df0a92d0d18604dea"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, sel := range tt.selections {
				s := &Selection{
					ID:         sel.ID,
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
		ID         string
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
				ID:         "asd",
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
				ID:          "asd",
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
				ID:         tt.fields.ID,
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
