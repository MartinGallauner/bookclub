# CLAUDE.md - Bookclub Project

This file provides guidance to Claude Code (claude.ai/code) when working with the Bookclub project.

## Important: Development Preference

**CRITICAL: Do NOT make any changes to code files.** This includes:

- No code edits or modifications
- No adding TODO comments or any other comments
- No writing new files
- The ONLY files you are allowed to edit are CLAUDE.md files (this file, `frontend/CLAUDE.md`, `backend/CLAUDE.md`) and `*_PLAN.md` files at the project root

**Act as a staff-level software engineer with a passion for teaching.** Instead of writing code:

- **Describe what to do** using clear, conceptual explanations
- **Identify missing concepts** when the user is stuck
- **Explain the "why"** behind architectural decisions
- **Guide through thinking** rather than giving answers
- **Use analogies and comparisons** to clarify complex ideas
- **Point to the right place** in files where changes should be made, but don't write the code

The user wants to learn full-stack development by doing the implementation themselves. Your role is to be a knowledgeable mentor who helps them understand concepts deeply, not a code generator.

---

## Project Overview

**Bookclub** is a full-stack web application that allows users to digitally catalog their physical book collections and make them searchable to friends. Users can scan book barcodes to add books to their personal library, connect with friends, and discover which friends own specific books.

### Learning Goals

This project is designed as a learning vehicle for:

1. **Frontend Development:** Flutter web, responsive UI, state management
2. **Backend Development:** Java, Spring Boot, REST API design
3. **Database Design:** PostgreSQL, schema normalization, migrations
4. **API Design:** OpenAPI-first development, contract-driven architecture
5. **Authentication:** OAuth 2.0 with Google, JWT tokens
6. **Testing:** Unit tests, integration tests, E2E testing
7. **Deployment:** GCP Cloud Run, Cloud SQL, Firebase Hosting
8. **DevOps:** Docker, CI/CD concepts

---

## Architecture Evolution

### Phase 1: Firebase Direct Access (Previous)

```
┌─────────────────────┐
│   Flutter Web UI    │
│                     │
│  Firebase SDK       │
└──────────┬──────────┘
           │ Direct access
           │ (Firestore, Auth)
           ▼
┌─────────────────────┐
│   Firebase/GCP      │
│                     │
│  • Firestore        │
│  • Firebase Auth    │
│  • Security Rules   │
└─────────────────────┘
```

**Characteristics:**
- Client-side Firebase SDK
- Real-time data synchronization
- Security rules in Firestore
- NoSQL document database
- Minimal backend complexity

**Limitations:**
- Complex queries difficult (N+1 problems)
- Limited control over auth flow
- Not suitable for public API
- Harder to learn traditional backend concepts

### Phase 2: REST API Backend (Current)

```
┌─────────────────────┐
│   Flutter Web UI    │
│  (frontend/)        │
│                     │
│  HTTP Client        │
└──────────┬──────────┘
           │
           │ REST API (HTTPS)
           │ OpenAPI Contract
           │
           ▼
┌─────────────────────┐
│  Spring Boot API    │
│  (backend/)         │
│                     │
│  • Controllers      │
│  • Services         │
│  • Repositories     │
│  • JPA/Hibernate    │
└──────────┬──────────┘
           │
           │ JDBC
           │
           ▼
┌─────────────────────┐
│    PostgreSQL       │
│  (Cloud SQL/Local)  │
│                     │
│  • Users            │
│  • Books            │
│  • Connections      │
└─────────────────────┘
```

**Characteristics:**
- Traditional 3-tier architecture
- Request-response pattern (no real-time)
- Server-side OAuth 2.0 + JWT
- Relational database with JOINs
- OpenAPI-first development

**Benefits:**
- Learning opportunity for backend development
- Complex queries with SQL JOINs
- Full control over authentication
- Industry-standard patterns
- Public API ready

---

## System Architecture

### The OpenAPI Contract: Single Source of Truth

The core of this architecture is the **OpenAPI specification** at `openapi/api.yaml`. This single file defines the contract between frontend and backend.

