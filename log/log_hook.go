package log

import (
	"net"
	"os"
	"path"
	"time"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/sohlich/elogrus.v2"
)

func lfsHook(cfg *hookConfig) logrus.Hook {
	logPathTmp := path.Join(cfg._logPath, cfg._server)
	if _, err := os.Stat(logPathTmp); os.IsNotExist(err) {
		os.MkdirAll(logPathTmp, 0644)
	}

	writer, err := rotatelogs.New(
		path.Join(logPathTmp, cfg._server+".log.%Y%m%d%H%M"),
		rotatelogs.WithMaxAge(168*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	})
}

func elasticHook(cfg *hookConfig) logrus.Hook {
	client, err := elastic.NewClient(elastic.SetURL(cfg._address))
	if err != nil {
		logrus.Error(err)
		return nil
	}
	result, err := elogrus.NewElasticHook(client, cfg._host, logrus.DebugLevel, "pear-api-log")
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return result
}

func logstashHook(cfg *hookConfig) logrus.Hook {
	conn, err := net.Dial("tcp", cfg._address)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "pear-api-log"}))
}
