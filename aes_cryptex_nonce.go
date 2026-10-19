package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
    "os"
)

var ivKey = make([]byte, aes.BlockSize)

func parseSeed(seed string) ([]byte, error) {
    return hex.DecodeString(seed)
}

func entangleNonce(key, seed []byte) ([]byte, error) {
    if len(key) != 16 {
        return nil, fmt.Errorf("AES-128 key must be 16 bytes")
    }
    if len(seed)%aes.BlockSize != 0 {
        return nil, fmt.Errorf("seed length must be a multiple of 16 bytes")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    encrypted := make([]byte, len(seed))
    cipher.NewCBCEncrypter(block, ivKey).CryptBlocks(encrypted, seed)

    fmt.Printf("Encrypted Seed: %x\n", encrypted)

    digest := sha512.Sum384(encrypted)
    return digest[:32], nil // First 64 hexadecimal characters.
}

func main() {
    if len(os.Args) != 3 {
        fmt.Printf("Error: Expected 2 arguments, got %d\n", len(os.Args)-1)
        os.Exit(1)
    }

    keyHex := os.Args[1]
    seedArg := os.Args[2]

    if len(keyHex) != 32 {
        fmt.Println("Error: Invalid input, check that '0x' is stripped. Or check key/seed!")
        os.Exit(1)
    }

    key, err := hex.DecodeString(keyHex)
    if err != nil {
        fmt.Println("Error: Invalid hexadecimal key")
        os.Exit(1)
    }

    seed, err := parseSeed(seedArg)
    if err != nil {
        fmt.Println("Error: Invalid hexadecimal seed")
        os.Exit(1)
    }

    // The original parse_key() swaps each 4-byte segment twice,
    // leaving the key unchanged.
    result, err := entangleNonce(key, seed)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Entangled Nonce: %x\n", result)
}
