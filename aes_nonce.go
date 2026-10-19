package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha512"
    "encoding/hex"
    "fmt"
    "os"
    "strings"
)

var ivKey = make([]byte, aes.BlockSize)

func hexswap(input string) string {
    var result strings.Builder
    for i := len(input); i > 0; {
        start := i - 2
        if start < 0 {
            start = 0
        }
        result.WriteString(input[start:i])
        i = start
    }
    return result.String()
}

func parseNonce(nonce string) ([]byte, error) {
    if len(nonce) > 16 {
        nonce = nonce[:16]
    }

    parsed := hexswap(nonce)
    parsed = strings.Repeat("0", 32-len(parsed)) + parsed

    return hex.DecodeString(parsed)
}

func entangleNonce(key, nonce []byte) ([]byte, error) {
    if len(key) != 16 {
        return nil, fmt.Errorf("AES-128 key must be 16 bytes")
    }
    if len(nonce)%aes.BlockSize != 0 {
        return nil, fmt.Errorf("nonce length must be a multiple of 16 bytes")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    encrypted := make([]byte, len(nonce))
    cipher.NewCBCEncrypter(block, ivKey).CryptBlocks(encrypted, nonce)

    fmt.Printf("Encrypted Generator: %x\n", encrypted)

    digest := sha512.Sum384(encrypted)
    return digest[:32], nil // First 64 hexadecimal characters.
}

func main() {
    if len(os.Args) != 3 {
        fmt.Printf("Error: Expected 2 arguments, got %d\n", len(os.Args)-1)
        os.Exit(1)
    }

    keyHex := os.Args[1]
    nonceArg := os.Args[2]

    if len(keyHex) != 32 {
        fmt.Println("Error: Invalid input, check that '0x' is stripped. Or check key/nonce!")
        os.Exit(1)
    }

    key, err := hex.DecodeString(keyHex)
    if err != nil {
        fmt.Println("Error: Invalid hexadecimal key")
        os.Exit(1)
    }

    nonce, err := parseNonce(nonceArg)
    if err != nil {
        fmt.Println("Error: Invalid hexadecimal nonce")
        os.Exit(1)
    }

    // parse_key() in the Python script swaps each 4-byte segment twice,
    // which leaves the key unchanged.
    result, err := entangleNonce(key, nonce)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Entangled Nonce: %x\n", result)
}
