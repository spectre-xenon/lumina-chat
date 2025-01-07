package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func GenerateHashString(password string) (string, error) {
	// Recommended params
	defaultParams := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		keyLength:   32,
	}

	salt := make([]byte, 16)

	// Generate a cryptographically secure random salt.
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// base64 encoded from of the hash
	hashString := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash)

	return hashString, nil
}

func CompareHashStrings(password, encodedHash string) bool {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash := decodeHash(encodedHash)

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}
	return false
}

func decodeHash(encodedHash string) (p *params, salt []byte, hash []byte) {
	vals := strings.Split(encodedHash, "$")

	var version int
	fmt.Sscanf(vals[2], "v=%d", &version)

	p = &params{}
	fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)

	salt, _ = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	p.saltLength = uint32(len(salt))

	hash, _ = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	p.keyLength = uint32(len(hash))

	return p, salt, hash
}
