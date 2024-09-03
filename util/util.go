package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func MergeFiberMaps(maps ...fiber.Map) *fiber.Map {
	merged := fiber.Map{}

	for _, m := range maps {
		for key, value := range m {
			if existingValue, exists := merged[key]; exists {
				if existingMap, ok := existingValue.(fiber.Map); ok {
					if valueMap, ok := value.(fiber.Map); ok {
						merged[key] = MergeFiberMaps(existingMap, valueMap)
						continue
					}
				}
			}
			merged[key] = value
		}
	}

	return &merged
}

func OverrideFiberMaps(init, override *fiber.Map) fiber.Map {
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