```
                   ┌─────────────────────┐
                   │  openapi/api.yaml   │
                   │                     │
                   │  • Endpoints        │
                   │  • Request schemas  │
                   │  • Response schemas │
                   │  • Validation rules │
                   └──────────┬──────────┘
                              │
                 ┌────────────┴────────────┐
                 │                         │
                 ▼                         ▼
       ┌──────────────────┐      ┌──────────────────┐
       │  Dart Generator  │      │  Java Generator  │
       │  (Frontend)      │      │  (Backend)       │
       └────────┬─────────┘      └────────┬─────────┘
                │                         │
                ▼                         ▼
       ┌──────────────────┐      ┌──────────────────┐
       │  API Client      │      │  API Interfaces  │
       │  Models          │      │  DTOs            │
       │  (Dart classes)  │      │  (Java classes)  │
       └────────┬─────────┘      └────────┬─────────┘
                │                         │
                ▼                         ▼
       ┌──────────────────┐      ┌──────────────────┐
       │  UI Components   │      │  Controllers     │
       │  call client     │      │  implement       │
       │                  │      │  interfaces      │
       └──────────────────┘      └──────────────────┘
```

**Why OpenAPI-First?**

1. **Contract Enforcement:** Frontend and backend must stay in sync
2. **Documentation:** The spec IS the documentation
3. **Type Safety:** Both sides use generated types
4. **Validation:** Request/response validation comes for free
5. **Testing:** Clear contracts make testing easier
6. **Prevents Drift:** Can't deploy mismatched client/server

### API Endpoints (REST Resources)

| Resource              | Endpoint                           | Methods      | Purpose                          |
|-----------------------|------------------------------------|--------------|----------------------------------|
| User Library          | `/api/users/{userId}/library`      | GET, POST    | Manage user's book collection    |
| User Profile          | `/api/users/{userId}`              | GET          | Fetch user profile               |
| Connections           | `/api/connections`                 | GET, POST    | Friend connections               |
| Book Search           | `/api/search/books`                | GET          | Find books in friends' libraries |
| Authentication        | `/api/auth/google`                 | POST         | OAuth 2.0 login                  |

### Authentication Flow (OAuth 2.0 + JWT)

```
┌──────────┐                                    ┌──────────┐
│ Flutter  │                                    │  Google  │
│   App    │                                    │  OAuth   │
└────┬─────┘                                    └────┬─────┘
     │                                                │
     │ 1. Redirect to Google consent screen          │
     │ ──────────────────────────────────────────────>
     │                                                │
     │            2. User approves                    │
     │                                                │
     │ 3. Redirect with authorization code            │
     │ <──────────────────────────────────────────────
     │                                                │
     │                    ┌─────────────┐             │
     │ 4. POST /api/auth  │   Spring    │             │
     │    {code: "..."}   │    Boot     │             │
     │ ──────────────────>│   Backend   │             │
     │                    └──────┬──────┘             │
     │                           │                    │
     │                           │ 5. Exchange code   │
     │                           │    for token       │
     │                           │ ──────────────────>│
     │                           │                    │
     │                           │ 6. Access token +  │
     │                           │    user profile    │
     │                           │ <──────────────────│
     │                           │                    │
     │                           │ 7. Create/update   │
     │                           │    user in         │
     │                           │    PostgreSQL      │
     │                           │                    │
     │ 8. JWT token              │                    │
     │    {token: "...",         │                    │
     │     userId: "..."}        │                    │
     │ <──────────────────       │                    │
     │                           │                    │
     │ 9. All requests:          │                    │
     │    Authorization:         │                    │
     │    Bearer {JWT}           │                    │
     │ ──────────────────────────>                    │
     │                                                │
```

**Key Points:**
- Backend validates Google credentials (not client-side trust)
- JWT is stateless (no session storage)
- JWT contains: user ID, email, expiration
- Expiration: 7 days (configurable)

---

## Project Structure

