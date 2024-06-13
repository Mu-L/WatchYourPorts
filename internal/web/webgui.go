package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aceberg/WatchYourPorts/internal/check"
	"github.com/aceberg/WatchYourPorts/internal/conf"
	"github.com/aceberg/WatchYourPorts/internal/models"
)

// Gui - start web server
func Gui(dirPath, nodePath string) {

	confPath := dirPath + "/config.yaml"
	check.Path(confPath)

	appConfig = conf.Get(confPath)

	appConfig.DirPath = dirPath
	appConfig.DBPath = dirPath + "/sqlite.db"
	check.Path(appConfig.DBPath)
	appConfig.ConfPath = confPath
	appConfig.NodePath = nodePath

	// log.Println("INFO: starting web gui with config", appConfig.ConfPath)
	log.Println("INFO: starting web gui with config", appConfig)

	// db.Create(appConfig.DBPath)
	tmpFill()

	address := appConfig.Host + ":" + appConfig.Port

	log.Println("=================================== ")
	log.Printf("Web GUI at http://%s", address)
	log.Println("=================================== ")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	templ := template.Must(template.New("").ParseFS(templFS, "templates/*"))
	router.SetHTMLTemplate(templ) // templates

	router.StaticFS("/fs/", http.FS(pubFS)) // public

	router.GET("/", indexHandler)         // index.go
	router.GET("/config/", configHandler) // config.go
	router.GET("/scan/", scanHandler)     // scanpage.go

	router.POST("/config/", saveConfigHandler)    // config.go
	router.POST("/scan_ports/", scanPortsHandler) // scanpage.go

	err := router.Run(address)
	check.IfError(err)
}

func tmpFill() {
	var oneAddr models.AddrToScan

	oneAddr.PortMap = make(map[int]models.PortItem)
	allAddrs = make(map[string]models.AddrToScan)

	oneAddr.Name = "Onslaught"
	oneAddr.Addr = "192.168.2.3"

	allAddrs[oneAddr.Addr] = oneAddr

	oneAddr.Name = "BlastOff"
	oneAddr.Addr = "192.168.2.2"

	allAddrs[oneAddr.Addr] = oneAddr
}
