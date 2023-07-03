package size

import (
	"math/rand"
	"time"
)

var sizeKeys = []string{
	"78F4021E82C64DBA858F49E210192EE0",
	"FB90965B96454DC7B246D78C7358D93B",
}

func getRandomString(stringList []string) string {
	rand.Seed(time.Now().UnixNano())

	// Generazione di un indice casuale basato sulla lunghezza della lista
	randomIndex := rand.Intn(len(stringList))

	// Restituzione della stringa corrispondente all'indice casuale
	return stringList[randomIndex]
}
