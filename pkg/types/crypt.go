package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"time"
)

/*
Crypt is a secure data store for sensitive data. It consists of two fenced memory buffers, one for the password, one for the actual symmetric ciphertext, and a regular buffer that stores the data encrypted. It is intended to be used to minimise the number of times sensitive data ends up being copied around in memory.

The accessors for internal values all make copies into the same type of encapsulation rather than returning a reference to prevent a programmer accidentally messing up a Crypt variable's internals and other parts of an application using the values, so good memory hygiene dictates one should make sure to destroy especially LockedBuffers out of these functions.

It is important to note that this encryption library is only for encrypting individual strings of data for use in a relatively small database, and not for securing communications streams, as it would need to be extended with a byte counter and IV regenerator for the main IV used with encryption/decryption that triggers after some number of bytes (less than 4gb).

A database that uses this encryption must store the initialisation vector, crypt and iteration count, and from the password a GCM blockmode cipher is created that can be used to encrypt and decrypt.
*/
type Crypt struct {
	password      *Password
	crypt         *Bytes
	iv            *Bytes
	ciphertext    *LockedBuffer
	iterations    int
	locked, armed bool
	gcm           *cipher.AEAD
	err           error
}
type crypt interface {
	Generate(*Password) *Crypt
	Copy() *Crypt
	Crypt() *Bytes
	Load(*Bytes) *Crypt
	Password() *Password
	Unlock(*Password) *Crypt
	Lock() *Crypt
	IsLocked() bool
	IV() *Bytes
	SetRandomIV() *Crypt
	SetIV(*Bytes) *Crypt
	Iterations() int
	SetIterations(int) *Crypt
	Ciphertext() *LockedBuffer
	Arm() *Crypt
	Disarm() *Crypt
	IsArmed() bool
	Encrypt(*LockedBuffer) *Bytes
	Decrypt(*Bytes) *LockedBuffer
	Delete()
}

// NewCrypt makes a new Crypt variable
func NewCrypt() (c *Crypt) {
	c = new(Crypt)
	return
}

/*
Generate creates a new crypt (password encrypted symmetric key), random initialisation vector, based on a password, encrypts the crypt.

Because all required data is available from this function it arms the crypt after encrypting the generated ciphertext enabling it to be immediately put to use. Note that one can immediately chain a decrypt or encrypt function after this function invocation if desired.

The consumer of this library should be saving the crypt, IV and iteration count alongside the data that has been encrypted with the ciphertext to enable the data to be recovered.

For this purpose, and the reason for making this library, the collection of structures containing data requiring encryption when not in use can have direct access to the GCM to provide the decrypted data to other functions using the data, by embedding this class alongside the crypt and the decrypt function returns a LockedBuffer containing the sensitive data, which can be wrapped in the function that returns the value in the crypt to callers.
*/
func (c *Crypt) Generate(p *Password) *Crypt {
	c.SetRandomIV()
	c.SetIterations(KDFBench(time.Second))
	c.ciphertext.FromRandomBytes(32)
	var LB *LockedBuffer
	var cipherIV *Bytes
	LB, cipherIV, c.err = kdf(p, c.IV(), c.Iterations())
	if c.err != nil {
		return c
	}
	defer LB.Delete()
	var block cipher.Block
	block, c.err = aes.NewCipher(*LB.Buffer())
	if c.err != nil {
		return c
	}
	var gcm cipher.AEAD
	gcm, c.err = cipher.NewGCM(block)
	crypt := gcm.Seal(nil, *cipherIV.Buffer(), *LB.Buffer(), nil)
	cipherIV.Zero()
	cipherIV = nil
	c.Load(NewBytes().FromByteSlice(&crypt))
	c.Unlock(p)
	c.Arm()
	return c
}

// Load moves the encrypted data into the crypt
func (c *Crypt) Load(b *Bytes) *Crypt {
	c.crypt.FromBytes(b)
	return c
}

// Copy makes a copy of a crypt
func (c *Crypt) Copy() (C *Crypt) {
	C = new(Crypt)
	C.crypt = c.crypt.Copy()
	C.password = c.password.Copy()
	C.ciphertext = c.ciphertext.Copy()
	return
}

// Password returns a copy of the password stored in the Crypt
func (c *Crypt) Password() (P *Password) {
	P = c.password.Copy()
	return
}

// Crypt returns a copy of the encrypted data in the Crypt
func (c *Crypt) Crypt() (B *Bytes) {
	B = c.crypt.Copy()
	return
}

// Ciphertext returns a copy of the LockedBuffer contining unencrypted data stored in the Crypt.
func (c *Crypt) Ciphertext() (LB *LockedBuffer) {
	LB = c.ciphertext.Copy()
	return
}

