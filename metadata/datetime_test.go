package metadata

import (
	"testing"
)

func TestDateTime_Parse(t *testing.T) {
	type fields struct {
		date   string
		time   string
		subsec string
		zone   string
	}
	tests := []struct {
		name    string
		arg     string
		out     fields
		wantErr bool
	}{
		{
			"empty",
			"",
			fields{},
			false,
		},
		{
			"date",
			"2021-06-21",
			fields{"2021-06-21", "00:00:00", "", ""},
			false,
		},
		{
			"datetime",
			"2021-06-21T03:55:00",
			fields{"2021-06-21", "03:55:00", "", ""},
			false,
		},
		{
			"subsec",
			"2021-06-21T03:55:00.123",
			fields{"2021-06-21", "03:55:00", "123", ""},
			false,
		},
		{
			"zone",
			"2021-06-21T03:55:00-07:00",
			fields{"2021-06-21", "03:55:00", "", "-07:00"},
			false,
		},
		{
			"all",
			"2021-06-21T03:55:00.123-07:00",
			fields{"2021-06-21", "03:55:00", "123", "-07:00"},
			false,
		},
		{
			"zoneZ",
			"2021-06-21T03:55:00Z",
			fields{"2021-06-21", "03:55:00", "", "Z"},
			false,
		},
		{
			"allZ",
			"2021-06-21T03:55:00.123Z",
			fields{"2021-06-21", "03:55:00", "123", "Z"},
			false,
		},
		{
			"zoneZero",
			"2021-06-21T03:55:00-00:00",
			fields{"2021-06-21", "03:55:00", "", "Z"},
			false,
		},
		{
			"allZero",
			"2021-06-21T03:55:00.123+00:00",
			fields{"2021-06-21", "03:55:00", "123", "Z"},
			false,
		},
		// ERROR CASES
		{
			"garbage",
			"garbage",
			fields{},
			true,
		},
		{
			"dateZ",
			"2021-06-21Z",
			fields{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &DateTime{}
			if err := dt.Parse(tt.arg); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (dt.date != tt.out.date || dt.time != tt.out.time || dt.subsec != tt.out.subsec || dt.zone != tt.out.zone) {
				t.Errorf("DateTime.Parse() result = %+v, wanted %+v", dt, tt.out)
			}
		})
	}
}

func TestDateTime_String(t *testing.T) {
	tests := []struct {
		name     string
		receiver DateTime
		want     string
	}{
		{
			"empty",
			DateTime{},
			"",
		},
		{
			"datetime",
			DateTime{"D", "T", "", ""},
			"DTT",
		},
		{
			"datetimesub",
			DateTime{"D", "T", "S", ""},
			"DTT.S",
		},
		{
			"datetimezone",
			DateTime{"D", "T", "", "Z"},
			"DTTZ",
		},
		{
			"all",
			DateTime{"D", "T", "S", "Z"},
			"DTT.SZ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.receiver.String(); got != tt.want {
				t.Errorf("DateTime.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