```
bookclub/
├── CLAUDE.md                    ← You are here (project overview)
│
├── openapi/
│   └── api.yaml                 ← API contract (SINGLE SOURCE OF TRUTH)
│
├── frontend/                    ← Flutter Web Application
│   ├── CLAUDE.md                ← Frontend-specific documentation
│   ├── pubspec.yaml
│   ├── lib/
│   │   ├── main.dart            ← App entry, navigation
│   │   ├── pages/               ← UI pages (Library, Network, Search)
│   │   ├── widgets/             ← Reusable components
│   │   ├── services/            ← Business logic (will use generated client)
│   │   └── generated/api/       ← Generated Dart client (from OpenAPI)
│   ├── test/
│   └── build/web/               ← Production build output
│
└── backend/                     ← Spring Boot REST API
    ├── CLAUDE.md                ← Backend-specific documentation
    ├── build.gradle             ← Gradle build config + OpenAPI plugin
    ├── src/main/
    │   ├── java/com/bookclub/
    │   │   ├── controller/      ← REST controllers (implement generated interfaces)
    │   │   ├── service/         ← Business logic
    │   │   ├── repository/      ← JPA repositories
    │   │   ├── entity/          ← JPA entities (database tables)
    │   │   ├── config/          ← Spring configuration (Security, etc.)
    │   │   └── exception/       ← Custom exceptions
    │   └── resources/
    │       ├── application.yml  ← Spring Boot config
    │       └── db/migration/    ← Flyway SQL migrations
    │           └── V1__initial_schema.sql
    ├── src/test/
    │   └── java/com/bookclub/   ← Unit + integration tests
    └── build/generated/         ← Generated Java interfaces + DTOs (from OpenAPI)
```

---

## Database Design (PostgreSQL)

### Schema: Normalized Relational Design

The database uses a **normalized schema** to avoid data duplication and enable complex queries.

```
┌─────────────────────────────────────────────────────┐
│  users                                              │
├─────────────────────────────────────────────────────┤
│  id (VARCHAR PK)          ← Firebase UID            │
│  email (VARCHAR UNIQUE)                             │
│  display_name (VARCHAR)                             │
│  photo_url (TEXT)                                   │
│  created_at (TIMESTAMP)                             │
│  updated_at (TIMESTAMP)                             │
└──────────────────┬──────────────────────────────────┘
                   │
                   │
       ┌───────────┴───────────┐
       │                       │
       ▼                       ▼
┌─────────────────┐   ┌─────────────────────────────┐
│ user_library    │   │  connections                │
├─────────────────┤   ├─────────────────────────────┤
│ id (UUID PK)    │   │ id (UUID PK)                │
│ user_id (FK)────┼   │ user_id_1 (FK)──────────────┼──┐
│ book_isbn (FK)──┼─┐ │ user_id_2 (FK)──────────────┼──┤
│ added_at        │ │ │ status (ENUM)               │  │
└─────────────────┘ │ │ requested_by (FK)───────────┼──┤
                    │ │ requested_at (TIMESTAMP)    │  │
                    │ │ accepted_at (TIMESTAMP)     │  │
                    │ │ rejected_at (TIMESTAMP)     │  │
                    │ └─────────────────────────────┘  │
                    │                                  │
                    │      ┌───────────────────────────┘
                    │      │
                    │      └──> references users(id)
                    │
                    ▼
┌──────────────────────────────────────────┐
│  books                                   │
├──────────────────────────────────────────┤
│  isbn (VARCHAR PK)                       │
│  title (VARCHAR)                         │
│  publisher (VARCHAR)                     │
│  published_year (INT)                    │
│  language (VARCHAR)                      │
│  cover_url (TEXT)                        │
│  description (TEXT)                      │
└────┬─────────────────────────────────────┘
     │
     ├──> book_authors (isbn, author_name, author_order)
     └──> book_genres (isbn, genre)
```

**Key Design Decisions:**

1. **Users:** Firebase UID as primary key (continuity with old system)
2. **Books:** Separate table (normalized) - book metadata stored once
3. **Authors/Genres:** Separate tables (not arrays) for better queryability
4. **User Library:** Junction table linking users to books
5. **Connections:** Unique constraint on `LEAST(user_id_1, user_id_2)` prevents duplicates

**Why Normalized?**
- One of the main reasons to move from Firestore to SQL
- Teaches proper database design
- Enables complex JOINs (e.g., "which friends have this book?")
- No data duplication

---

## Development Workflow

### 1. API-First Development Cycle

When adding a new feature:

