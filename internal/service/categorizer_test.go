package service

import "testing"

func TestCategorizerMatchesDictionaryEntries(t *testing.T) {
	categorizer, err := NewCategorizer("../../data/products.json")
	if err != nil {
		t.Fatalf("NewCategorizer() error = %v", err)
	}

	tests := map[string]string{
		"banana":            "Fruit",
		"fresh banana":      "Fruit",
		"organic milk":      "Dairy",
		"whole grain bread": "Bakery",
		"bananna":           "Fruit",
		"milkk":             "Dairy",
		"mystery item":      "Other",
	}

	for name, want := range tests {
		if got := categorizer.Categorize(name); got != want {
			t.Errorf("Categorize(%q) = %q, want %q", name, got, want)
		}
	}
}
