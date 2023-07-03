package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/modules/adidas"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/bootz"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/footlocker"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/manashop"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/footpatrol"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/jd"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/size"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/rinascente"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/skatedeluxe"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/zalando"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
)

var location string = "Dashboard"

func main() {

	//go protection.GetProcess("HM", "OCTA")

	var selection int
	log.General("Select an option: ", location)
	log.General("1 - Footlocker", location)
	log.General("2 - Zalando", location)
	log.General("3 - Adidas", location)
	log.General("4 - Bootz", location)
	log.General("5 - Rinascente", location)
	log.General("6 - Manashop", location)
	log.General("7 - SkateDeluxe", location)
	log.General("8 - JD", location)
	log.General("9 - Size", location)
	log.General("10 - Footpatrol", location)
	//log.General("999 - AIO sites [Soon]", location)
	log.General(">>> ", location)

	if _, err := fmt.Scan(&selection); err != nil || selection < 0 || selection > 1000 {
		fmt.Printf(log.Colors.Red+"\r[%s] [%s] Wrong input! \n"+log.Colors.Reset, location, log.GetTime())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	switch selection {
	case 1:
		footlocker.InitFootlocker()
		break
	case 2:
		zalando.InitZalando()
		break
	case 3:
		adidas.InitAdidas()
		break
	case 4:
		bootz.InitBootz()
		break
	case 5:
		rinascente.InitRina()
		break
	case 6:
		manashop.InitMana()
		break
	case 7:
		skatedeluxe.InitSkate()
		break
	case 8:
		jd.InitMesh()
		break
	case 9:
		size.InitSize()
		break
	case 10:
		footpatrol.InitFootpatrol()
		break
	case 999:
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			rinascente.InitRina()
		}()

		go func() {
			defer wg.Done()
			bootz.InitBootz()
		}()

		wg.Wait()
	default:
		log.Error("No module found!", "Error")
		os.Exit(0)
		break
	}
}
