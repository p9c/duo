// Duo is a new client for the Parallelcoin network written in Golang and using a html/json/dom/ajax progressive web application for its GUI interface
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/awnumar/memguard"
	_ "gitlab.com/parallelcoin/duo/pkg/cmds"
	"gitlab.com/parallelcoin/duo/pkg/iniflags"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/server"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/server/state"
	"gitlab.com/parallelcoin/duo/pkg/subcmd"
	"gitlab.com/parallelcoin/duo/pkg/version"
	"log"
	"os"
)

var (
	Version, GitCommit, GitBranch, GitState, BuildDate string
)

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
	fmt.Println()
	logger.Debug("Securely disposing of sensitive memory storage...")
	memguard.DestroyAll()
	server.Shutdown()
	memguard.SafeExit(0)
}

func main() {
	memguard.DisableUnixCoreDumps()
	state.Init()
	if _, err := os.Stat(*args.Conf); os.IsNotExist(err) {
		createconf()
	}
	iniflags.Parse()
	memguard.CatchInterrupt(cleanup)
	switch {
	case *args.Version:
		printversioninfo()
		memguard.SafeExit(0)
	case *args.Help:
		version.Print()
		fmt.Printf("%s [-options] [command] [command args...]\n", os.Args[0])
		fmt.Printf("Usage:\n\n")
		flag.Usage()
		memguard.SafeExit(0)
	case *args.TestNet && *args.RegTest:
		printerr("Error: testnet and regtest cannot be set at the same time\n")
		memguard.SafeExit(1)
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