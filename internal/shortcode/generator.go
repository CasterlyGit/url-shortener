package shortcode

import (
    "crypto/rand"
    "fmt"
    "strings"
)

const (
    // Base62 characters for short codes
    base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    shortCodeLength = 6
)

// GenerateRandom generates a random short code using crypto/rand
func GenerateRandom() (string, error) {
    bytes := make([]byte, shortCodeLength)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate random bytes: %w", err)
    }
    
    // Convert to base62
    for i := range bytes {
        bytes[i] = base62Chars[bytes[i]%byte(len(base62Chars))]
    }
    
    return string(bytes), nil
}

// GenerateFromID generates a short code from a numeric ID (for future use with Snowflake)
func GenerateFromID(id int64) string {
    if id == 0 {
        return string(base62Chars[0])
    }
    
    var result strings.Builder
    base := int64(len(base62Chars))
    
    for id > 0 {
        result.WriteByte(base62Chars[id%base])
        id = id / base
    }
    
    // Reverse the string
    runes := []rune(result.String())
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    return string(runes)
}