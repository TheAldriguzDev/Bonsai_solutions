package jd

import (
	"math/rand"
	"time"
)

var jdKeys = []string{
	"BEAA04857083438297405124CE3517D0",
	"A06043AAAD2B4AF3B009B46BE13404C2",
}

func getRandomString(stringList []string) string {
	rand.Seed(time.Now().UnixNano())

	// Generazione di un indice casuale basato sulla lunghezza della lista
	randomIndex := rand.Intn(len(stringList))

	// Restituzione della stringa corrispondente all'indice casuale
	return stringList[randomIndex]
}
