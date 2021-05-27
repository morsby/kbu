package kbu

import (
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
		URL              string
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
			want: []string{"e6921b08e6178c70a679b3914cd34b44"},
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
			want: []string{"5a1d7de97a5334df51331d3ae442723c", "a05b9fdd3ceb488c682efe9b88509508"},
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
			want: []string{"0350503fd797bab6a11b539a17e1a958", "0350503fd797bab6a11b539a17e1a958a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, sel := range tt.selections {
				s := &Selection{
					ID:          sel.ID,
					Round:       sel.Runde,
					Date:        sel.Dato,
					University:  sel.Universitet,
					Number:      sel.Nummer,
					Region:      sel.Region,
					Start:       sel.Startdato,
					Place1:      sel.Uddannelsessted,
					Department1: sel.Afdeling,
					Specialty1:  sel.Speciale,
					Place2:      sel.Uddannelsessted2,
					Department2: sel.Afdeling2,
					Specialty2:  sel.Speciale2,
					URL:         sel.URL,
				}
				if got := s.GenerateID(); got != tt.want[i] {
					t.Errorf("Selection.GenerateID() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}
