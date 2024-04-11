package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func MergeFiberMaps(init, override *fiber.Map) fiber.Map {
	allData := make(fiber.Map)

	if init != nil {
		for key, value := range *init {
			allData[key] = value
		}
	}

	if override != nil {
		for key, value := range *override {
			allData[key] = value
		}
	}

	return allData
}

func GetMD5Hash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
