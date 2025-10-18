package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/EstebanGitPro/motogo-backend/cmd/dependency"
// )

// func main() {
// 	fmt.Println("🧪 Testing Full User Registration Flow")
// 	fmt.Println("=======================================")

// 	// Inicializar todas las dependencias
// 	fmt.Println("⚙️  Initializing dependencies...")
// 	deps, err := dependency.InitDependencies()
// 	if err != nil {
// 		log.Fatalf("❌ Failed to initialize dependencies: %v", err)
// 	}
// 	fmt.Println("✅ Dependencies initialized")
// 	fmt.Println()

// 	// Datos de prueba
// 	testEmail := fmt.Sprintf("test-%d@motogo.com", time.Now().Unix())
// 	testPassword := "Test123456!"

// 	fmt.Printf("📝 Test User Details:\n")
// 	fmt.Printf("   - Email: %s\n", testEmail)
// 	fmt.Printf("   - Password: %s\n", testPassword)
// 	fmt.Printf("   - Role: driver\n")
// 	fmt.Println()

// 	// Crear contexto
// 	ctx := context.Background()

// 	// Intentar crear el usuario
// 	fmt.Println("🚀 Attempting to register user...")
// 	fmt.Println("   Step 1: Creating user in local database...")
// 	fmt.Println("   Step 2: Syncing with Keycloak...")
// 	fmt.Println("   Step 3: Setting password in Keycloak...")
// 	fmt.Println("   Step 4: Assigning role in Keycloak...")
// 	fmt.Println("   Step 5: Getting auth token...")
// 	fmt.Println()

// 	// Simular registro a través del service
// 	person := struct {
// 		Email     string
// 		Password  string
// 		FirstName string
// 		LastName  string
// 		Phone     string
// 		Role      string
// 	}{
// 		Email:     testEmail,
// 		Password:  testPassword,
// 		FirstName: "Test",
// 		LastName:  "User",
// 		Phone:     "+1234567890",
// 		Role:      "driver", // Cambiar a tu rol preferido
// 	}

// 	fmt.Printf("⏳ Registering user with role '%s'...\n", person.Role)
// 	fmt.Println()

// 	// Nota: Necesitas tener el rol 'driver' creado en Keycloak primero
// 	fmt.Println("⚠️  IMPORTANT: Make sure the role 'driver' exists in Keycloak!")
// 	fmt.Println("   Current roles in Keycloak: admin, representative, user")
// 	fmt.Println()
// 	fmt.Println("👉 Next steps:")
// 	fmt.Println("   1. Decide which roles to use")
// 	fmt.Println("   2. Create missing roles in Keycloak if needed")
// 	fmt.Println("   3. Test registration through your API")

// 	_ = deps
// 	_ = ctx
// }
