package tokens

import "testing"

func TestSelect(t *testing.T) {
	p := Select("original")
	if p.Accent != Original.Accent {
		t.Errorf("original: accent mismatch")
	}
	p = Select("corporate")
	if p.Accent != Corporate.Accent {
		t.Errorf("corporate: accent mismatch")
	}
	p = Select("cyber")
	if p.Accent != Cyber.Accent {
		t.Errorf("cyber: accent mismatch")
	}
	p = Select("unknown")
	if p.Accent != Original.Accent {
		t.Errorf("unknown should fallback to original")
	}
}
