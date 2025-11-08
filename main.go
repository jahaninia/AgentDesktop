package main

import (
	"log"
	"sync"

	"github.com/getlantern/systray"
	"jahaninia.ir/agentDesktop/jolConfigurtion"
	"jahaninia.ir/agentDesktop/jolSystry"
)

var err error

func main1() {

	configSetting, err := jolConfigurtion.LoadConfigFile()
	if err != nil {
		log.Panic(err)
	}
	app := &jolConfigurtion.App{
		Config: configSetting,
		Wg:     &sync.WaitGroup{},
	}
	// // تنظیم لاگ
	// logFile, err := os.OpenFile("sse_tray.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	log.SetOutput(logFile)
	// 	defer logFile.Close()
	// }

	log.Println("برنامه در حال شروع...")
	app.Wg.Add(1)

	systray.Run(func() { jolSystry.OnReady(app) }, jolSystry.OnExit)
	app.Wg.Wait()
}
