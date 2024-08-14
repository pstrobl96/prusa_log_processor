package cmd

import (
	"net/http"
	"strconv"

	"pstrobl96/prusa_log_processor/syslog"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	metricsPath         = kingpin.Flag("exporter.metrics-path", "Path where to expose metrics.").Default("/metrics").String()
	metricsPort         = kingpin.Flag("exporter.metrics-port", "Port where to expose metrics.").Default("10010").Int()
	syslogListenAddress = kingpin.Flag("processor.address", "Address where to expose port for gathering logs.").Default("0.0.0.0:8514").String()
	syslogToFile        = kingpin.Flag("processor.log-to-file", "Write logs to file.").Default("false").Bool()
	syslogDirectory     = kingpin.Flag("processor.directory", "Directory where to store logs.").Default("./logs").String()
	syslogFilename      = kingpin.Flag("processor.filename", "Filename where to store logs.").Default("prusa.log").String()
	syslogMaxSize       = kingpin.Flag("processor.max-size", "Maximum size of log file.").Default("10").Int()
	syslogMaxBackups    = kingpin.Flag("processor.max-backups", "Maximum number of backups.").Default("3").Int()
	syslogMaxAge        = kingpin.Flag("processor.max-age", "Maximum age of log file.").Default("28").Int()
	logLevel            = kingpin.Flag("log.level", "Log level for prusa_log_processor.").Default("info").String()
)

// Run function to start the exporter
func Run() {
	kingpin.Parse()
	log.Info().Msg("Prusa log processor starting")

	logLevel, err := zerolog.ParseLevel(*logLevel)

	if err != nil {
		logLevel = zerolog.InfoLevel // default log level
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	//var collectors []prometheus.Collector

	log.Info().Msg("Syslog logs server starting at: " + *syslogListenAddress)
	go syslog.HandleLogs(*syslogListenAddress, *syslogDirectory, *syslogFilename, *syslogMaxSize, *syslogMaxBackups, *syslogMaxAge, *syslogToFile)

	http.Handle(*metricsPath, promhttp.Handler())
	http.ListenAndServe(":"+strconv.Itoa(*metricsPort), nil)

}
