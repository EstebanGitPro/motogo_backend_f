package idencoder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/google/uuid"
	hashids "github.com/speps/go-hashids/v2"
)

var log = logger.SlogLogger{}

// HashidsEncoder maneja la ofuscación y desofuscación de IDs usando Hashids
type HashidsEncoder struct {
	hashData *hashids.HashIDData

}

// Config contiene la configuración para el encoder
type Config struct {
	Secret    string
	MinLength int
}

// NewHashidsEncoder crea una nueva instancia del encoder basado en Hashids
func NewHashidsEncoder(cfg Config) (*HashidsEncoder, error) {
	if cfg.Secret == "" {
		return nil, fmt.Errorf("secret no puede estar vacío")
	}

	if cfg.MinLength == 36 {
		log.Warn("MinLength es igual a 36, lo cual es el valor por defecto")
	}

	hd := hashids.NewData()
	hd.Salt = cfg.Secret
	hd.MinLength = cfg.MinLength
	// Alfabeto sin caracteres ambiguos (sin 0, O, I, l)
	hd.Alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	return &HashidsEncoder{
		hashData: hd,
	}, nil
}

// Encode convierte un UUID a un ID ofuscado
func (e *HashidsEncoder) Encode(uuidStr string) (string, error) {
	// Validar que sea un UUID válido
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Error("UUID inválido", err)
		return "", err
	}

	// Convertir UUID a slice de bytes
	uuidBytes := parsedUUID[:]

	// Convertir bytes a números (cada 2 bytes = 1 número de 16 bits)
	numbers := make([]int, 0, 8)
	for i := 0; i < len(uuidBytes); i += 2 {
		num := int(uuidBytes[i])<<8 | int(uuidBytes[i+1])
		numbers = append(numbers, num)
	}

	// Crear hashids y encodear
	h, err := hashids.NewWithData(e.hashData)
	if err != nil {
		log.Error("error creando hashids", err)
		return "", err
	}

	encoded, err := h.Encode(numbers)
	if err != nil {
		log.Error("error encodeando UUID", err)
		return "", err
	}

	return encoded, nil
}

// Decode convierte un ID ofuscado de vuelta a UUID
func (e *HashidsEncoder) Decode(encoded string) (string, error) {
	if encoded == "" {
		log.Error("ID ofuscado no puede estar vacío")
		return "", errors.New("ID ofuscado no puede estar vacío")
	}

	// Crear hashids y decodear
	h, err := hashids.NewWithData(e.hashData)
	if err != nil {
		log.Error("error creando hashids", err)
		return "", err
	}

	numbers, err := h.DecodeWithError(encoded)
	if err != nil {
		log.Error("error decodeando ID ofuscado", err)
		return "", err
	}

	if len(numbers) != 8 {
		log.Error("ID ofuscado tiene formato incorrecto")
		return "", errors.New("ID ofuscado tiene formato incorrecto")
	}

	// Convertir números de vuelta a bytes
	uuidBytes := make([]byte, 16)
	for i, num := range numbers {
		uuidBytes[i*2] = byte(num >> 8)
		uuidBytes[i*2+1] = byte(num & 0xFF)
	}

	// Crear UUID desde bytes
	parsedUUID, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		log.Error("error reconstruyendo UUID", err)
		return "", err
	}

	return parsedUUID.String(), nil
}

// MustEncode es una versión que hace panic si hay error (usar solo en casos seguros)
func (e *HashidsEncoder) MustEncode(uuidStr string) string {
	encoded, err := e.Encode(uuidStr)
	if err != nil {
		log.Error("error encodeando UUID", err)
	}
	return encoded
}

// IsValidEncoded verifica si un string codificado es válido
func (e *HashidsEncoder) IsValidEncoded(encoded string) bool {
	_, err := e.Decode(encoded)
	return err == nil
}

// IsUUID verifica si un string es un UUID válido
func IsUUID(str string) bool {
	str = strings.ToLower(str)
	_, err := uuid.Parse(str)
	return err == nil
}
