package util

import "github.com/gofiber/fiber/v2"

func DeepMergeFiberMaps(maps ...fiber.Map) *fiber.Map {
	merged := fiber.Map{}

	var pointer *fiber.Map
	for _, m := range maps {
		for key, value := range m {
			if pointer == nil {
				pointer = &merged
			}

			valueMap, isValueMap := value.(fiber.Map)
			if _, exist := (*pointer)[key]; !exist {
				(*pointer)[key] = fiber.Map{}
			}

			if !isValueMap {
				(*pointer)[key] = value
				pointer = &merged
				continue
			}

			data := (*pointer)[key].(fiber.Map)
			pointer = &data

			deepMergeFiberMaps(pointer, valueMap)

			pointer = &merged
		}
	}

	return &merged
}

func deepMergeFiberMaps(dst *fiber.Map, src interface{}) {
	if dst == nil {
		return
	}

	valueMap, isValueMap := src.(fiber.Map)
	if !isValueMap {
		return
	}

	pointer := dst
	for key, value := range valueMap {
		valueMap, isValueMap := value.(fiber.Map)
		if _, exist := (*pointer)[key]; !exist {
			(*pointer)[key] = fiber.Map{}
		}

		if !isValueMap {
			(*pointer)[key] = value
			pointer = dst
			continue
		}

		data := (*pointer)[key].(fiber.Map)
		pointer = &data

		deepMergeFiberMaps(pointer, valueMap)
		pointer = dst
	}
}
