package kbu

import (
	"reflect"
	"testing"
	"time"
)

func Test_calculateRound(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    Round
		wantErr bool
	}{
		{
			name: "Correct format",
			args: args{"https://kbu.logbog.net/Ajax_get2010v2.asp"},
			want: Round{Year: 2010, Season: SeasonFall},
		},
		{
			name: "Correct format",
			args: args{"https://kbu.logbog.net/Ajax_get2020v1.asp"},
			want: Round{Year: 2020, Season: SeasonSpring},
		},
		{
			name:    "Invalid format",
			args:    args{"https://kbu.logbog.net/Ajax_get2010v.asp"},
			want:    Round{},
			wantErr: true,
		},
		{
			name: "Latest round",
			args: args{"https://kbu.logbog.net/AJAX_Timelines.asp"},
			want: Round{Year: 2020, Season: SeasonFall},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateRound(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateRound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateRound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateNumberUni(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		number  int
		uni     string
		wantErr bool
	}{
		{
			name:   "should pass, regular space, no number",
			args:   args{" KU"},
			number: 0,
			uni:    "KU",
		},
		{
			name:   "should pass, regular space, with number",
			args:   args{"123 KU"},
			number: 123,
			uni:    "KU",
		},
		{
			name: "should pass, tab",
			args: args{"3	KU"},
			number: 3,
			uni:    "KU",
		},
		{
			name:   "should pass, weird space",
			args:   args{" KU"},
			number: 0,
			uni:    "KU",
		},
		{
			name:   "should work, no space",
			args:   args{"1KU"},
			number: 1,
			uni:    "KU",
		},
		{
			name:    "should fail, two words",
			args:    args{"AU KU"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := calculateNumberUni(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateNumberUni() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.number {
				t.Errorf("calculateNumberUni() got number = %v, want %v", got, tt.number)
			}
			if got1 != tt.uni {
				t.Errorf("calculateNumberUni() got uni = %v, want %v", got1, tt.uni)
			}
		})
	}
}

func Test_calculateDate(t *testing.T) {
	type args struct {
		str          string
		yearOverride int
	}
	tests := []struct {
		name     string
		args     args
		wantDate time.Time
		wantErr  bool
	}{
		{
			name:     "should work, no year",
			args:     args{str: "9. apr 15", yearOverride: 0},
			wantDate: time.Date(2015, time.April, 9, 0, 0, 0, 0, location),
		},
		{
			name:     "should work, with year",
			args:     args{str: "9. sep 13:24", yearOverride: 2019},
			wantDate: time.Date(2019, time.September, 9, 13, 24, 0, 0, location),
		},
		{
			name:     "should work, has both time of day and year",
			args:     args{str: "9. apr 05 13:24"},
			wantDate: time.Date(2005, time.April, 9, 13, 24, 0, 0, location),
		},
		{
			name:     "should work, but failed previously",
			args:     args{str: "30. mar. 16:00"},
			wantDate: time.Date(0, time.March, 30, 16, 0, 0, 0, location),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, err := calculateDate(tt.args.str, tt.args.yearOverride)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDate, tt.wantDate) {
				t.Errorf("calculateDate() = %v, want %v", gotDate, tt.wantDate)
			}
		})
	}
}