```
┌─────────────────────────────────────────────────────┐
│  Step 1: Design API in openapi/api.yaml            │
│  - Define endpoint                                  │
│  - Define request/response schemas                  │
│  - Add validation rules                             │
└────────────────────┬────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────┐
│  Step 2: Generate Code                              │
│  - Backend: ./gradlew openApiGenerate               │
│  - Frontend: (Dart generator - to be set up)       │
└────────────────────┬────────────────────────────────┘
                     │
           ┌─────────┴─────────┐
           │                   │
           ▼                   ▼
┌────────────────────┐   ┌─────────────────────┐
│  Step 3a: Backend  │   │  Step 3b: Frontend  │
│  - Implement       │   │  - Use generated    │
│    generated       │   │    client           │
│    interface       │   │  - Update UI        │
│  - Write service   │   │  - Update state     │
│  - Write tests     │   │    management       │
└────────────────────┘   └─────────────────────┘
           │                   │
           └─────────┬─────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────┐
│  Step 4: Test Integration                           │
│  - Backend tests with Testcontainers               │
│  - Frontend manual testing                          │
│  - E2E tests (future)                               │
└─────────────────────────────────────────────────────┘
```

### 2. Making Changes to Existing APIs

**IMPORTANT:** Never change the API spec in a breaking way without versioning.

**Non-breaking changes (safe):**
- Adding new endpoints
- Adding optional fields to requests
- Adding fields to responses

**Breaking changes (requires versioning):**
- Removing endpoints
- Removing required fields
- Changing field types
- Renaming fields

### 3. Database Schema Changes

Use Flyway migrations for schema versioning:

```bash
# Create new migration
backend/src/main/resources/db/migration/V2__add_book_ratings.sql

# Flyway automatically runs migrations on app startup
# Migrations are NEVER edited after deployment
```

---

## Local Development Setup

### Prerequisites

- **Java 25** (configured in backend/build.gradle)
- **Flutter SDK** (latest stable)
- **Docker** (for PostgreSQL)
- **PostgreSQL 15+** (local or Docker)
- **Mise** (task runner - optional, for convenience)

### Backend Setup

1. **Start PostgreSQL:**
   ```bash
   docker run --name bookclub-postgres \
     -e POSTGRES_DB=bookclub \
     -e POSTGRES_USER=dev \
     -e POSTGRES_PASSWORD=dev \
     -p 5432:5432 \
     -d postgres:15-alpine
   ```

2. **Set environment variables** (create `backend/.env`):
   ```bash
   DATABASE_URL=jdbc:postgresql://localhost:5432/bookclub
   DATABASE_USER=dev
   DATABASE_PASSWORD=dev
   GOOGLE_CLIENT_ID=your-client-id
   GOOGLE_CLIENT_SECRET=your-client-secret
   JWT_SECRET_KEY=your-secret-key
   ```

3. **Generate code and run:**
   ```bash
   cd backend
   ./gradlew openApiGenerate
   ./gradlew bootRun
   ```

4. **Verify it's running:**
   ```bash
   curl http://localhost:8080/actuator/health
   # Should return: {"status":"UP"}
   ```

### Frontend Setup

1. **Install dependencies:**
   ```bash
   cd frontend
   flutter pub get
   ```

2. **Run locally:**
   ```bash
   flutter run -d chrome
   # Or use mise:
   mise run start
   ```

3. **Access at:** http://localhost:8080

### Full Stack Development

**Terminal 1:** Backend
```bash
cd backend && ./gradlew bootRun
```

**Terminal 2:** Frontend
```bash
cd frontend && flutter run -d chrome
```

**Terminal 3:** API changes
```bash
# Edit openapi/api.yaml
cd backend && ./gradlew openApiGenerate
# (Frontend codegen when set up)
```

---

## Testing Strategy

### Backend Testing Pyramid

```
                    ┌─────────────┐
                    │   E2E Tests │  ← (Future) Full system test
                    │   (Slow)    │
                    └─────────────┘
                   /               \
              ┌────────────────────────┐
              │ Integration Tests      │  ← @SpringBootTest + Testcontainers
              │ (Medium)               │     Real PostgreSQL, HTTP layer
              └────────────────────────┘
             /                          \
    ┌────────────────────────────────────────┐
    │         Unit Tests                     │  ← @WebMvcTest, @DataJpaTest
    │         (Fast)                         │     Mockito, JUnit 5
    └────────────────────────────────────────┘
```

**Run backend tests:**
```bash
cd backend
./gradlew test                           # All tests
./gradlew test --tests LibraryServiceTest  # Specific test
```

### Frontend Testing

```bash
cd frontend
flutter test                    # All tests
flutter test test/widget_test.dart  # Specific test
flutter test --coverage         # With coverage report
```

