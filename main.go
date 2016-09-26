package main

import (
	"flag"
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/docker/go-plugins-helpers/graphdriver"
	"github.com/wrouesnel/go.log"

	"github.com/wrouesnel/docker-bittorrent-graph-driver/graphdriver/bittorrent_overlay2"
)

const (
	PluginName              string = "bittorrent-overlay2"
	DefaultDockerPluginPath        = "/run/docker/plugins"
)

var Version string = "dev"

func main() {
	dockerPluginPath := kingpin.Flag("docker-net-plugins", "Listen path for the plugin.").Default(fmt.Sprintf("unix://%s/%s.sock", DefaultDockerPluginPath, PluginName)).URL()
	loglevel := kingpin.Flag("log-level", "Logging Level").Default("info").String()
	logformat := kingpin.Flag("log-format", "If set use a syslog logger or JSON logging. Example: logger:syslog?appname=bob&local=7 or logger:stdout?json=true. Defaults to stderr.").Default("stderr").String()
	kingpin.Parse()

	flag.Set("log.level", *loglevel)
	flag.Set("log.format", *logformat)

	log.Infoln("Docker Plugin Path:", *dockerPluginPath)
	driver, err := bittorrent_overlay2.NewBitTorrentOverlay2GraphDriver()
	if err != nil {
		log.Errorln("Failed to initialize graph driver:", err)
	}
	log.Infoln("Graph driver initialized")

	handler := graphdriver.NewHandler(driver)
	handler.ServeUnix("root", *dockerPluginPath)
}
