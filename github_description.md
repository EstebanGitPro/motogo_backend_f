# MotoGo Backend

This project is a Go-based backend service for the MotoGo platform, built with a hexagonal architecture to ensure clean separation of concerns and maintainability.

## Key Features:
*   **User Management:** Robust APIs for user registration, authentication, email verification, and password reset flows.
*   **Identity Provider Integration:** Seamless integration with Keycloak for secure identity and access management.
*   **Dynamic Messaging:** Administrative APIs for managing localizable system messages, enhancing user communication.
*   **Comprehensive Observability:**
    *   **Metrics:** Prometheus integration via a `/metrics` endpoint for real-time monitoring.
    *   **Logging:** Structured logging using `slog` and Loki for centralized log management.
    *   **Tracing:** Request ID middleware for end-to-end request tracing.
    *   **Dashboarding:** Grafana for visualizing metrics and logs.
*   **API Documentation:** Automatically generated OpenAPI (Swagger) specifications for clear API understanding.
*   **Efficient Caching:** In-memory caching for system messages with automatic refresh, reducing database load.
*   **Load Testing:** Includes k6 scripts for performance and soak testing.

## Technology Stack:
*   **Backend:** Go (Gin Web Framework)
*   **Database:** MySQL
*   **Identity Management:** Keycloak
*   **Monitoring:** Prometheus, Grafana
*   **Logging:** Loki, `slog`
*   **Load Testing:** k6
*   **Architecture:** Hexagonal (Ports & Adapters)

## Getting Started:
(Further instructions on setup and running the project would go here)
