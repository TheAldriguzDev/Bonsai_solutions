package main

import (
	"fmt"
	"os"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/api"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/adidas"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/asos"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/bootz"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/footlocker"
	ftlNewAdd "github.com/TheAldriguzDev/Bonsai_solutions/modules/footlockerNewAdd"
	ftlRestock "github.com/TheAldriguzDev/Bonsai_solutions/modules/footlockerRestock"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/manashop"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/footpatrol"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/jd"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/mesh/size"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/rinascente"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/skatedeluxe"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/sotf"
	"github.com/TheAldriguzDev/Bonsai_solutions/modules/zalando"
	zalando_kids "github.com/TheAldriguzDev/Bonsai_solutions/modules/zalando-kids"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
)

var location string = "Dashboard"

func main() {

	//go protection.GetProcess("HM", "OCTA")

	var selection int
	fmt.Printf(log.Colors.DarkGrey+"\r[%s] [%s] Select an option: \n"+log.Colors.Reset, location, log.GetTime())
	log.General("1 - Footlocker", location)
	log.General("2 - Zalando", location)
	log.General("3 - Adidas", location)
	log.General("4 - Bootz", location)
	log.General("5 - Rinascente", location)
	log.General("6 - Manashop", location)
	log.General("7 - SkateDeluxe", location)
	log.General("8 - JD", location)
	log.General("9 - JD restock", location)
	log.General("10 - Size", location)
	log.General("11 - Size Restock", location)
	log.General("12 - Footpatrol", location)
	log.General("13 - Footpatrol Restock", location)
	log.General("14 - Sotf", location)
	log.General("15 - Ftl New Add", location)
	log.General("16 - Ftl Restock", location)
	log.General("17 - Asos", location)
	log.General("999 - Start apis", location)
	//log.General("999 - AIO sites [Soon]", location)
	fmt.Printf(log.Colors.DarkGrey+"\r[%s] [%s] >>> "+log.Colors.Reset, location, log.GetTime())

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
		zalando.InitZalandoWoman()
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
		jd.InitMeshNewAdd()
		break
	case 9:
		jd.InitJdRestock()
		break
	case 10:
		size.InitMeshNewAdd()
		break
	case 11:
		size.InitJdRestock()
		break
	case 12:
		footpatrol.InitMeshNewAdd()
		break
	case 13:
		footpatrol.InitJdRestock()
		break
	case 14:
		sotf.InitSotf()
		break
	case 15:
		ftlNewAdd.InitFootlocker()
		break
	case 16:
		ftlRestock.InitRestock()
		break
	case 17:
		asos.InitAsos()
		break
	case 18:
		zalando_kids.InitZalando()
		break
	case 999:
		api.ApiServer()
	default:
		log.Error("No module found!", "Error")
		os.Exit(0)
		break
	}
}
