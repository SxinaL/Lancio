package test

import (
	"Lancio/services"

	"testing"
)

func TestFetchOnline(t *testing.T) {
	word := "say"
	entries, err := services.FetchOnline(word, "en")
	
	if err != nil {
		t.Errorf("FetchOnline failed: %v", err)
	}
	t.Logf("FetchOnline success: %v", entries)
}