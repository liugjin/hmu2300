package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unicode"

	"clc.hmu/app/frp/frpc"
	"clc.hmu/app/public/store/muid"
	"clc.hmu/app/supd/config"
	"clc.hmu/app/supd/logger"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

func RunTools(cfg *config.Config) {
	supdCfg, ok := cfg.GetSupervisord()
	if !ok {
		panic("need supervisord setion in supd.ini")
	}
	// work for frp
	go frpc.StartSSH(&frpc.FrpcSSH{
		ServerHost:  frpc.DefaultServerAddr,
		ServerPort:  frpc.DefaultServerPort,
		ServerToken: frpc.DefaultServerToken,
		MuID:        muid.GetMuID(),
		MuSK:        frpc.DefaultSSHSK,
		SSHPort:     supdCfg.GetString("local_port", "22"),
	})
}

type Options struct {
	Configuration string `short:"c" long:"configuration" description:"the configuration file"`
	Daemon        bool   `short:"d" long:"daemon" description:"run as daemon"`
	EnvFile       string `long:"env-file" description:"the environment file"`
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logger.CustomFormatter{FullTimestamp: true})
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func initSignals(s *Supervisor) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.WithFields(log.Fields{"signal": sig}).Info("receive a signal to stop all process & exit")
		s.procMgr.StopAllProcesses()
		os.Exit(-1)
	}()

}

var options Options

var parser = flags.NewParser(&options, flags.Default & ^flags.PrintErrors)

func LoadEnvFile() {
	if len(options.EnvFile) <= 0 {
		return
	}
	//try to open the environment file
	f, err := os.Open(options.EnvFile)
	if err != nil {
		log.WithFields(log.Fields{"file": options.EnvFile}).Error("Fail to open environment file")
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		//for each line
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		//if line starts with '#', it is a comment line, ignore it
		line = strings.TrimSpace(line)
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		//if environment variable is exported with "export"
		if strings.HasPrefix(line, "export") && len(line) > len("export") && unicode.IsSpace(rune(line[len("export")])) {
			line = strings.TrimSpace(line[len("export"):])
		}
		//split the environment variable with "="
		pos := strings.Index(line, "=")
		if pos != -1 {
			k := strings.TrimSpace(line[0:pos])
			v := strings.TrimSpace(line[pos+1:])
			//if key and value are not empty, put it into the environment
			if len(k) > 0 && len(v) > 0 {
				os.Setenv(k, v)
			}
		}
	}
}

// find the supervisord.conf in following order:
//
// 1. $CWD/supervisord.conf
// 2. $CWD/etc/supervisord.conf
// 3. /etc/supervisord.conf
// 4. /etc/supervisor/supervisord.conf (since Supervisor 3.3.0)
// 5. ../etc/supervisord.conf (Relative to the executable)
// 6. ../supervisord.conf (Relative to the executable)
func findSupervisordConf() (string, error) {
	possibleSupervisordConf := []string{options.Configuration,
		// for clc.hmu
		os.ExpandEnv("$PRJ_ROOT/app/supd/etc/supd.ini"),
		// keep for old
		"./supd.ini",
		"./etc/supd.ini",
		"/etc/supd.ini",
		"/etc/supd/supd.ini",
		"../etc/supd.ini",
		"../supd.ini",
	}

	for _, file := range possibleSupervisordConf {
		if _, err := os.Stat(file); err == nil {
			abs_file, err := filepath.Abs(file)
			if err == nil {
				return abs_file, nil
			} else {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf("fail to find supervisord.conf")
}

func RunServer() {
	// infinite loop for handling Restart ('reload' command)
	LoadEnvFile()
	for true {
		options.Configuration, _ = findSupervisordConf()
		s := NewSupervisor(options.Configuration)
		initSignals(s)
		if sErr, _, _, _ := s.Reload(); sErr != nil {
			panic(sErr)
		}
		RunTools(s.GetConfig())

		s.WaitForExit()
	}
}

func main() {

	// for supervisor
	ReapZombie()

	if _, err := parser.Parse(); err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok {
			switch flagsErr.Type {
			case flags.ErrHelp:
				fmt.Fprintln(os.Stdout, err)
				os.Exit(0)
			case flags.ErrCommandRequired:
				// work for supervisor
				if options.Daemon {
					Deamonize(RunServer)
				} else {
					RunServer()
				}
			default:
				fmt.Fprintf(os.Stderr, "error when parsing command: %s\n", err)
				os.Exit(1)
			}
		}
	}
}
