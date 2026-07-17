package service

import "testing"

func TestAmbroseValleyCoordinateMapping(t *testing.T) {
	reg := NewTransformerRegistry()
	tr, err := reg.Get("AmbroseValley")
	if err != nil {
		t.Fatal(err)
	}

	// README worked example: x=-301.45, z=-355.55 → pixel (78, 890)
	px, py := tr.WorldToPixel(-301.45, -355.55)
	if px < 76 || px > 80 {
		t.Errorf("pixelX = %v, want ~78", px)
	}
	if py < 888 || py > 892 {
		t.Errorf("pixelY = %v, want ~890", py)
	}
}

func TestGrandRiftTransformerRegistered(t *testing.T) {
	reg := NewTransformerRegistry()
	for _, id := range []string{"AmbroseValley", "GrandRift", "Lockdown"} {
		if _, err := reg.Get(id); err != nil {
			t.Errorf("map %s not registered: %v", id, err)
		}
	}
}

func TestUnknownMapReturnsError(t *testing.T) {
	reg := NewTransformerRegistry()
	if _, err := reg.Get("UnknownMap"); err == nil {
		t.Error("expected error for unknown map")
	}
}
