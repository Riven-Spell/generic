package enumerable

// AB is an identity type; effectively a KeyValue pair in a dictionary, but abstractly, just two values without inherent association.
type AB[A, B any] struct {
	A A
	B B
}

// ABC is just three values with no inherent association.
type ABC[A, B, C any] struct {
	AB[A, B]
	C C
}

// ABCD for the especially daring, is just four values with no inherent association.
type ABCD[A, B, C, D any] struct {
	ABC[A, B, C]
	D D
}

// ABCDE Free will sure is a thing, huh? What a thing to do with it.
type ABCDE[A, B, C, D, E any] struct {
	ABCD[A, B, C, D]
	E E
}

// ABCDEF We're just doing things around here.
type ABCDEF[A, B, C, D, E, F any] struct {
	ABCDE[A, B, C, D, E]
	F F
}
