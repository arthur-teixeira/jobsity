package service

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type Argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

type HashSalt struct {
	Hash []byte
	Salt []byte
}

func NewHasher(saltLen, keyLen uint32) *Argon2Hasher {
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#introduction
	return &Argon2Hasher{
		time:    3,
		memory:  19 * 1024,
		threads: 1,
		keyLen:  keyLen,
		saltLen: saltLen,
	}
}

func randomSecret(len uint32) ([]byte, error) {
	secret := make([]byte, len)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (hasher *Argon2Hasher) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if len(salt) == 0 {
		salt, err = randomSecret(hasher.saltLen)
	}

	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, hasher.time, hasher.memory, hasher.threads, hasher.keyLen)
	return &HashSalt{Hash: hash, Salt: salt}, nil
}

func (hasher *Argon2Hasher) Compare(hash, salt, pwd []byte) bool {
	hashSalt, err := hasher.GenerateHash(pwd, salt)
	if err != nil {
		return false
	}

	return bytes.Equal(hash, hashSalt.Hash)
}
