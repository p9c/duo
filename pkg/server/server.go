package server

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	"gitlab.com/parallelcoin/duo/pkg/wallet"
)

const (
	// LockFilename is the name of the lock file that indicates another duo is running on that directory
	LockFilename = "/.lock"
)

var (
	// Walletdb is a package centralised store for the server wallet
	Walletdb *wallet.DB
)

func lockDataDir() bool {
	if _, err := os.Stat(*args.DataDir + LockFilename); os.IsNotExist(err) {
		f, err := os.Create(*args.DataDir + LockFilename)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err)
		}
		f.Close()
	} else {
		logger.Debug("Lock already exists, server is probably already running")
		logger.Debug("If you are sure a server is not running, delete", *args.DataDir+LockFilename)
		fmt.Println()
		os.Exit(2)
	}
	logger.Debug("Locked data directory")
	return true
}

func unlockDataDir() bool {
	if _, err := os.Stat(*args.DataDir + LockFilename); os.IsNotExist(err) {
		return false
	}
	err := os.Remove(*args.DataDir + LockFilename)
	if err != nil {
		return false
	}
	logger.Debug("Unlocked data directory")
	return true
}

// Start up a server (network client)
func Start() {
	logger.Debug("Starting up server ...")
	wallet.Db.SetFilename(*args.DataDir + "/" + *args.Wallet)
	if lockDataDir() {
		err := wallet.Db.Verify()
		if err != nil {
			logger.Debug("Wallet did not verify", err)
			// Add automatic salvage here ... (maybe only possible with newer bdb?)
			os.Exit(1)
		}
		err = wallet.Db.Open()
		if err != nil {
			logger.Debug("Could not open wallet", err)
			os.Exit(1)
		}
		select {}
	} else {
		Shutdown()
	}
}

// Shutdown a server
func Shutdown() {
	if err := wallet.Db.Close(); err != nil {
		logger.Debug("error closing DB", err)
	}
	if !unlockDataDir() {
		logger.Debug("lock must have been removed by another process")
	}
	logger.Debug("Completed shutdown")
	os.Exit(0)
}

// GetWarnings -
func GetWarnings(s string) string {
	return ""
}
