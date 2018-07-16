package bdb

/*
	#cgo LDFLAGS: -ldb
	#include <stdlib.h>
	#include <errno.h>
	#include <db.h>

	int db_env_set_encrypt(DB_ENV *env, const char *passwd, u_int32_t flags) {
		return env->set_encrypt(env, passwd, flags);
	}
	int db_env_open(DB_ENV *env, const char *home, u_int32_t flags, int mode) {
		return env->open(env, home, flags, mode);
	}
	int db_env_close(DB_ENV *env, u_int32_t flags) {
		return env->close(env, flags);
	}
	int db_set_encrypt(DB *db, const char *passwd, u_int32_t flags) {
		return db->set_encrypt(db, passwd, flags);
	}
	int db_open(DB *db, DB_TXN *txn, const char *file, const char *database, DBTYPE type, u_int32_t flags, int mode) {
		return db->open(db, txn, file, database, type, flags, mode);
	}
	int db_close(DB *db, u_int32_t flags) {
		return db->close(db, flags);
	}
	int db_verify(DB *db, const char *file, const char *database, FILE *outfile, u_int32_t flags) {
		return db->verify(db, file, database, outfile, flags);
	}
	int db_get_type(DB *db, DBTYPE *type) {
		return db->get_type(db, type);
	}
	int db_put(DB *db, DB_TXN *txn, DBT *key, DBT *data, u_int32_t flags) {

		return db->put(db, txn, key, data, flags);
	}
	int db_get(DB *db, DB_TXN *txn, DBT *key, DBT *data, u_int32_t flags) {
		return db->get(db, txn, key, data, flags);
	}
	int db_del(DB *db, DB_TXN *txn, DBT *key, u_int32_t flags) {
		return db->del(db, txn, key, flags);
	}
	int db_cursor(DB *db, DB_TXN *txn, DBC **cursor, u_int32_t flags) {
		return db->cursor(db, txn, cursor, flags);
	}
	int db_cursor_close(DBC *cur) {
		return cur->close(cur);
	}
	int db_cursor_get(DBC *cur, DBT *key, DBT *data, u_int32_t flags) {
		return cur->get(cur, key, data, flags);
	}
	int db_cursor_del(DBC *cur, u_int32_t flags) {
		return cur->del(cur, flags);
	}

	int db_env_txn_begin(DB_ENV *env, DB_TXN *parent, DB_TXN **txn, u_int32_t flags) {
		return env->txn_begin(env, parent, txn, flags);
	}
	int db_txn_abort(DB_TXN *txn) {
		return txn->abort(txn);
	}
	int db_txn_commit(DB_TXN *txn, u_int32_t flags) {
		return txn->commit(txn, flags);
	}
	int db_is_alive(DB_ENV *dbenv, int is_alive) {
		return dbenv->set_isalive(dbenv, 0);
	}
*/
import "C"

import (
	"os"
	"unsafe"
)

// DatabaseType is the type of database being used
type DatabaseType int

// Available database types.
const (
	BTree    = DatabaseType(C.DB_BTREE)
	Hash     = DatabaseType(C.DB_HASH)
	Numbered = DatabaseType(C.DB_RECNO)
	Queue    = DatabaseType(C.DB_QUEUE)
	Unknown  = DatabaseType(C.DB_UNKNOWN)
)

// DatabaseConfig - Database configuration.
type DatabaseConfig struct {
	Create          bool         // Create the database, if necessary.
	Mode            os.FileMode  // File creation mode for the database.
	Password        string       // Encryption password or an empty string.
	Name            string       // Identifier of the database inside the file.
	Type            DatabaseType // Type of database to create
	ReadUncommitted bool         // Enable support for read-uncommitted isolation.
	Snapshot        bool         // Enable support for snapshot isolation.
}

// Database -
type Database struct {
	ptr *C.DB
}

// Errno is the error return value for barkeleydb functions
type Errno int

// EnvFlag is the settings for a barkeleydb database environment
type EnvFlag int

