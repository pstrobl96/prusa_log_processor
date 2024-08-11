package syslog

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

func startSyslogServer(listenUDP string) (syslog.LogPartsChannel, *syslog.Server) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(listenUDP)
	server.Boot()
	return channel, server
}

func getSeverity(severity int) string {
	switch severity {
	case 1:
		messageCounterDebug.Inc()
		return "Debug"
	case 2:
		messageCounterInfo.Inc()
		return "Info"
	case 3:
		messageCounterWarn.Inc()
		return "Warning"
	case 4:
		messageCounterError.Inc()
		return "Error"
	case 5:
		messageCounterCrit.Inc()
		return "Critical"
	case 6:
		messageCounterInfo.Inc()
		return "Informational"
	default:
		messageCounterUnkn.Inc()
		return "Unknown - " + strconv.Itoa(severity)
	}
}

// HandleLogs is a function to handle logs from syslog and send them to Loki or Promtail - promtail does not work because printers send logs in a different format than it should and Promtails throws EOF error
func HandleLogs(listenUDP string, directory string, filename string, maxSize int, maxBackups int, maxAge int) {
	channel, server := startSyslogServer(listenUDP)
	log.Debug().Msg("Syslog server for logs started at: " + listenUDP)

	if err := os.MkdirAll(directory, 0744); err != nil {
		log.Panic().Err(err).Str("path", directory).Msg("Can't create log directory")
		return
	}

	writers := []io.Writer{&lumberjack.Logger{
		Filename:   path.Join(directory, filename),
		MaxBackups: maxBackups, // maximum number of backups
		MaxSize:    maxSize,    // in MB
		MaxAge:     maxAge,     // in Days
	}}

	mw := io.MultiWriter(writers...)

	syslogLogger := zerolog.New(mw).With().Timestamp().Logger()

	log.Debug().Msg("Syslog logs are being written to: " + path.Join(directory, filename))

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {

			log.Trace().Msg(fmt.Sprintf("%v", logParts))
			syslogLogger.Info().
				Str("app_name", logParts["app_name"].(string)).
				Str("client", strings.Split(logParts["client"].(string), ":")[0]).
				Str("hostname", logParts["hostname"].(string)).
				Str("priority", strconv.Itoa(logParts["priority"].(int))).
				Str("proc_id", logParts["proc_id"].(string)).
				Str("msg_id", logParts["msg_id"].(string)).
				Str("severity", getSeverity(logParts["severity"].(int))).
				Str("facility", strconv.Itoa(logParts["facility"].(int))).
				Str("structured_data", logParts["structured_data"].(string)).
				Str("tls_peer", logParts["tls_peer"].(string)).
				Str("version", strconv.Itoa(logParts["version"].(int))).
				Msg(logParts["message"].(string))

		}
	}(channel)

	server.Wait()
}
