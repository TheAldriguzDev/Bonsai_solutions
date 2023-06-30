package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/modules/adidas"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/footlocker"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
)

var location string = "Dashboard"

func main() {
	var selection int
	log.General("Select an option: ", location)
	log.General("1 - Footlocker", location)
	log.General("2 - Zalando", location)
	log.General("3 - Adidas", location)
	fmt.Printf(log.Colors.DarkGrey+"\r[%s] [%s] >>> "+log.Colors.Reset, location, log.GetTime())

	if _, err := fmt.Scan(&selection); err != nil || selection < 0 || selection > 4 {
		fmt.Printf(log.Colors.Red+"\r[%s] [%s] Wrong input! \n"+log.Colors.Reset, location, log.GetTime())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	switch selection {
	case 1:
		footlocker.InitFootlocker()
		break
	case 2:
		//zalando.InitZalando()
		break
	case 3:
		adidas.InitAdidas()
		break
	default:
		log.Error("No module found!", "Error")
		os.Exit(0)
		break
	}
}