// Status codes representing common errors.
const (
	ErrAgain           = Errno(C.EAGAIN)
	ErrInvalid         = Errno(C.EINVAL)
	ErrNoEntry         = Errno(C.ENOENT)
	ErrExists          = Errno(C.EEXIST)
	ErrAccess          = Errno(C.EACCES)
	ErrNoSpace         = Errno(C.ENOSPC)
	ErrPermission      = Errno(C.EPERM)
	ErrRunRecovery     = Errno(C.DB_RUNRECOVERY)
	ErrVersionMismatch = Errno(C.DB_VERSION_MISMATCH)
	ErrOldVersion      = Errno(C.DB_OLD_VERSION)
	ErrLockDeadlock    = Errno(C.DB_LOCK_DEADLOCK)
	ErrLockNotGranted  = Errno(C.DB_LOCK_NOTGRANTED)
	ErrBufferTooSmall  = Errno(C.DB_BUFFER_SMALL)
	ErrSecondaryBad    = Errno(C.DB_SECONDARY_BAD)
	ErrForeignConflict = Errno(C.DB_FOREIGN_CONFLICT)
	ErrKeyExists       = Errno(C.DB_KEYEXIST)
	ErrKeyEmpty        = Errno(C.DB_KEYEMPTY)
	ErrNotFound        = Errno(C.DB_NOTFOUND)
	// Environment
	DbAggressive = EnvFlag(C.DB_AGGRESSIVE)
	DbSalvage    = EnvFlag(C.DB_SALVAGE)
)

func (db Database) is_alive(env *C.DB_ENV, pid C.pid_t, threadid C.db_threadid_t, i C.u_int32_t) (r int) {
	return 1
}

// Error turns a status code into a human readable message.
func (err Errno) Error() string {
	return C.GoString(C.db_strerror(C.int(err)))
}

// Check a function result and return an error if necessary.
func check(rc C.int) (err error) {
	if rc != 0 {
		err = Errno(rc)
	}
	return
}

type EnvironmentConfig struct {
	Create        bool        // Create the environment, if necessary.
	Mode          os.FileMode // File creation mode for the environment.
	Password      string      // Encryption password or an empty string.
	Recover       bool        // Run recovery on the environment, if necessary.
	Transactional bool        // Enable transactions in the environment.
	NoSync        bool        // Do not flush to log when committing.
	WriteNoSync   bool        // Do not flush log when committing.
}

// Database environment.
type Environment struct {
	ptr *C.DB_ENV
}

// Special constant to indicate no environment should be used.
var NoEnvironment = Environment{ptr: nil}

type IsolationLevel int

// Available transaction isolation levels.
const (
	ReadCommitted   = IsolationLevel(C.DB_READ_COMMITTED)
	ReadUncommitted = IsolationLevel(C.DB_READ_UNCOMMITTED)
	Snapshot        = IsolationLevel(C.DB_TXN_SNAPSHOT)
)

// Transaction configuration.
type TransactionConfig struct {
	Parent      Transaction    // Parent transaction.
	Isolation   IsolationLevel // Transaction isolation level.
	Bulk        bool           // Optimize for bulk insertions.
	NoWait      bool           // Fail instead of waiting for locks.
	NoSync      bool           // Do not flush to log when committing.
	WriteNoSync bool           // Do not flush log when committing.
}

// Transaction in a database environment.
type Transaction struct {
	ptr *C.DB_TXN
}

// Special constant indicating no transaction should be used.
var NoTransaction = Transaction{ptr: nil}

func OpenEnvironment(home string, config *EnvironmentConfig) (env Environment, err error) {
	err = check(C.db_env_create(&env.ptr, 0))
	if err == nil {
		defer func() {
			if err != nil && env.ptr != nil {
				C.db_env_close(env.ptr, 0)
				env.ptr = nil
			}
		}()
	} else {
		return
	}

	var mode C.int = 0
	var flags C.u_int32_t = C.DB_THREAD
	var chome, cpassword *C.char

	chome = C.CString(home)
	defer C.free(unsafe.Pointer(chome))

	if config != nil {
		if config.Create {
			flags |= C.DB_CREATE
		}
		if config.Mode != 0 {
			mode = C.int(config.Mode)
		}
		if len(config.Password) > 0 {
			cpassword := C.CString(config.Password)
			C.free(unsafe.Pointer(cpassword))
		}
		if config.Recover {
			flags |= C.DB_REGISTER | C.DB_RECOVER // | C.DB_FAILCHK
		}
		if config.Transactional {
			flags |= C.DB_INIT_TXN | C.DB_INIT_MPOOL
		}
		if config.NoSync {
			flags |= C.DB_TXN_NOSYNC
		}
		if config.WriteNoSync {
			flags |= C.DB_TXN_WRITE_NOSYNC
		}
	}

	if cpassword != nil {
		err = check(C.db_env_set_encrypt(env.ptr, cpassword, 0))
		if err != nil {
			return
		}
	}
	err = check(C.db_env_open(env.ptr, chome, flags, mode))

	return
}

