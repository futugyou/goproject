package main

type collection interface {
	createiterator() iterator
}
