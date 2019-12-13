package main

import (
	"environment/cfgargs"
	"fmt"
	"singledb/app"
	"singledb/server"
)

var (
	BuildVersion = ""
)

func main() {
	srvCfg, err := cfgargs.InitSrvConfig(BuildVersion, func() {
		// user flag binding code
	})
	if nil != err {
		fmt.Println(err)
		return
	}
	app.GetApp().InitApp(srvCfg)

	srv := server.NewServer()
	srv.Run(srvCfg.Info.Addr)
	// genserver.TestCodecServe(srvCfg.Info.Addr, true, false, nil)
}
