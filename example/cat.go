package main

import (
	"log"
)

//go:generate optionGen

func _CatOptionDeclaration() interface{} {
	return map[string]interface{}{
		"sounds": string("Meow"),
		"food":   (*string)(nil),
		"Walk": func() {
			log.Println("Walking")
		},
	}
}

type Cat struct {
	options CatOptions
}

func NewCat(option ... CatOp) *Cat {
	cat := Cat{
		options: _NewCatOptions(),
	}

	for _, op := range option {
		op(&cat.options)
	}
	return &cat
}

func (c *Cat) Play() {
	// Yell
	log.Println(c.options.sounds)

	// Eat
	if c.options.food != nil {
		log.Println("Eating", *c.options.food)
	} else {
		log.Println("There is no food")
	}

	// Walk
	c.options.Walk()
}

func main() {
	log.SetFlags(0)

	// Normal Cat
	log.SetPrefix("Normal Cat: ")
	cat := NewCat()
	cat.Play()

	// Optional Cat
	log.SetPrefix("Optional Cat: ")
	food := "Cake"

	optionalCat := NewCat(
		CatOpWith_sounds("Purr"),
		CatOpWith_food(&food),
		CatOpWith_Walk(func() {
			log.Println("Flying")
		}),
	)
	optionalCat.Play()
}
