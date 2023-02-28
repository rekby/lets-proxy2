package fastuuid

import (
	"errors"
	"github.com/rekby/fastuuid/internal"
)

// UUIDv4 generate uuid as byte array without allocations in heap
//
// It is safe calling this function from concurrent goroutines.
func UUIDv4() (binaryUUID [16]byte, err error) {
	randomReader := internal.GetRandomReader()
	err = randomReader.ReadFull(binaryUUID[:])
	randomReader.Release()

	if err != nil {
		return binaryUUID, err
	}

	binaryUUID[6] = (binaryUUID[6] & 0x0f) | 0x40 // Version 4
	binaryUUID[8] = (binaryUUID[8] & 0x3f) | 0x80 // Variant is 10

	return binaryUUID, nil
}

// UUIDv4String generate uuid random string with only one allocation
//
// It is safe calling this function from concurrent goroutines.
func UUIDv4String() (string, error) {
	var buf = make([]byte, 36)

	// use buf for binary and text uuid same time
	// reuse same buffer need for prevent second allocation by escape analise for buffer
	// it can leak through bufio.Reader
	binaryUUID := buf[20:]

	randomReader := internal.GetRandomReader()
	err := randomReader.ReadFull(binaryUUID)
	randomReader.Release()
	if err != nil {
		return "", err
	}

	binaryUUID[6] = (binaryUUID[6] & 0x0f) | 0x40 // Version 4
	binaryUUID[8] = (binaryUUID[8] & 0x3f) | 0x80 // Variant is 10

	internal.EncodeHex(buf[:], binaryUUID)

	return internal.UnsafeConvertToString(buf), nil
}

// UUIDv4StringBytes generate uuid and render it as string to dst buffer without allocations
//
// It is safe calling this function from concurrent goroutines.
func UUIDv4StringBytes(buf []byte) error {
	if len(buf) < 36 {
		return errors.New("fastuuid: buf size must be 36 bytes or more")
	}

	// use buf for binary and text uuid same time
	// reuse same buffer need for prevent second allocation by escape analise for buffer
	// it can leak through bufio.Reader
	binaryUUID := buf[20:]

	randomReader := internal.GetRandomReader()
	err := randomReader.ReadFull(binaryUUID)
	randomReader.Release()
	if err != nil {
		return err
	}

	binaryUUID[6] = (binaryUUID[6] & 0x0f) | 0x40 // Version 4
	binaryUUID[8] = (binaryUUID[8] & 0x3f) | 0x80 // Variant is 10

	internal.EncodeHex(buf[:], binaryUUID)

	return nil
}