// Close the environment
func (env Environment) Close() (err error) {
	err = check(C.db_env_close(env.ptr, C.u_int32_t(C.DB_FORCESYNC)))
	return
}

// OpenDatabase -
func OpenDatabase(env Environment, txn Transaction, file string, config *DatabaseConfig) (db Database, err error) {
	err = check(C.db_create(&db.ptr, env.ptr, 0))
	if err == nil {
		defer func() {
			if err != nil && db.ptr != nil {
				C.db_close(db.ptr, 0)
				db.ptr = nil
			}
		}()
	} else {
		return
	}

	var mode C.int = 0
	var flags C.u_int32_t = C.DB_THREAD
	var cfile, cpassword, cname *C.char
	var dbtype C.DBTYPE = C.DB_UNKNOWN

	if len(file) > 0 {
		cfile = C.CString(file)
		defer C.free(unsafe.Pointer(cfile))
	}

	if config != nil {
		if config.Create {
			flags |= C.DB_CREATE
		}
		if config.Mode != 0 {
			mode = C.int(config.Mode)
		}
		if len(config.Password) > 0 {
			cpassword := C.CString(config.Password)
			defer C.free(unsafe.Pointer(cpassword))
		}
		if len(config.Name) > 0 {
			cname = C.CString(config.Name)
			defer C.free(unsafe.Pointer(cname))
		}
		if config.Type != 0 {
			dbtype = C.DBTYPE(config.Type)
		}
		if config.ReadUncommitted {
			flags |= C.DB_READ_UNCOMMITTED
		}
		if config.Snapshot {
			flags |= C.DB_MULTIVERSION
		}
	}

	if cpassword != nil {
		err = check(C.db_set_encrypt(db.ptr, cpassword, 0))
		if err != nil {
			return
		}
	}

	err = check(C.db_open(db.ptr, txn.ptr, cfile, cname, dbtype, flags, mode))

	return
}

// Close the database.
func (db Database) Close() (err error) {
	err = check(C.db_close(db.ptr, 0))
	return
}

// Verify the database
func Verify(file string) (err error) {
	var db Database
	err = check(C.db_create(&db.ptr, nil, 0))
	if err == nil {
		defer func() {
			if err != nil && db.ptr != nil {
				C.db_close(db.ptr, 0)
				db.ptr = nil
			}
		}()
	} else {
		return
	}

	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	err = check(C.db_verify(db.ptr, cfile, nil, nil, 0))
	return
}

// Get the type of the database.
func (db Database) Type() (dbtype DatabaseType, err error) {
	var cdbtype C.DBTYPE
	err = check(C.db_get_type(db.ptr, &cdbtype))
	dbtype = DatabaseType(cdbtype)
	return
}

func (db Database) Put(txn Transaction, append bool, recs ...[2][]byte) (err error) {
	dbtype, err := db.Type()
	if err != nil {
		return
	}
	var key, data C.DBT
	var flags C.u_int32_t = 0
	if append {
		key.flags |= C.DB_DBT_USERMEM

		switch dbtype {
		case Numbered, Queue:
			flags |= C.DB_APPEND
		default:
			flags |= C.DB_NOOVERWRITE
		}
	} else {
		key.flags |= C.DB_DBT_READONLY
	}
	data.flags |= C.DB_DBT_READONLY
	for rec := range recs {
		key.size = C.u_int32_t(len(recs[rec][0]))
		key.data = unsafe.Pointer(C.CString(string(recs[rec][0])))
		data.size = C.u_int32_t(len(recs[rec][1]))
		data.data = unsafe.Pointer(C.CString(string(recs[rec][1])))
		r := C.db_put(db.ptr, txn.ptr, &key, &data, flags)
		err = check(r)

		if err != nil {
			return
		}
	}
	return
}

func (db Database) Get(txn Transaction, consume bool, recs ...[2][]byte) (err error) {
	var key, data C.DBT
	var flags C.u_int32_t = 0

	if consume {
		key.flags |= C.DB_DBT_USERMEM
		flags |= C.DB_CONSUME_WAIT
	} else {
		key.flags |= C.DB_DBT_READONLY
	}
	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)
	for rec := range recs {
		key.size = C.u_int32_t(len(recs[0][rec]))
		key.data = unsafe.Pointer(C.CString(string(recs[0][rec])))
		data.size = C.u_int32_t(len(recs[1][rec]))
		data.data = unsafe.Pointer(C.CString(string(recs[1][rec])))

		err = check(C.db_get(db.ptr, txn.ptr, &key, &data, flags))
		if err != nil {
			return
		}
	}
	return
}