---

## Deployment Architecture (Google Cloud Platform)

### Production Environment

```
┌──────────────────────────────────────────────────────────┐
│                    Internet (HTTPS)                      │
└────────────────┬─────────────────────┬───────────────────┘
                 │                     │
                 │                     │
        ┌────────▼─────────┐  ┌───────▼────────────┐
        │ Firebase Hosting │  │  Cloud Run         │
        │                  │  │  (Backend API)     │
        │ • Static files   │  │                    │
        │ • CDN            │  │  • Auto-scaling    │
        │ • SSL/TLS        │  │  • Serverless      │
        │                  │  │  • Docker container│
        └──────────────────┘  └─────────┬──────────┘
                                        │
                                        │ Cloud SQL Connector
                                        │
                                ┌───────▼────────────┐
                                │  Cloud SQL         │
                                │  (PostgreSQL)      │
                                │                    │
                                │  • Managed DB      │
                                │  • Backups         │
                                │  • High availability│
                                └────────────────────┘

                          ┌─────────────────────────┐
                          │  Secret Manager         │
                          │  • DB credentials       │
                          │  • JWT signing key      │
                          │  • OAuth secrets        │
                          └─────────────────────────┘
```

### Deployment Commands

**Frontend (Firebase Hosting):**
```bash
cd frontend
mise run deploy  # Runs: analyze → test → build → firebase deploy
```

**Backend (Cloud Run):**
```bash
cd backend
# Build Docker image
docker build -t gcr.io/bookclub-PROJECT_ID/backend:latest .

# Push to Google Container Registry
docker push gcr.io/bookclub-PROJECT_ID/backend:latest

# Deploy to Cloud Run
gcloud run deploy bookclub-backend \
  --image gcr.io/bookclub-PROJECT_ID/backend:latest \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="DATABASE_URL=..." \
  --add-cloudsql-instances=PROJECT_ID:REGION:INSTANCE
```

---

## Key Learning Concepts

### 1. Contract-First Development

The OpenAPI spec is written **before** any implementation code. This is opposite to traditional "code-first" where you write code and document later.

**Benefits:**
- Forces you to think about API design upfront
- Prevents backend/frontend mismatches
- Documentation is always up-to-date
- Both teams can work in parallel

### 2. DTO vs Entity Pattern