// IV returns the initialisation vector used by the ciphertext, and concatenated with the password in each cycle of the hashchain iteration.
func (c *Crypt) IV() *Bytes {
	return c.iv.Copy()
}

// SetRandomIV returns a new 16 byte initialisation vector from a cryptographically secure random source
func (c *Crypt) SetRandomIV() *Crypt {
	i := *c.iv.Buffer()
	rand.Read(i)
	return c
}

// SetIV sets the initialisation vector from a Bytes
func (c *Crypt) SetIV(B *Bytes) *Crypt {
	c.iv.FromBytes(B)
	return c
}

// Iterations returns the number of iterations configured in the Crypt
func (c *Crypt) Iterations() int {
	return c.iterations
}

// SetIterations sets the number of iterations to use in the KDF
func (c *Crypt) SetIterations(I int) *Crypt {
	c.iterations = I
	return c
}

// Unlock loads the password into its LockedBuffer and arms it by default if the crypt has been loaded with encrypted data
func (c *Crypt) Unlock(p *Password) *Crypt {
	if c.password != nil {
		c.password.Delete()
	}
	c.password = p
	if !c.armed {
		c.Arm()
	}
	c.locked = false
	return c
}

// Lock disarms the crypt if it is armed, and removes the password
func (c *Crypt) Lock() *Crypt {
	c.password.Delete()
	c.Disarm()
	c.locked = true
	return c
}

// IsLocked returns true if the crypt is locked (does not have the password)
func (c *Crypt) IsLocked() bool {
	return c.locked
}

/*
Arm runs the KDF on the password and uses the resultant key to unlock the ciphertext from the crypt enabling the encryption and decryption of data.

This function assumes the IV, encrypted form of the key is present, that it has not been marked as armed, the password has been loaded in, and that the iteration count is valid.

If the content of the IV, crypt and password are not correct the result is not defined.
*/
func (c *Crypt) Arm() *Crypt {
	var ciphertext *LockedBuffer
	var cipherIV *Bytes
	if c.iv != nil &&
		c.crypt != nil &&
		c.password != nil &&
		c.iterations > 0 &&
		!c.IsLocked() {
		ciphertext, cipherIV, c.err = kdf(c.password, c.iv, c.iterations)
	} else {
		c.err = errors.New("Crypt is not fully populated")
		return c
	}
	var block cipher.Block
	block, c.err = aes.NewCipher(*ciphertext.Buffer())
	if c.err != nil {
		return c
	}
	var gcm cipher.AEAD
	gcm, c.err = cipher.NewGCM(block)
	if c.err != nil {
		return c
	}
	ctb := *c.ciphertext.WithSize(32).Buffer()
	_, c.err = gcm.Open(ctb, *cipherIV.Buffer(), *c.crypt.Buffer(), nil)
	if c.err != nil {
		return c
	}
	cipherIV.Zero()
	cipherIV = nil
	block, c.err = aes.NewCipher(*c.ciphertext.Buffer())
	if c.err != nil {
		return c
	}
	gcm, c.err = cipher.NewGCM(block)
	if c.err != nil {
		return c
	}
	c.gcm = &gcm
	c.armed = true
	return c
}

// Disarm clears the unlocked ciphertext and disables encryption/decryption
func (c *Crypt) Disarm() *Crypt {
	c.ciphertext.Delete()
	c.armed = false
	return c
}

// IsArmed returns true if the ciphertext has been decrypted
func (c *Crypt) IsArmed() bool {
	return c.armed
}

// Encrypt encrypts the contents of a LockedBuffer and returns the encrypted form as Bytes
func (c *Crypt) Encrypt(lb *LockedBuffer) (B *Bytes) {
	b := (*c.gcm).Seal(nil, *c.IV().Buffer(), *lb.Buffer(), nil)
	B.FromByteSlice(&b)
	return
}

// Decrypt decrypts the contents of a Bytes buffer and returns a LockedBuffer containing the decrypted plaintext
func (c *Crypt) Decrypt(ciphertext *Bytes) (LB *LockedBuffer) {
	var b []byte
	b, c.err = (*c.gcm).Open(nil, *c.IV().Buffer(), *ciphertext.Buffer(), nil)
	LB.FromByteSlice(&b)
	return
}

// Delete destroys all the lockedbuffers in a Crypt
func (c *Crypt) Delete() {
	c.ciphertext.Delete()
	c.password.Delete()
	iv := *c.iv.Buffer()
	for i := range iv {
		iv[i] = 0
	}
	crypt := *c.crypt.Buffer()
	for i := range crypt {
		crypt[i] = 0
	}
}
