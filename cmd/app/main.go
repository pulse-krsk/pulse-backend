package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/kurochkinivan/pulskrsk/config"
	_ "github.com/kurochkinivan/pulskrsk/docs"
	"github.com/kurochkinivan/pulskrsk/internal/app"
	"github.com/sirupsen/logrus"
)

const oauthToken = "y0_AgAAAAB2JboiAAzppgAAAAEbVlN_AAD8gNKsL0tNoaKIMd-HJwnWIu7IyA"
const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	logrus.Info("loading config")
	cfg := config.MustLoad()

	logrus.Info("setup logger")
	setupLogger(cfg.Environment)

	a, err := app.NewApp(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	err = a.StartHTTP()
	if err != nil {
		logrus.Fatal(err)
	}
}

func setupLogger(env string) {
	callerPrettyfier := func(f *runtime.Frame) (string, string) {
		filename := path.Base(f.File)
		funcName := path.Base(f.Function)
		return fmt.Sprintf("%s()", funcName), fmt.Sprintf("%s:%d", filename, f.Line)
	}

	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	switch env {
	case envLocal:
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			CallerPrettyfier: callerPrettyfier,
		})
		logrus.SetLevel(logrus.TraceLevel)
	case envProd:
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			CallerPrettyfier: callerPrettyfier,
		})
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.Fatal("unknown environment")
	}
}
