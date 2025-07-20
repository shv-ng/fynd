package indexer

import (
	"reflect"
	"slices"
	"testing"
)

func TestStemmer(t *testing.T) {
	input := []string{"running", "jumps", "easily", "fishing"}
	expected := []string{"run", "jump", "easili", "fish"}

	output := Stemmer(input)

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestRemoveStopWords(t *testing.T) {
	input := []string{"the", "quick", "brown", "fox", "jumps", "on", "the", "lazy", "dog"}
	expected := []string{"quick", "brown", "fox", "jumps", "lazy", "dog"}

	output := RemoveStopWords(input)

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v, got %v", expected, output)
	}
}

func TestTokenizer(t *testing.T) {
	content := "hello i run 24 ke34a)(adad)"
	expected := []string{"hello", "run", "ke34a", "adad"}
	found := Tokenise(content)
	if !slices.Equal(expected, found) {
		t.Errorf("expected: %v; found: %v", expected, found)
	}
}
