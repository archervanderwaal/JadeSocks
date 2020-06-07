package main

import (
	"flag"
	"fmt"
	"github.com/archervanderwaal/JadeSocks/config"
	"github.com/archervanderwaal/JadeSocks/logger"
	"github.com/archervanderwaal/JadeSocks/socks5"
	"github.com/archervanderwaal/JadeSocks/utils"
	"github.com/aybabtme/rgbterm"
	"os"
)

const (
	Version = "1.0"
	Usage   = "Usage of JadeSocks: JadeSocks <options>"
	Logo    = `
      _           _       _____            _        
     | |         | |     / ____|          | |       
     | | __ _  __| | ___| (___   ___   ___| | _____ 
 _   | |/ _| |/ _| |/ _ \\___ \ / _ \ / __| |/ / __|
| |__| | (_| | (_| |  __/____) | (_) | (__|   <\__ \
 \____/ \__,_|\__,_|\___|_____/ \___/ \___|_|\_\___/
	:: 一北 ::					 			(1.0)
`
)

var (
	v bool
	h bool
	s bool
	c bool
)

func init() {
	flag.BoolVar(&h, "h", false, "Show usage of JadeSocks and exit")
	flag.BoolVar(&v, "v", false, "Show version of JadeSocks and exit")
	flag.BoolVar(&s, "s", false, "Use the server mode")
	flag.BoolVar(&c, "c", false, "Use the client mode")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	_, args := utils.ParseArgs(os.Args)
	if v {
		showVersion()
		return
	}
	if h || len(args) == 0 {
		flag.Usage()
		return
	}
	if s && c {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf(rgbterm.FgString("The s and client cannot be started "+
			"at the same time", 255, 0, 0)))
		return
	}
	conf := config.LoadConfig(s)
	if s {
		startServerMode(conf)
		return
	}
	if c {
		//startClientMode(conf)
		return
	}
}

func startServerMode(config *config.Config) {
	accounts := socks5.Accounts{MemoryUser: config.Users}
	serverConf := &socks5.ServerConfig{
		AuthMethods: []socks5.Authenticator{socks5.UserPassAuthenticator{Accounts: accounts}},
		Resolver:    socks5.DNSResolver{},
		ListenAddr:  config.ListenAddr,
		Network:     "tcp",
		Logger:      logger.Logger,
	}
	serve, err := socks5.New(serverConf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, rgbterm.FgString("Internal error: "+err.Error(), 255, 0, 0))
		return
	}
	if err := serve.ListenAndServe(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, rgbterm.FgString("Internal error: "+err.Error(), 255, 0, 0))
		return
	}
}

//func startClientMode(config *config.ServerConfig) {
//	cli, err := client.NewClient(config)
//	logger.Logger.Infof("Create client success %s\n", cli)
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, rgbterm.FgString("Internal error: "+err.Error(), 255, 0, 0))
//		return
//	}
//	err = cli.Listen(func(listenerAddr *net.TCPAddr) {
//		log.Println(rgbterm.FgString("Start client success", 0, 255, 0))
//	})
//	if err != nil {
//		_, _ = fmt.Fprintf(os.Stderr, rgbterm.FgString("Internal error: "+err.Error(), 255, 0, 0))
//		return
//	}
//}

func usage() {
	// #00FF00
	logo := rgbterm.FgString(Logo, 0, 255, 0)
	// #FF42E1
	usage := rgbterm.FgString(Usage, 255, 66, 225)
	_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n\n%s\n", logo, usage))
	flag.PrintDefaults()
	os.Exit(0)
}

func showVersion() {
	version := rgbterm.FgString(Version, 0, 255, 0)
	fmt.Println(version)
	os.Exit(0)
}
