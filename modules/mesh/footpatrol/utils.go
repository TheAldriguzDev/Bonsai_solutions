package footpatrol

import (
	"math/rand"
	"time"
)

var footPatrolKeys = []string{
	"D28D0A3C97CF457D957758834BF261F7",
	"679B43969CB94EBCB2F6B67946E96951",
}

func getRandomString(stringList []string) string {
	rand.Seed(time.Now().UnixNano())

	// Generazione di un indice casuale basato sulla lunghezza della lista
	randomIndex := rand.Intn(len(stringList))

	// Restituzione della stringa corrispondente all'indice casuale
	return stringList[randomIndex]
}
