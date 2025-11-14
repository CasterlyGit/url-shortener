package shortcode

import (
    "fmt"
    "strings"

    "github.com/CasterlyGit/url-shortener/internal/snowflake"
)

const (
    base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
    node *snowflake.Node
)

func InitSnowflake(nodeID int64) error {
    var err error
    node, err = snowflake.NewNode(nodeID)
    return err
}

func GenerateFromSnowflake() (int64, error) {
    if node == nil {
        return 0, fmt.Errorf("snowflake node not initialized")
    }
    
    return node.Generate(), nil
}

// EncodeBase62 converts a int64 to base62 string
func EncodeBase62(num int64) string {
    if num == 0 {
        return string(base62Chars[0])
    }
    
    var result strings.Builder
    base := int64(len(base62Chars))
    
    for num > 0 {
        result.WriteByte(base62Chars[num%base])
        num = num / base
    }
    
    // Reverse the string
    runes := []rune(result.String())
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    return string(runes)
}

// GenerateRandom is kept for backward compatibility
func GenerateRandom() (string, error) {
    id, err := GenerateFromSnowflake()
    if err != nil {
        return "", err
    }
    return EncodeBase62(id), nil
}