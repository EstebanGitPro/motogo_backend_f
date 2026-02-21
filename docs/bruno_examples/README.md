# Bruno Examples - MotoGo API

Colección de ejemplos para probar la API de MotoGo usando Bruno.

## 📁 Estructura

```
docs/
├── bruno_examples/     # Ejemplos básicos de la API
├── bruno_messages/     # Colección completa de Messages API
└── bruno_persons/      # Colección de Persons API
```

## 🚀 Cómo Usar

### Importar en Bruno

1. Abre Bruno
2. Click en "Open Collection"
3. Selecciona una de las carpetas:
   - `/docs/bruno_examples` - Ejemplos básicos
   - `/docs/bruno_messages` - API de Mensajes completa
   - `/docs/bruno_persons` - API de Personas

### Variables de Entorno

Cada colección tiene configuradas las variables necesarias en `environments/local.bru`:

```
base_url: http://localhost:8080/motogo/api/v1
```

## 📝 Colecciones Disponibles

### bruno_examples
Ejemplos básicos para empezar:
- Create Message
- Get Message by ID
- List Messages

### bruno_messages
Colección completa de la API de Mensajes:
- CRUD completo
- Filtros
- Casos de error
- Tests automáticos

### bruno_persons
API de Personas/Usuarios:
- Create Person
- (más endpoints por agregar)

## ⚙️ Prerequisitos

1. **Docker Desktop** corriendo
2. **MySQL** en puerto 3309
3. **API MotoGo** en puerto 8080

```bash
# Verificar
docker ps
lsof -i :8080
lsof -i :3309
```

## 📚 Más Información

Para documentación detallada de cada API, revisa los archivos dentro de cada colección.
