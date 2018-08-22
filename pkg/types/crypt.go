package types

// Crypt is a secure data store for sensitive data. It consists of two fenced memory buffers, one for the password, one for the actual symmetric ciphertext, and a regular buffer that stores the data encrypted. It is intended to be used to minimise the number of times sensitive data ends up being copied around in memory.
type Crypt struct {
	password      *Password
	encrypted     *Bytes
	ciphertext    *LockedBuffer
	locked, armed bool
}
type crypt interface {
	Password() *Password
	Crypt() *Bytes
	Ciphertext() *LockedBuffer
	Unlock(*Password) *Crypt
	Lock() *Crypt
	IsLocked() *Crypt
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed()
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
}

// NewCrypt makes a new Crypt variable
func NewCrypt() *Crypt {
	return new(Crypt)
}

// Password returns the password stored in the Crypt
func (c *Crypt) Password() (P *Password) {
	return c.password
}

// Crypt returns the encrypted data in the Crypt
func (c *Crypt) Crypt() (B *Bytes) {
	return c.encrypted
}

// Ciphertext returns the unencrypted data stored in the Crypt
func (c *Crypt) Ciphertext() (LB *LockedBuffer) {
	return c.ciphertext
}

func (c *Crypt) Unlock(*Password) *Crypt {
	return c
}

func (c *Crypt) Lock() *Crypt {
	return c
}

func (c *Crypt) IsLocked() *Crypt {
	return c
}

func (c *Crypt) Arm() *Crypt {
	return c
}

func (c *Crypt) Disarm() *Crypt {
	return c
}

func (c *Crypt) IsArmed() bool {
	return c.armed
}

func (c *Crypt) Encrypt(*LockedBuffer) (B *Bytes) {
	return
}

func (c *Crypt) Decrypt(*Bytes) *LockedBuffer {
	return c
}
