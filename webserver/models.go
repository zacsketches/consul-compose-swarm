package main

type product struct {
	Name        string
	Description string
}

func defaultProduct() product {
	return product{
		Name:        "Default Product",
		Description: "This is the default product that is really cool.",
	}
}
