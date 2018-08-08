package wallet
import (
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"os"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
)
// A berkeley DB database
type BDB struct {
	*bdb.Database
	Filename      string
}
type bDB interface {
	Open() error
	Close() error
	Verify() error
}
// Open a wallet.dat file
func (db *BDB) Open() (err error) {
	dbenvconf := bdb.EnvironmentConfig{
		Create:        true,
		Recover:       true,
		Mode:          0600,
		Transactional: true,
	}
	dbenv, err := bdb.OpenEnvironment(*args.DataDir, &dbenvconf)
	if err != nil {
		return
	}
	dbconfig := bdb.DatabaseConfig{
		Create: false,
		Mode:   0600,
		Name:   "main",
	}
	db1, err := bdb.OpenDatabase(dbenv, bdb.NoTransaction, db.Filename, &dbconfig)
	if err == nil {
		db.Database = &db1
	} else {
		logger.Debug("Failed to open database", err)
		return
	}
	return
}
// Close an wallet.dat file
func (db *BDB) Close() (err error) {
	err = db.Database.Close()
	return
}
// Verify the consistency of a wallet.dat database
func (db *BDB) Verify() (err error) {
	if _, err = os.Stat(db.Filename); os.IsNotExist(err) {
		logger.Debug(err)
		return
	}
	if err = bdb.Verify(db.Filename); err != nil {
		logger.Debug(err)
		return
	}
	return
}
// SetFilename changes the name of the database we want to open
func (db *BDB) SetFilename(filename string) {
	db.Filename = filename
}