func (db Database) Del(txn Transaction, recs ...[]byte) (err error) {
	var key C.DBT
	key.flags |= C.DB_DBT_READONLY
	for rec := range recs {
		key.size = C.u_int32_t(len(recs[rec]))
		key.data = unsafe.Pointer(C.CString(string(recs[rec])))
		err = check(C.db_del(db.ptr, txn.ptr, &key, 0))
		if err != nil {
			return
		}
	}
	return
}

type Cursor struct {
	db  Database
	ptr *C.DBC
}

// Obtain a cursor over the database.
func (db Database) Cursor(txn Transaction) (cur Cursor, err error) {
	cur.db = db
	err = check(C.db_cursor(db.ptr, txn.ptr, &cur.ptr, 0))
	return
}

// Close the cursor.
func (cur Cursor) Close() (err error) {
	err = check(C.db_cursor_close(cur.ptr))
	return
}
func (cur Cursor) Set(rec [2][]byte, exact bool) (err error) {
	var key, data C.DBT
	var flags C.u_int32_t = 0

	if exact {
		key.flags |= C.DB_DBT_READONLY
		flags |= C.DB_SET
	} else {
		key.flags |= C.DB_DBT_MALLOC
		flags |= C.DB_SET_RANGE
	}

	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)

	odata := key.data
	defer func() {
		if key.data != odata {
			C.free(data.data)
		}
	}()

	err = check(C.db_cursor_get(cur.ptr, &key, &data, flags))
	if err != nil {
		return
	}
	rec[0] = *(*[]byte)(unsafe.Pointer(&key.data))
	rec[1] = *(*[]byte)(unsafe.Pointer(&data.data))
	return
}

// Retrieve the first record of the database.
func (cur Cursor) First(rec *[2][]byte) (err error) {
	var key, data C.DBT

	key.flags |= C.DB_DBT_REALLOC
	defer C.free(key.data)
	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)

	err = check(C.db_cursor_get(cur.ptr, &key, &data, C.DB_FIRST))
	if err != nil {
		return
	}
	rec[0] = *(*[]byte)(unsafe.Pointer(&key.data))
	rec[1] = *(*[]byte)(unsafe.Pointer(&data.data))
	return
}

// Retrieve the next record from the cursor.
func (cur Cursor) Next(rec *[2][]byte) (err error) {
	var key, data C.DBT

	key.flags |= C.DB_DBT_REALLOC
	defer C.free(key.data)
	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)

	err = check(C.db_cursor_get(cur.ptr, &key, &data, C.DB_NEXT))
	if err != nil {
		return
	}
	rec[0] = *(*[]byte)(unsafe.Pointer(&key.data))
	rec[1] = *(*[]byte)(unsafe.Pointer(&data.data))
	return
}

// Retrieve the last record of the database.
func (cur Cursor) Last(rec [2][]byte) (err error) {
	var key, data C.DBT

	key.flags |= C.DB_DBT_REALLOC
	defer C.free(key.data)
	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)

	err = check(C.db_cursor_get(cur.ptr, &key, &data, C.DB_LAST))
	if err != nil {
		return
	}
	rec[0] = *(*[]byte)(unsafe.Pointer(&key.data))
	rec[1] = *(*[]byte)(unsafe.Pointer(&data.data))
	return
}

// Retrieve the previous record from the cursor.
func (cur Cursor) Prev(rec [2][]byte) (err error) {
	var key, data C.DBT

	key.flags |= C.DB_DBT_REALLOC
	defer C.free(key.data)
	data.flags |= C.DB_DBT_REALLOC
	defer C.free(data.data)

	err = check(C.db_cursor_get(cur.ptr, &key, &data, C.DB_PREV))
	if err != nil {
		return
	}
	rec[0] = *(*[]byte)(unsafe.Pointer(&key.data))
	rec[1] = *(*[]byte)(unsafe.Pointer(&data.data))
	return
}

// Delete the current record at the cursor.
func (cur Cursor) Del() (err error) {
	err = check(C.db_cursor_del(cur.ptr, 0))
	return
}
