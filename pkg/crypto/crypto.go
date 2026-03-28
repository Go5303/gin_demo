package crypto

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_."
const ikey = "-x6g6ZWm2G9g_vr0Bo.pOq3kRIxsZ6rm"

// MD5 returns md5 hash of a string
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// MD6 matches PHP md6() function
func MD6(str, salt, preStr string) string {
	if preStr == "" {
		preStr = "agg_"
	}
	return MD5(preStr + MD5(str) + salt)
}

// Encrypt matches PHP encrypt() function exactly
// This is needed for compatibility with existing cookies/tokens
func Encrypt(txt, key string) string {
	if txt == "" {
		return txt
	}
	if key == "" {
		key = MD5(key)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nh1 := r.Intn(65)
	nh2 := r.Intn(65)
	nh3 := r.Intn(65)
	ch1 := chars[nh1]
	ch2 := chars[nh2]
	ch3 := chars[nh3]
	nhnum := nh1 + nh2 + nh3

	knum := 0
	for i := 0; i < len(key); i++ {
		knum += int(key[i])
	}

	mdKeyFull := MD5(MD5(MD5(key+string(ch1)) + string(ch2) + ikey) + string(ch3))
	start := nhnum % 8
	length := knum%8 + 16
	mdKey := mdKeyFull[start : start+length]

	// base64 encode with timestamp prefix
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%s", time.Now().Unix(), txt)))
	encoded = strings.NewReplacer("+", "-", "/", "_", "=", ".").Replace(encoded)

	tmp := make([]byte, len(encoded))
	k := 0
	klen := len(mdKey)
	for i := 0; i < len(encoded); i++ {
		if k == klen {
			k = 0
		}
		j := (nhnum + strings.IndexByte(chars, encoded[i]) + int(mdKey[k])) % 64
		k++
		tmp[i] = chars[j]
	}

	result := string(tmp)
	tmplen := len(result)
	result = result[:nh2%(tmplen+1)] + string(ch3) + result[nh2%(tmplen+1):]
	tmplen++
	result = result[:nh1%(tmplen+1)] + string(ch2) + result[nh1%(tmplen+1):]
	tmplen++
	result = result[:knum%(tmplen+1)] + string(ch1) + result[knum%(tmplen+1):]

	return result
}

// Decrypt matches PHP decrypt() function exactly
func Decrypt(txt, key string) string {
	if txt == "" {
		return txt
	}
	if key == "" {
		key = MD5(key)
	}

	knum := 0
	for i := 0; i < len(key); i++ {
		knum += int(key[i])
	}

	tlen := len(txt)
	if tlen == 0 {
		return ""
	}

	// Extract ch1
	pos := knum % tlen
	if pos >= len(txt) {
		return ""
	}
	ch1 := txt[pos]
	nh1 := strings.IndexByte(chars, ch1)
	txt = txt[:pos] + txt[pos+1:]
	tlen--

	// Extract ch2
	if tlen == 0 {
		return ""
	}
	pos = nh1 % tlen
	if pos >= len(txt) {
		return ""
	}
	ch2 := txt[pos]
	nh2 := strings.IndexByte(chars, ch2)
	txt = txt[:pos] + txt[pos+1:]
	tlen--

	// Extract ch3
	if tlen == 0 {
		return ""
	}
	pos = nh2 % tlen
	if pos >= len(txt) {
		return ""
	}
	ch3 := txt[pos]
	nh3 := strings.IndexByte(chars, ch3)
	txt = txt[:pos] + txt[pos+1:]

	nhnum := nh1 + nh2 + nh3
	mdKeyFull := MD5(MD5(MD5(key+string(ch1)) + string(ch2) + ikey) + string(ch3))
	start := nhnum % 8
	length := knum%8 + 16
	mdKey := mdKeyFull[start : start+length]

	result := make([]byte, len(txt))
	k := 0
	klen := len(mdKey)
	for i := 0; i < len(txt); i++ {
		if k == klen {
			k = 0
		}
		j := strings.IndexByte(chars, txt[i]) - nhnum - int(mdKey[k])
		k++
		for j < 0 {
			j += 64
		}
		result[i] = chars[j]
	}

	decoded := strings.NewReplacer("-", "+", "_", "/", ".", "=").Replace(string(result))
	raw, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return ""
	}

	s := strings.TrimSpace(string(raw))
	// Remove timestamp prefix (format: 1234567890_content)
	if len(s) > 11 && s[10] == '_' {
		allDigits := true
		for i := 0; i < 10; i++ {
			if s[i] < '0' || s[i] > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			s = s[11:]
		}
	}
	return s
}
