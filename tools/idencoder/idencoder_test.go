package idencoder

import (
	"testing"

	"github.com/google/uuid"
)

func TestIDEncoder_EncodeDecodeUUID(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	})
	if err != nil {
		t.Fatalf("Error creating encoder: %v", err)
	}

	// Generar un UUID de prueba
	testUUID := uuid.New().String()

	// Encodear
	encoded, err := encoder.Encode(testUUID)
	if err != nil {
		t.Fatalf("Error encoding UUID: %v", err)
	}

	// Verificar que el ID ofuscado tiene la longitud mínima
	if len(encoded) < 10 {
		t.Errorf("Encoded ID length is %d, expected at least 10", len(encoded))
	}

	// Decodear
	decoded, err := encoder.Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding ID: %v", err)
	}

	// Verificar que el UUID decodificado es igual al original
	if decoded != testUUID {
		t.Errorf("Decoded UUID %s doesn't match original %s", decoded, testUUID)
	}

	t.Logf("Original UUID: %s", testUUID)
	t.Logf("Encoded ID: %s", encoded)
	t.Logf("Decoded UUID: %s", decoded)
}

func TestIDEncoder_InvalidUUID(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	})
	if err != nil {
		t.Fatalf("Error creating encoder: %v", err)
	}

	// Intentar encodear un UUID inválido
	_, err = encoder.Encode("not-a-valid-uuid")
	if err == nil {
		t.Error("Expected error for invalid UUID, got nil")
	}
}

func TestIDEncoder_InvalidEncodedID(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	})
	if err != nil {
		t.Fatalf("Error creating encoder: %v", err)
	}

	// Intentar decodear un ID inválido
	_, err = encoder.Decode("invalid-encoded-id")
	if err == nil {
		t.Error("Expected error for invalid encoded ID, got nil")
	}
}

func TestIDEncoder_Consistency(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	})
	if err != nil {
		t.Fatalf("Error creating encoder: %v", err)
	}

	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Encodear el mismo UUID múltiples veces
	encoded1, _ := encoder.Encode(testUUID)
	encoded2, _ := encoder.Encode(testUUID)
	encoded3, _ := encoder.Encode(testUUID)

	// Verificar que siempre genera el mismo ID ofuscado
	if encoded1 != encoded2 || encoded2 != encoded3 {
		t.Errorf("Encoding is not consistent: %s, %s, %s", encoded1, encoded2, encoded3)
	}

	t.Logf("UUID %s always encodes to: %s", testUUID, encoded1)
}
