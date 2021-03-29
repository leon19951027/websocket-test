package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	//neverdone := make(chan bool)
	//	num := make(chan int)
	DogChan := make(chan int, 1)
	CatChan := make(chan int)
	FishChan := make(chan int)

	go Dog(DogChan, CatChan)
	go Cat(CatChan, FishChan)
	go Fish(FishChan, DogChan)
	log.Println("---------")
	DogChan <- 1
	time.Sleep(10 * time.Minute)

}

func Dog(DogChan, CatChan chan int) {
	i := 0
	for {
		select {
		case <-DogChan:
			fmt.Println("Dog")
			CatChan <- 1
			i++
			if i > 100 {
				return
			}
		}
	}

}

func Cat(CatChan, FishChan chan int) {
	i := 0
	for {
		select {
		case <-CatChan:
			fmt.Println("Cat")
			FishChan <- 1
			i++
			if i > 100 {
				return
			}
		}

	}

	return
}
func Fish(FishChan, DogChan chan int) {
	i := 0
	for {
		select {
		case <-FishChan:
			fmt.Println("Fish")
			DogChan <- 1
			i++
			if i > 100 {
				return
			}
		}
	}
}
