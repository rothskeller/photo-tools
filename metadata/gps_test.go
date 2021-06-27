package metadata

import "testing"

func TestGPS(t *testing.T) {
	var gc, gc2 GPSCoords
	if err := gc.Parse("37.33544, -122.01990, 200ft"); err != nil {
		t.Errorf("gc.Parse failed")
	}
	{
		latref, lat, longref, long, altref, alt := gc.AsEXIF()
		if err := gc2.ParseEXIF(latref, lat, longref, long, altref, alt); err != nil {
			t.Errorf("gc.ParseEXIF failed")
		}
	}
	{
		lat, long, altref, alt := gc2.AsXMP()
		if err := gc2.ParseXMP(lat, long, altref, alt); err != nil {
			t.Errorf("gc.ParseXMP failed")
		}
	}
	if !gc.Equivalent(&gc2) {
		t.Errorf("result is not equivalent: %s", gc2.String())
	}
	if gc2.String() != "37.33544, -122.0199, 200ft" {
		t.Errorf("result is wrong: %s", gc2.String())
	}
}
