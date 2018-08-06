package main
import (
	"context"
	"flag"
	"fmt"
	"gitlab.com/parallelcoin/duo/pkg/server/state"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "gitlab.com/parallelcoin/duo/pkg/walletdat"
	_ "gitlab.com/parallelcoin/duo/pkg/cmd"
	"gitlab.com/parallelcoin/duo/pkg/iniflags"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/server"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/subcmd"
	"gitlab.com/parallelcoin/duo/pkg/version"
)
var (
	Version, GitCommit, GitBranch, GitState, BuildDate string
)
const ()
// PrintGitInfo prints out information from git about the current build
// Use govvv to get this to output correctly
func PrintGitInfo() {
	fmt.Printf("    Repo URL  https://%s\n"+
		"  Build date  %s\n"+
		"      Branch  %s\n"+
		"      Commit  %s\n"+
		"       State  %s\n\n",
		Version, BuildDate, GitBranch, GitCommit, GitState)
}
func printerr(format string, vars ...interface{}) {
	fmt.Fprintf(os.Stderr, format, vars...)
}
func printversioninfo() {
	version.Print()
	PrintGitInfo()
}
func createconf() {
	var confFile *os.File
	fmt.Println(*args.DataDir)
	log.Println("Creating data dir at", *args.DataDir)
	os.Mkdir(*args.DataDir, os.ModePerm)
	log.Println("Creating file", *args.Conf)
	confFile, err := os.Create(*args.Conf)
	if err != nil {
		log.Fatal(err)
	}
	defer confFile.Close()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name != "config" && f.Name != "dumpflags" && f.Name != "configure" {
			fmt.Fprintf(confFile, "%s = %s  # %s\n", f.Name, iniflags.QuoteValue(f.Value.String()), iniflags.EscapeUsage(f.Usage))
		}
	})
}
func startServer() {
	if *args.Debug || *args.DebugNet {
		file, err := os.OpenFile(*args.DataDir+"/debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err)
		} else {
			logger.Output = file
			switch {
			case *args.DebugNet:
				logger.DebugNetFlag = true
				fallthrough
			case *args.Debug:
				logger.DebugFlag = true
			}
			defer file.Close()
		}
	}
	server.Start()
}
func cleanup() {
	// Shut everything down
	fmt.Println("")
	server.Shutdown()
}
func main() {
	state.Init()
	if _, err := os.Stat(*args.Conf); os.IsNotExist(err) {
		createconf()
	}
	iniflags.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
	switch {
	case *args.Version:
		printversioninfo()
		os.Exit(1)
	case *args.Help:
		version.Print()
		fmt.Printf("%s [-options] [command] [command args...]\n", os.Args[0])
		fmt.Printf("Usage:\n\n")
		flag.Usage()
		os.Exit(0)
	case *args.TestNet && *args.RegTest:
		printerr("Error: testnet and regtest cannot be set at the same time\n")
		os.Exit(1)
	case *args.CreateConf:
		fmt.Println(*args.CreateConf)
		createconf()
	}
	ctx := context.Background()
	status := subcmd.Execute(ctx)
	if status == subcmd.ExitQuit || status == subcmd.ExitFailure {
		os.Exit(0)
	}
	if len(os.Args) == 1 {
		startServer()
	}
	os.Exit(int(status))
}
