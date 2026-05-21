# StartTech Application Architecture

## Application Structure
```
starttech-application/
├── frontend/ # React + TypeScript (Vite)
│ ├── src/
│ │ ├── components/ # UI components
│ │ │ ├── ui/ # shadcn/ui primitives
│ │ │ ├── CreateTodo.tsx # Todo creation form
│ │ │ └── TodoItem.tsx # Individual todo display
│ │ ├── context/ # React Context providers
│ │ │ └── AuthContext.tsx # Authentication state
│ │ ├── hooks/ # Custom React hooks
│ │ │ └── useAuth.ts # Auth hook
│ │ ├── lib/ # Utilities
│ │ │ ├── apiClient.ts # Axios API client
│ │ │ └── utils.ts # Helper functions
│ │ ├── routes/ # Page components
│ │ │ ├── login.tsx # Login page
│ │ │ ├── register.tsx # Registration page
│ │ │ ├── todos.tsx # Todo list page
│ │ │ └── profile.tsx # User profile page
│ │ └── types/ # TypeScript type definitions
│ └── vite.config.ts # Vite configuration
│
├── backend/ # Go API (Gin)
│ ├── cmd/api/ # Application entry point
│ ├── internal/
│ │ ├── auth/ # JWT token service
│ │ ├── cache/ # Redis cache layer
│ │ ├── config/ # Configuration loader
│ │ ├── database/ # MongoDB connection
│ │ ├── handlers/ # HTTP handlers
│ │ │ ├── health.go # Health check endpoint
│ │ │ ├── todo.go # CRUD operations for todos
│ │ │ └── user.go # User management
│ │ ├── logger/ # Structured logging
│ │ ├── middleware/ # Gin middleware
│ │ │ ├── auth.go # JWT authentication
│ │ │ ├── cors.go # CORS configuration
│ │ │ └── logger.go # Request logging
│ │ ├── models/ # Data models
│ │ └── routes/ # Route definitions
│ └── Dockerfile # Multi-stage build
│
└── .github/workflows/ # CI/CD pipelines
├── frontend-ci-cd.yml
└── backend-ci-cd.yml
```
## Key Design Decisions

### Authentication Flow
1. User registers/login → Backend validates → JWT token generated
2. Token set as httpOnly cookie (prevents XSS)
3. Frontend uses `withCredentials: true` for all API calls
4. Auth middleware validates token on protected routes

### API Design
- **RESTful**: Standard HTTP methods (GET, POST, PUT, DELETE)
- **Versioning**: Routes grouped under `/auth`, `/tasks`, `/users`
- **Documentation**: Auto-generated Swagger/OpenAPI specs
- **Validation**: Server-side input validation with DTOs

### Caching Strategy
- **Redis**: Optional caching layer (togglable via `ENABLE_CACHE`)
- **Cache Keys**: User-specific keys for todo lists
- **Invalidation**: Cache cleared on todo/user mutations

### Logging Strategy
- **Format**: Structured JSON for machine parsing
- **Levels**: DEBUG, INFO, WARN, ERROR
- **Context**: Request ID, method, path, status code
- **Output**: stdout → CloudWatch Logs via agent