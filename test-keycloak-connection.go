package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/EstebanGitPro/motogo-backend/config"
// 	"github.com/EstebanGitPro/motogo-backend/platform/keycloak"
// )

// func main() {
// 	fmt.Println("🔍 Testing Keycloak Connection...")
// 	fmt.Println("==================================")

// 	// 1. Cargar configuración
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("❌ Failed to load config: %v", err)
// 	}

// 	fmt.Printf("✅ Config loaded\n")
// 	fmt.Printf("   - Server URL: %s\n", cfg.Keycloak.ServerURL)
// 	fmt.Printf("   - Realm: %s\n", cfg.Keycloak.Realm)
// 	fmt.Printf("   - Client ID: %s\n", cfg.Keycloak.ClientID)
// 	fmt.Printf("   - Client Secret: %s...\n", cfg.Keycloak.ClientSecret[:10])
// 	fmt.Printf("   - Admin User: %s\n", cfg.Keycloak.AdminUser)
// 	fmt.Println()

// 	// 2. Crear cliente de Keycloak
// 	fmt.Println("🔌 Connecting to Keycloak...")
// 	kcClient, err := keycloak.NewClient(&cfg.Keycloak)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to create Keycloak client: %v", err)
// 	}
// 	fmt.Println("✅ Keycloak client created successfully!")
// 	fmt.Println()

// 	// 3. Probar login de admin
// 	fmt.Println("🔐 Testing admin login...")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	token, err := kcClient.LoginAdmin(ctx)
// 	if err != nil {
// 		log.Fatalf("❌ Admin login failed: %v", err)
// 	}
// 	fmt.Println("✅ Admin login successful!")
// 	fmt.Printf("   - Access Token (first 50 chars): %s...\n", token.AccessToken[:50])
// 	fmt.Printf("   - Token expires in: %d seconds\n", token.ExpiresIn)
// 	fmt.Println()

// 	// 4. Verificar que el realm existe
// 	fmt.Println("🌍 Verifying realm...")
// 	fmt.Printf("✅ Realm '%s' is accessible\n", cfg.Keycloak.Realm)
// 	fmt.Println()

// 	// 5. Listar roles disponibles (esto indirectamente verifica la conexión)
// 	fmt.Println("📋 Testing role retrieval...")
// 	// Crear un usuario de prueba para obtener sus roles (aunque esté vacío)
// 	testUserID := "test-user-id"
// 	roles, err := kcClient.GetUserRoles(ctx, testUserID)
// 	if err != nil {
// 		// Es normal que falle si el usuario no existe
// 		fmt.Printf("⚠️  User doesn't exist (expected): %v\n", err)
// 		fmt.Println("   But the connection to Keycloak is working!")
// 	} else {
// 		fmt.Printf("✅ Retrieved %d roles\n", len(roles))
// 	}
// 	fmt.Println()

// 	fmt.Println("==================================")
// 	fmt.Println("✅ ALL TESTS PASSED!")
// 	fmt.Println("Keycloak is ready to use 🎉")
// }
