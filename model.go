package main

type Game struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
	Price int    `json:"price"`
}