**Entity** = Database representation (JPA)
**DTO** = API representation (OpenAPI generated)

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│  Database   │  ←──→ │  Entity     │  ←──→ │  DTO        │  ←──→  Network
│  (SQL)      │  JDBC │  (JPA)      │  Map  │  (JSON)     │  HTTP
└─────────────┘       └─────────────┘       └─────────────┘
```

**Why separate?**
- Entities may have lazy-loaded relationships (can't serialize)
- Database structure may differ from API structure
- API may expose subset of entity fields (security)
- API may combine multiple entities

### 3. Layer Separation

```
Controller  → "HTTP concerns" (routing, validation, status codes)
Service     → "Business logic" (rules, calculations, orchestration)
Repository  → "Data access" (queries, CRUD)
Entity      → "Data structure" (table mapping)
```

**Each layer has one responsibility.**

### 4. Migration Strategy: Big Bang vs Strangler Fig

**Big Bang** (chosen):
- Rewrite everything at once
- High risk, but faster learning
- All-or-nothing deployment

**Strangler Fig** (alternative):
- Migrate one feature at a time
- Lower risk, but more complex
- Two systems running simultaneously

---

## Success Metrics

### Backend Ready When:

- ✅ OpenAPI spec defines all endpoints
- ✅ Controllers implement generated interfaces
- ✅ Services contain business logic
- ✅ Repositories query PostgreSQL
- ✅ Unit tests pass (services, repositories)
- ✅ Integration tests pass (Testcontainers)
- ✅ OAuth flow works with Google
- ✅ JWT tokens are generated/validated
- ✅ Flyway migrations run successfully
- ✅ Docker image builds and runs
- ✅ Deployed to Cloud Run
- ✅ Health check endpoint responds

### Frontend Ready When:

- ✅ Dart client generated from OpenAPI
- ✅ Services use generated client (not Firebase SDK)
- ✅ OAuth flow redirects to backend
- ✅ JWT tokens stored and sent with requests
- ✅ All pages fetch from REST API
- ✅ Error handling for API failures
- ✅ Loading states during API calls
- ✅ Widget tests pass
- ✅ Build succeeds for web
- ✅ Deployed to Firebase Hosting

### Integration Ready When:

- ✅ Frontend can authenticate via backend
- ✅ Frontend can CRUD books via API
- ✅ Frontend can manage connections via API
- ✅ Frontend can search books via API
- ✅ E2E user flows work
- ✅ Error handling works across stack
- ✅ Performance is acceptable (<500ms API responses)

---

## Common Pitfalls

### 1. Breaking API Changes

**Problem:** Changing OpenAPI spec breaks deployed frontend

**Solution:**
- Always add fields, never remove (deprecate instead)
- Use API versioning (`/api/v1/`, `/api/v2/`)
- Test frontend compatibility before deploying backend

### 2. Forgetting to Regenerate Code

**Problem:** Edit OpenAPI spec but forget to run code generators

**Solution:**
- Make codegen part of build process (already configured)
- `compileJava.dependsOn tasks.openApiGenerate` ensures this

### 3. Exposing Entities in Controllers

**Problem:** Returning JPA entities directly from REST endpoints

**Solution:**
- Always map entities to DTOs
- DTOs are generated from OpenAPI
- Controllers convert: Entity → DTO → JSON

### 4. Missing Database Indexes

**Problem:** Queries slow down as data grows

**Solution:**
- Add indexes to foreign keys
- Add indexes to WHERE/JOIN columns
- Monitor query performance

### 5. Hardcoding Configuration

**Problem:** Database URLs, secrets in code

**Solution:**
- Use environment variables
- Use GCP Secret Manager for production
- Never commit secrets to Git

---

## Component Documentation

For detailed component-specific documentation:

- **Frontend Details:** See `frontend/CLAUDE.md`
  - Flutter widgets
  - State management
  - Routing
  - Testing

- **Backend Details:** See `backend/CLAUDE.md`
  - Spring Boot configuration
  - JPA entities
  - Service layer patterns
  - Testing strategies
  - Database schema design

- **API Contract:** See `openapi/api.yaml`
  - Endpoint definitions
  - Request/response schemas
  - Validation rules

---

## Next Steps

### Phase 1: OpenAPI Spec (Current)

1. Define all endpoints in `openapi/api.yaml`
2. Define all request/response schemas
3. Add validation rules (min/max length, patterns, required fields)

### Phase 2: Backend Implementation

1. Generate Java code: `./gradlew openApiGenerate`
2. Implement database schema (Flyway migration)
3. Create JPA entities
4. Create repositories
5. Implement services (business logic)
6. Implement controllers (implement generated interfaces)
7. Write unit tests
8. Write integration tests (Testcontainers)
9. Set up OAuth 2.0 + JWT
10. Test with Postman/curl

### Phase 3: Frontend Integration

1. Generate Dart client from OpenAPI
2. Replace Firebase SDK with generated client
3. Update authentication to use backend OAuth
4. Update all services to call REST API
5. Handle loading/error states
6. Write widget tests
7. Manual E2E testing

### Phase 4: Deployment

1. Build Docker image for backend
2. Deploy backend to Cloud Run
3. Set up Cloud SQL for PostgreSQL
4. Configure Secret Manager
5. Build Flutter web app
6. Deploy frontend to Firebase Hosting
7. Configure CORS
8. Test production environment

### Phase 5: Enhancements (Future)

- Pagination for large lists
- Book cover image upload (Cloud Storage)
- Full-text search (PostgreSQL tsvector)
- Rate limiting
- Caching (Redis)
- Monitoring (Cloud Monitoring)
- CI/CD pipeline (Cloud Build)

---

## Resources

### Documentation

- **Spring Boot:** https://spring.io/projects/spring-boot
- **OpenAPI Specification:** https://swagger.io/specification/
- **Flutter:** https://flutter.dev/docs
- **PostgreSQL:** https://www.postgresql.org/docs/
- **GCP Cloud Run:** https://cloud.google.com/run/docs

### Tools

- **OpenAPI Generator:** https://openapi-generator.tech/
- **Testcontainers:** https://www.testcontainers.org/
- **Flyway:** https://flywaydb.org/
- **Postman:** https://www.postman.com/ (API testing)

---

*This document is a living guide. Update it as architectural decisions evolve.*