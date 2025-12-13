package idencoder

import (
	"testing"

	"github.com/google/uuid"
)

func TestIDEncoder_EncodeDecodeUUID(t *testing.T) {
	encoder, err := NewHashidsEncoder(Config{
		Secret:    "test-secret-12345",
		MinLength: 10,
	}, nil)
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
	}, nil)
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
	}, nil)
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
	}, nil)
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

func TestIDEncoder_EmptySecret(t *testing.T) {
	_, err := NewHashidsEncoder(Config{
		Secret:    "",
		MinLength: 10,
	}, nil)

	if err == nil {
		t.Error("Expected error for empty secret, got nil")
	}
}

func TestIDEncoder_DifferentSecrets_DifferentEncodings(t *testing.T) {
	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	encoder1, _ := NewHashidsEncoder(Config{
		Secret:    "secret-1",
		MinLength: 10,
	}, nil)

	encoder2, _ := NewHashidsEncoder(Config{
		Secret:    "secret-2",
		MinLength: 10,
	}, nil)

	encoded1, _ := encoder1.Encode(testUUID)
	encoded2, _ := encoder2.Encode(testUUID)

	if encoded1 == encoded2 {
		t.Error("Different secrets should produce different encodings")
	}
}

func TestIDEncoder_DecodeWithDifferentSecret_Fails(t *testing.T) {
	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	encoder1, _ := NewHashidsEncoder(Config{
		Secret:    "secret-1",
		MinLength: 10,
	}, nil)

	encoder2, _ := NewHashidsEncoder(Config{
		Secret:    "secret-2",
		MinLength: 10,
	}, nil)

	encoded, _ := encoder1.Encode(testUUID)

	// Try to decode with different encoder
	decoded, err := encoder2.Decode(encoded)

	// Should either error or produce wrong result
	if err == nil && decoded == testUUID {
		t.Error("Decode with different secret should fail or produce wrong result")
	}
}

func TestIsUUID_ValidFormats(t *testing.T) {
	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"F47AC10B-58CC-4372-A567-0E02B2C3D479",
		"f47ac10b-58cc-4372-a567-0e02b2c3d479",
	}

	for _, uuid := range validUUIDs {
		if !IsUUID(uuid) {
			t.Errorf("Expected %s to be valid UUID", uuid)
		}
	}
}

func TestIsUUID_InvalidFormats(t *testing.T) {
	invalidValues := []string{
		"not-a-uuid",
		"",
		"12345",
		"550e8400-e29b-41d4-a716",
	}

	for _, val := range invalidValues {
		if IsUUID(val) {
			t.Errorf("Expected %s to be invalid UUID", val)
		}
	}
}

func TestIDEncoder_MustEncode_WithError(t *testing.T) {
	encoder, _ := NewHashidsEncoder(Config{
		Secret:    "test-secret",
		MinLength: 10,
	}, nil)

	// MustEncode with invalid UUID should return empty string (doesn't panic)
	result := encoder.MustEncode("invalid-uuid")
	if result != "" {
		t.Logf("MustEncode with error returned: %s", result)
	}
}

func TestIDEncoder_IsValidEncoded(t *testing.T) {
	encoder, _ := NewHashidsEncoder(Config{
		Secret:    "test-secret",
		MinLength: 10,
	}, nil)

	testUUID := "550e8400-e29b-41d4-a716-446655440000"
	encoded, _ := encoder.Encode(testUUID)

	// Valid encoded ID
	if !encoder.IsValidEncoded(encoded) {
		t.Error("Valid encoded ID should return true")
	}

	// Invalid encoded ID
	if encoder.IsValidEncoded("definitely-invalid-id") {
		t.Error("Invalid encoded ID should return false")
	}

	// Empty string
	if encoder.IsValidEncoded("") {
		t.Error("Empty string should return false")
	}
}

func TestIDEncoder_MinLength(t *testing.T) {
	testUUID := "550e8400-e29b-41d4-a716-446655440000"

	testCases := []int{5, 10, 20, 30}

	for _, minLen := range testCases {
		encoder, _ := NewHashidsEncoder(Config{
			Secret:    "test-secret",
			MinLength: minLen,
		}, nil)

		encoded, err := encoder.Encode(testUUID)
		if err != nil {
			t.Fatalf("Error encoding with minLength %d: %v", minLen, err)
		}

		if len(encoded) < minLen {
			t.Errorf("Expected encoded length >= %d, got %d", minLen, len(encoded))
		}

		t.Logf("MinLength %d: encoded length = %d", minLen, len(encoded))
	}
}
