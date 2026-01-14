package main

import
	(	
		"testing"
		"reflect"
	)

func TestInput(t *testing.T) {

	type test struct {
		input string
		want []string
	}

	tests := []test{
		{input: "    hello  world", want: []string{"hello", "world"}},
		{input: "   bulbasaur  ", want: []string{"bulbasaur"}},
		{input: "  ", want: []string{}},
		{input: "Charmander Bulbasaur Squirtle Pikachu", want: []string{"charmander", "bulbasaur", "squirtle", "pikachu"}},
	}

	for _, testCase := range tests {
		got := cleanInput(testCase.input)
		if !reflect.DeepEqual(testCase.want, got) {
			t.Fatalf("expected %v, got: %v", testCase.want, got)
		}
	}
}