package clients

import (
	"log"
	"os"
)

//ProjectID is your google cloud proyect ID
var ProjectID = ""

//dev feature production
func init() {
	gitFlow := os.Getenv("gitFlow") //what branch your are working
	log.Println("GitFlow: ", gitFlow)
	switch gitFlow {
	case "master":
		ProjectID = "proyectmaster" //proyectID for master branch
	case "dev":
		ProjectID = "proyectdev" //proyectID for dev branch
	case "qa":
		ProjectID = "proyectqa" //proyectID for qa branch
	case "feature":
		ProjectID = "proyectfeature" //proyectID for feature branch
	default:
		ProjectID = "proyectdev" //proyectID for default branch(dev) branch
	}
	log.Println("ProjectID: ", ProjectID)
}
