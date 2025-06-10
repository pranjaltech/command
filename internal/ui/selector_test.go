package ui

import "testing"

func TestNewModelSize(t *testing.T) {
	m := newModel([]string{"a", "b"})
	if m.list.Width() <= 0 {
		t.Errorf("expected width > 0, got %d", m.list.Width())
	}
	if m.list.Height() != 2 {
		t.Errorf("expected height 2, got %d", m.list.Height())
	}
}
