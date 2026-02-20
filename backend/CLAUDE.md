# CLAUDE.md - Backend

This file provides guidance to Claude Code (claude.ai/code) when working with the backend code in this repository.

## Important: Development Preference

**CRITICAL: Do NOT make any changes to code files.** This includes:

- No code edits or modifications
- No adding TODO comments or any other comments
- No writing new files
- The ONLY file you are allowed to edit is this CLAUDE.md file

**Act as a staff-level software engineer with a passion for teaching.** Instead of writing code:

- **Describe what to do** using clear, conceptual explanations
- **Identify missing concepts** when the user is stuck
- **Explain the "why"** behind architectural decisions
- **Guide through thinking** rather than giving answers
- **Use analogies and comparisons** to clarify complex ideas
- **Point to the right place** in files where changes should be made, but don't write the code

The user wants to learn Java and Spring Boot by doing the implementation themselves. Your role is to be a knowledgeable
mentor who helps them understand concepts deeply, not a code generator.

---

This file also documents the architectural design decisions for the Bookclub Java backend.

## Project Overview

This is a **Spring Boot REST API backend** for the Bookclub application. It replaces the previous Firebase/Firestore
direct-access architecture with a traditional 3-tier application stack.

### Migration Context

**Previous Architecture:**

- Flutter app used Firebase Client SDK to directly access Firestore
- Real-time data synchronization via Firestore streams
- Firebase Auth for authentication
- Client-side security rules

**New Architecture:**

- Flutter app calls REST API endpoints (HTTP client)
- PostgreSQL relational database
- Java backend handles authentication (OAuth 2.0 with Google)
- Server-side authorization logic

## Design Decisions

### 1. Database: PostgreSQL

**Decision:** Use PostgreSQL for persistent storage

**Rationale:**

- Need for complex join queries (e.g., searching books across friends' libraries)
- Better query optimization compared to Firestore's N+1 pattern in `SearchService.searchByISBN()`
- Learning opportunity: Understanding relational database design
- Strong data consistency and ACID transactions
- Public API ready: SQL databases are industry standard for APIs

**Key Schema Challenges:**

- **Subcollections → Foreign Keys:** Firestore's `users/{uid}/library/{isbn}` becomes foreign key relationships
- **Arrays:** Firestore stores `authors: ['Author 1', 'Author 2']`. PostgreSQL options:
    - TEXT[] array columns (simpler)
    - Normalized `book_authors` table (more queryable)
- **Document IDs:** Firestore uses composite IDs like `{uid1}_{uid2}` for connections. PostgreSQL uses UUIDs with unique
  constraints.

### 2. Real-Time Updates: NOT Implemented

**Decision:** No real-time data synchronization

**Rationale:**

- REST API uses request-response pattern, not streaming
- Simplifies architecture significantly
- User-initiated refresh pattern (pull-to-refresh in Flutter)
- WebSockets/Server-Sent Events add complexity not needed for MVP

**Flutter Impact:**

- Services no longer extend `ChangeNotifier` with stream subscriptions
- Manual refresh: Call API endpoints when data is needed
- `notifyListeners()` called after successful HTTP responses

**Future Consideration:** If real-time becomes critical, add WebSocket layer later

### 3. Authentication: Java-Managed OAuth 2.0

**Decision:** Java backend handles complete OAuth flow and issues JWT tokens

**Rationale:**

- Learning opportunity: Understand OAuth 2.0 flow end-to-end
- Full control over user session management
- Public API ready: Standard JWT-based auth
- Security: Validate Google credentials server-side, not client-side

**Flow:**

1. Flutter redirects user to Google OAuth consent screen
2. User approves, Google redirects back with authorization code
3. Flutter sends code to backend: `POST /api/auth/google`
4. Backend exchanges code for Google access token (server-to-server call)
5. Backend fetches user profile from Google
6. Backend creates/updates user in PostgreSQL
7. Backend generates JWT token signed with secret key
8. Backend returns JWT to Flutter
9. Flutter includes JWT in `Authorization: Bearer {token}` header for all subsequent requests

**Components:**

- **Spring Security OAuth2 Client:** Handles Google OAuth integration
- **JWT Library:** Generate/validate tokens (e.g., `java-jwt` or Spring Security OAuth)
- **SecurityFilterChain:** Define public (`/api/auth/**`) vs protected endpoints
- **JWT Filter:** Intercept requests, validate token, populate Spring Security context

**Token Strategy:**

- Stateless JWTs (server doesn't store sessions)
- Include claims: user ID, email, expiration
- Expiration: 7 days (configurable)
- Refresh strategy: Re-authenticate with Google (simpler than refresh tokens for MVP)

### 4. Deployment: Google Cloud Platform (GCP)

**Decision:** Deploy to Google Cloud Platform

**Target Service:** Cloud Run (serverless containers)

**Rationale:**

- Simplest GCP deployment: Build Docker image → Push to GCR → Deploy
- Auto-scaling (including scale-to-zero for cost savings)
- Automatic HTTPS
- Pay-per-use pricing
- Easier than GKE (Kubernetes) for single application

**Alternative Considered:**

- **Google Kubernetes Engine (GKE):** Overkill for one app, but valuable for learning K8s
- **App Engine Standard:** Simpler but less flexible

**Deployment Workflow:**

1. Package Spring Boot as executable JAR
2. Create Dockerfile (using multi-stage build)
3. Build Docker image
4. Push to Google Container Registry
5. Deploy to Cloud Run via `gcloud run deploy`

**Infrastructure:**

- **Cloud SQL for PostgreSQL:** Managed PostgreSQL instance
- **Secret Manager:** Store database credentials, JWT signing key, OAuth client secrets
- **Cloud Build:** CI/CD pipeline (optional, can start with local builds)

### 5. API Design: OpenAPI-First with Code Generation

**Decision:** Write OpenAPI 3.0 spec first, generate Java server interfaces and Dart client code

**Rationale:**

- **Contract-first development:** API contract is source of truth
- **Prevents drift:** Java and Flutter always in sync with spec
- **Documentation:** OpenAPI spec IS the documentation
- **Validation:** Generated code includes Bean Validation annotations
- **Type safety:** Both sides use generated models

**Tooling:**

- **Java:** `openapi-generator-maven-plugin` in `pom.xml`
    - Generates Spring Boot controller interfaces
    - Generates model POJOs with Jackson annotations
    - Controllers implement generated interfaces
- **Dart/Flutter:** `openapi-generator-cli` or Flutter package
    - Generates HTTP client with all API methods
    - Generates Dart model classes with serialization

**Spec Location:** `backend/src/main/resources/openapi.yaml`

**Generated Artifacts (NOT hand-edited):**

- Java: `target/generated-sources/openapi/`
- Dart: `frontend/lib/generated/api/` (or similar)

**API Design Pattern:** RESTful resource-based routing

**Endpoint Mapping:**

| Flutter Service Method                 | REST Endpoint                      | HTTP Method | Description                      |
|----------------------------------------|------------------------------------|-------------|----------------------------------|
| `LibraryService.addBook(book)`         | `/api/users/{userId}/library`      | POST        | Add book to user's library       |
| `LibraryService.books` getter          | `/api/users/{userId}/library`      | GET         | List all books in user's library |
| `ContactService.addContact(userToAdd)` | `/api/connections`                 | POST        | Create connection request        |
| `ContactService.connections` getter    | `/api/connections?userId={userId}` | GET         | List user's connections          |
| `SearchService.searchByISBN(isbn)`     | `/api/search/books?isbn={isbn}`    | GET         | Find which friends have a book   |
| `AuthService.signInWithGoogle()`       | `/api/auth/google`                 | POST        | Exchange Google code for JWT     |
| `AuthService.fetchProfile(uid)`        | `/api/users/{userId}`              | GET         | Get user profile                 |

**Error Handling:**

- Use standard HTTP status codes (400, 401, 403, 404, 500, etc.)
- Consistent error response schema in OpenAPI spec:
  ```json
  {
    "timestamp": "2025-02-20T10:30:00Z",
    "status": 400,
    "error": "Bad Request",
    "message": "ISBN is required",
    "path": "/api/users/123/library"
  }
  ```

### 6. Migration Strategy: Big Bang

**Decision:** Rewrite entire backend at once, switch Flutter app in single deployment

**Rationale:**

- **Learning focus:** Want to dive into Java immediately, not maintain two backends
- **Scope:** Application is small enough for complete rewrite
- **Clean break:** No hybrid complexity of Firestore + PostgreSQL

**Risks:**

- High-risk deployment (all-or-nothing)
- Testing burden is higher (must verify all features work)

**Mitigation:**

- Comprehensive integration tests with Testcontainers
- Local testing with production-like environment
- Firebase remains available as rollback option initially

**Alternative (Not Chosen):** Strangler Fig Pattern

- Migrate one service at a time
- Lower risk but slower learning
- More complex: maintaining two backends simultaneously

### 7. Testing Strategy

**Decision:** Three-tier testing pyramid with JUnit 5, Mockito, and Testcontainers

#### Level 1: Unit Tests (Fast, Isolated)

**Purpose:** Test business logic in service layer

**Tools:**

- JUnit 5
- Mockito (mock repositories and external dependencies)
- AssertJ (fluent assertions)

**Pattern:**

```java
@ExtendWith(MockitoExtension.class)
class LibraryServiceTest {
    @Mock
    private LibraryRepository libraryRepository;

    @InjectMocks
    private LibraryService libraryService;

    @Test
    void addBook_savesBookToRepository() {
        // Test service logic without database
    }
}
```

**What to test:**

- Service methods with mocked repositories
- Business logic validation
- Error handling (e.g., duplicate book, unauthorized access)

**What NOT to test:**

- Framework behavior (Spring's DI, JPA's queries)
- Database interactions (use slice tests for that)

#### Level 2: Slice Tests (Medium Speed, Spring Context)

**Purpose:** Test one layer of Spring Boot in isolation

**A. Web Layer (`@WebMvcTest`)**

**Tests:** Controller request/response mapping

**Pattern:**

```java
@WebMvcTest(LibraryController.class)
class LibraryControllerTest {
    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private LibraryService libraryService;

    @Test
    void addBook_returnsCreatedStatus() throws Exception {
        mockMvc.perform(post("/api/users/{userId}/library", "123")
            .contentType(MediaType.APPLICATION_JSON)
            .content("{ \"isbn\": \"123456\" }"))
            .andExpect(status().isCreated());
    }
}
```

**What to test:**

- HTTP status codes
- Request/response JSON mapping
- Validation errors
- Controller logic (if any, though should be minimal)

**B. Data Layer (`@DataJpaTest`)**

**Tests:** JPA repositories and queries

**Pattern:**

```java
@DataJpaTest
class LibraryRepositoryTest {
    @Autowired
    private LibraryRepository libraryRepository;

    @Autowired
    private TestEntityManager entityManager;

    @Test
    void findByIsbn_returnsBook() {
        // Test queries against in-memory H2 database
    }
}
```

**What to test:**

- Custom query methods
- JPA relationships
- Database constraints

**Database:** Uses H2 in-memory by default (fast, but not production-identical)

#### Level 3: Integration Tests (`@SpringBootTest` + Testcontainers)

**Purpose:** Test full application stack with real PostgreSQL

**Tools:**

- `@SpringBootTest` (loads entire Spring context)
- Testcontainers (spins up PostgreSQL in Docker)
- RestAssured or TestRestTemplate (HTTP client for tests)

**Pattern:**

```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@Testcontainers
class LibraryIntegrationTest {
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine")
        .withDatabaseName("testdb")
        .withUsername("test")
        .withPassword("test");

    @LocalServerPort
    private int port;

    @Autowired
    private TestRestTemplate restTemplate;

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Test
    void addBook_persistsToDatabase() {
        // Test entire flow: HTTP → Controller → Service → Repository → PostgreSQL
    }
}
```

**What to test:**

- End-to-end API workflows
- Database transactions
- PostgreSQL-specific features (e.g., JSON columns, full-text search)
- Authentication/authorization flows

**Why Testcontainers?**

- **Production parity:** Tests run against actual PostgreSQL, not H2
- **Catches bugs:** PostgreSQL constraints, data types, and query syntax differ from H2
- **Confidence:** If integration tests pass, deployment is likely to succeed

**Trade-off:** Slower than unit/slice tests (Docker startup overhead)

---

## Spring Boot Architecture

### Layer Separation (Separation of Concerns)

```
┌─────────────────────────────────────────────────────────┐
│  Controller Layer (@RestController)                     │
│  - Implements OpenAPI-generated interfaces              │
│  - HTTP request/response handling                       │
│  - Input validation (Bean Validation)                   │
│  - NO business logic                                    │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Service Layer (@Service)                               │
│  - Business logic and rules                             │
│  - Transaction boundaries (@Transactional)              │
│  - Coordinates multiple repositories                    │
│  - Authorization checks                                 │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Repository Layer (extends JpaRepository)               │
│  - Database queries                                     │
│  - CRUD operations                                      │
│  - Custom query methods                                 │
│  - NO business logic                                    │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│  Entity Layer (@Entity)                                 │
│  - JPA entity mapping to database tables                │
│  - Relationship definitions (@ManyToOne, etc.)          │
│  - Just data structure                                  │
└─────────────────────────────────────────────────────────┘
```

### Key Principles

**Single Responsibility:**

- Each layer has one reason to change
- Database changes → affect entities/repositories only
- Business rules → affect services only
- API contract → affect controllers and OpenAPI spec only

**DTOs vs Entities:**

- **Entity:** Database row mapping (annotated with `@Entity`)
- **DTO (Data Transfer Object):** API contract model (generated from OpenAPI)
- **NEVER expose entities in API responses** (breaks encapsulation, leaks DB structure)
- Controllers map between DTOs and entities

---

## Data Model Design

### Firebase/Firestore → PostgreSQL Schema Mapping

#### Firestore Collections (Current)

```
users/{uid}
├── Fields: uid, email, displayName, photoURL, createdAt, updatedAt
└── library/{isbn}
    └── Fields: title, authors[], isbn, coverUrl, language, publisher,
                publishedIn, genre[], description, addedAt

connections/{uid1}_{uid2}
└── Fields: uid, users[], requestedBy, requestedAt, status,
            acceptedAt, rejectedAt
```

#### PostgreSQL Tables (Design Decisions Needed)

The schema design requires several key decisions:

##### Decision Point 1: Users Table

**Primary Key Choice:**

- **Option A:** Use Firebase UID as primary key (VARCHAR)
    - Pros: Continuity with old system, easier migration
    - Cons: Less conventional, string comparison slower than UUID
- **Option B:** Generate new UUIDs (UUID type)
    - Pros: Database independence, better performance
    - Cons: Need to migrate user IDs

**Recommended:** Option A for initial migration (use Firebase UID), can migrate to UUIDs later if needed

**Schema:**

```sql
CREATE TABLE users (
    id VARCHAR(128) PRIMARY KEY,  -- Firebase UID
    email VARCHAR(255) NOT NULL UNIQUE,
    display_name VARCHAR(255),
    photo_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

##### Decision Point 2: Books and Library

**Critical Design Choice:**

**Option A: Normalized Schema (Separate Books Table)**

Books table stores unique books by ISBN:

```sql
CREATE TABLE books (
    isbn VARCHAR(13) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    publisher VARCHAR(255),
    published_year INTEGER,
    language VARCHAR(10),
    cover_url TEXT,
    description TEXT
);

CREATE TABLE book_authors (
    book_isbn VARCHAR(13) REFERENCES books(isbn) ON DELETE CASCADE,
    author_name VARCHAR(255) NOT NULL,
    author_order INTEGER NOT NULL,
    PRIMARY KEY (book_isbn, author_order)
);

CREATE TABLE book_genres (
    book_isbn VARCHAR(13) REFERENCES books(isbn) ON DELETE CASCADE,
    genre VARCHAR(100) NOT NULL,
    PRIMARY KEY (book_isbn, genre)
);

CREATE TABLE user_library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) REFERENCES users(id) ON DELETE CASCADE,
    book_isbn VARCHAR(13) REFERENCES books(isbn) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, book_isbn)
);

CREATE INDEX idx_user_library_user ON user_library(user_id);
CREATE INDEX idx_user_library_isbn ON user_library(book_isbn);
```

**Pros:**

- No data duplication (book metadata stored once)
- Easy to update book info globally
- Query "which users have this book" is simple
- Normalized design (database best practice)

**Cons:**

- Requires JOINs to get user's library with book details
- More complex JPA entity relationships
- More tables to manage

**Option B: Denormalized Schema (All-in-One Library Table)**

```sql
CREATE TABLE library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) REFERENCES users(id) ON DELETE CASCADE,
    isbn VARCHAR(13) NOT NULL,
    title VARCHAR(500) NOT NULL,
    authors TEXT[] NOT NULL,  -- PostgreSQL array type
    genres TEXT[],
    publisher VARCHAR(255),
    published_year INTEGER,
    language VARCHAR(10),
    cover_url TEXT,
    description TEXT,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, isbn)
);

CREATE INDEX idx_library_user ON library(user_id);
CREATE INDEX idx_library_isbn ON library(isbn);
```

**Pros:**

- Simple queries (no JOINs needed)
- Mirrors Firestore subcollection structure
- Faster reads (single table scan)
- Easier JPA mapping

**Cons:**

- Book metadata duplicated per user
- Updating book info requires updating all entries
- Denormalized (violates database normal forms)

**Recommended:** **Option A (Normalized)** for learning purposes

- **Rationale:** One of the key reasons for moving to SQL is to leverage relational design
- Teaches proper normalization
- More representative of real-world backend development
- Can always denormalize later if performance requires it

##### Decision Point 3: Connections Table

**Schema:**

```sql
CREATE TABLE connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id_1 VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
    requested_by VARCHAR(128) NOT NULL REFERENCES users(id),
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    rejected_at TIMESTAMP,
    -- Ensure no duplicate connections (regardless of order)
    CONSTRAINT unique_connection UNIQUE (
        LEAST(user_id_1, user_id_2),
        GREATEST(user_id_1, user_id_2)
    )
);

CREATE INDEX idx_connections_user1 ON connections(user_id_1);
CREATE INDEX idx_connections_user2 ON connections(user_id_2);
CREATE INDEX idx_connections_status ON connections(status);
```

**Key Features:**

- **Unique constraint:** Uses `LEAST`/`GREATEST` to prevent duplicate connections regardless of user order
- **Composite indexes:** Fast lookups for "all connections for user X"
- **Status check constraint:** Enforces valid enum values at database level
- **Nullable timestamps:** `accepted_at` and `rejected_at` only set when applicable

**Query Pattern for "User X's Connections":**

```sql
SELECT * FROM connections
WHERE (user_id_1 = ? OR user_id_2 = ?)
  AND status = 'accepted';
```

---

## Database Migrations with Flyway

**Decision:** Use Flyway for schema versioning

**Rationale:**

- Version-controlled schema changes
- Automated migration on application startup
- Safe production deployments (can rollback if needed)
- Audit trail of all schema changes

**Location:** `src/main/resources/db/migration/`

**Naming Convention:**

- `V1__initial_schema.sql` (two underscores after version)
- `V2__add_book_ratings.sql`
- `V3__add_indexes.sql`

**Spring Boot Integration:**

```yaml
# application.yml
spring:
  flyway:
    enabled: true
    baseline-on-migrate: true
    locations: classpath:db/migration
```

**First Migration (V1__initial_schema.sql):**
Should contain all CREATE TABLE statements for users, books, book_authors, book_genres, user_library, connections.

---

## Configuration Management

**Principle:** Never hardcode configuration values

### Application Properties Structure

Use YAML for configuration: `src/main/resources/application.yml`

**Profiles:**

- `application.yml` - default/shared config
- `application-dev.yml` - local development
- `application-prod.yml` - production (GCP Cloud Run)

**Environment Variables for Secrets:**

```yaml
spring:
  datasource:
    url: ${DATABASE_URL}  # Environment variable
    username: ${DATABASE_USER}
    password: ${DATABASE_PASSWORD}
  security:
    oauth2:
      client:
        registration:
          google:
            client-id: ${GOOGLE_CLIENT_ID}
            client-secret: ${GOOGLE_CLIENT_SECRET}

jwt:
  secret: ${JWT_SECRET_KEY}
  expiration: 604800000  # 7 days in milliseconds
```

**GCP Secret Manager Integration:**
Use Spring Cloud GCP to fetch secrets:

```xml
<dependency>
    <groupId>com.google.cloud</groupId>
    <artifactId>spring-cloud-gcp-starter-secretmanager</artifactId>
</dependency>
```

**Never Commit:**

- Database passwords
- OAuth client secrets
- JWT signing keys
- API keys

Use `.env` files for local development (added to `.gitignore`).

---

## Common Pitfalls & Best Practices

### 1. N+1 Query Problem in JPA

**Problem:**

```java
List<User> users = userRepository.findAll();
for (User user : users) {
    user.getBooks();  // Triggers separate SELECT for each user!
}
```

**Solution:**
Use `@EntityGraph` or JOIN FETCH:

```java
@Query("SELECT u FROM User u LEFT JOIN FETCH u.books WHERE u.id = :id")
User findByIdWithBooks(@Param("id") String id);
```

Or configure `@OneToMany(fetch = FetchType.LAZY)` and use JOIN FETCH in queries.

### 2. Entity vs DTO Confusion

**Wrong:**

```java
@GetMapping("/users/{id}")
public User getUser(@PathVariable String id) {
    return userRepository.findById(id);  // Exposes entity!
}
```

**Right:**

```java
@GetMapping("/users/{id}")
public UserResponse getUser(@PathVariable String id) {
    User entity = userRepository.findById(id);
    return mapToDTO(entity);  // Convert to DTO
}
```

**Why:** Entities may contain sensitive data, lazy-loaded collections (causes serialization issues), or internal
database structure.

### 3. Missing @Transactional

**Problem:** Multiple writes without transaction boundary:

```java
public void addBookAndNotifyFriends(Book book) {
    libraryRepository.save(book);  // Write 1
    notificationRepository.save(notification);  // Write 2
    // If Write 2 fails, Write 1 is still committed!
}
```

**Solution:**

```java
@Transactional
public void addBookAndNotifyFriends(Book book) {
    // Both commit together or rollback together
}
```

### 4. Hardcoded Configuration

**Bad:**

```java
String dbUrl = "jdbc:postgresql://localhost:5432/bookclub";
```

**Good:**

```java
@Value("${spring.datasource.url}")
private String dbUrl;
```

Or use `@ConfigurationProperties` for grouped settings.

### 5. Ignoring Database Indexes

**Problem:** Querying without indexes causes full table scans (slow at scale)

**Solution:** Add indexes for:

- Foreign keys (usually indexed automatically)
- WHERE clause columns (e.g., `isbn`, `status`)
- ORDER BY columns
- JOIN columns

**Example:**

```sql
CREATE INDEX idx_library_user_isbn ON user_library(user_id, isbn);
```

### 6. Leaking Exception Details

**Problem:**

```java
catch (SQLException e) {
    throw new RuntimeException(e.getMessage());  // Leaks database details!
}
```

**Solution:** Use custom exceptions and global exception handler:

```java
@ControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(BookNotFoundException.class)
    public ResponseEntity<ErrorResponse> handleNotFound(BookNotFoundException e) {
        return ResponseEntity.status(404)
            .body(new ErrorResponse("Book not found"));
    }
}
```

---

## Development Workflow

### Local Development Setup

1. **Start PostgreSQL:**
   ```bash
   docker run --name bookclub-postgres \
     -e POSTGRES_DB=bookclub \
     -e POSTGRES_USER=dev \
     -e POSTGRES_PASSWORD=dev \
     -p 5432:5432 \
     -d postgres:15-alpine
   ```

2. **Set Environment Variables:**
   Create `.env` file (gitignored):
   ```
   DATABASE_URL=jdbc:postgresql://localhost:5432/bookclub
   DATABASE_USER=dev
   DATABASE_PASSWORD=dev
   GOOGLE_CLIENT_ID=your-client-id
   GOOGLE_CLIENT_SECRET=your-client-secret
   JWT_SECRET_KEY=your-secret-key
   ```

3. **Run Application:**
   ```bash
   ./mvnw spring-boot:run
   ```

4. **Run Tests:**
   ```bash
   ./mvnw test                    # All tests
   ./mvnw test -Dtest=LibraryServiceTest  # Specific test
   ```

5. **Generate Code from OpenAPI:**
   ```bash
   ./mvnw generate-sources
   ```

### Iterative Development Process

1. **Update OpenAPI spec** (`src/main/resources/openapi.yaml`)
2. **Run `mvn generate-sources`** to regenerate interfaces/models
3. **Implement controller** methods (implement generated interface)
4. **Implement service** business logic
5. **Write unit tests** for service
6. **Write integration test** with Testcontainers
7. **Run tests:** `mvn test`
8. **Test locally** via Postman/curl

### Database Schema Changes

1. **Create migration file:** `V{n}__{description}.sql`
2. **Write SQL** (CREATE, ALTER, etc.)
3. **Restart app** (Flyway runs migration automatically)
4. **Update JPA entities** to match schema
5. **Commit migration file** to version control

---

## Learning Resources

### Phase 1: Java Fundamentals (if needed)

- **Records** (Java 14+): Immutable data classes
- **Streams API**: Functional programming (`map`, `filter`, `collect`)
- **Optional**: Null safety
- **CompletableFuture**: Async patterns

### Phase 2: Spring Core

- **Dependency Injection**: `@Autowired`, constructor injection
- **Component Scanning**: `@Component`, `@Service`, `@Repository`
- **Configuration**: `@Value`, `@ConfigurationProperties`
- **Profiles**: `@Profile`, `application-{profile}.yml`

### Phase 3: Spring Data JPA

- **Entity Mapping**: `@Entity`, `@Id`, `@Column`, `@Table`
- **Relationships**: `@ManyToOne`, `@OneToMany`, `@ManyToMany`, `@JoinColumn`
- **Repository Methods**: Query derivation (e.g., `findByIsbn`)
- **JPQL**: Java Persistence Query Language
- **Lazy vs Eager Loading**: `FetchType.LAZY`, `@EntityGraph`

### Phase 4: Spring Security

- **SecurityFilterChain**: Define protected/public endpoints
- **OAuth2 Client**: Google OAuth integration
- **JWT**: Token structure, signing, validation
- **Authentication vs Authorization**: Who you are vs what you can do

### Phase 5: Testing

- **JUnit 5**: `@Test`, `@BeforeEach`, lifecycle hooks
- **Mockito**: `@Mock`, `@InjectMocks`, `when()...thenReturn()`
- **AssertJ**: Fluent assertions (`assertThat(x).isEqualTo(y)`)
- **Spring Boot Test**: `@SpringBootTest`, `@WebMvcTest`, `@DataJpaTest`
- **Testcontainers**: Docker-based integration testing

### Recommended Documentation

- **Spring Boot Reference:** https://docs.spring.io/spring-boot/docs/current/reference/html/
- **Spring Data JPA:** https://docs.spring.io/spring-data/jpa/docs/current/reference/html/
- **Spring Security OAuth2:** https://docs.spring.io/spring-security/reference/servlet/oauth2/index.html
- **OpenAPI Generator:** https://openapi-generator.tech/docs/generators/spring/

---

## Next Steps: Getting Started

### Step 1: Schema Design (Do This First!)

Before writing any code, finalize the PostgreSQL schema:

**Tasks:**

1. Draw ERD (Entity-Relationship Diagram) with tables, columns, relationships
2. Decide: Normalized books table vs denormalized library table
3. Define all indexes
4. Write `V1__initial_schema.sql` migration file

**Tool Suggestion:** Use https://dbdiagram.io to visualize schema

**Questions to Answer:**

- How to handle `authors` array? (TEXT[] or separate table?)
- How to handle `genres` array? (TEXT[] or separate table?)
- What indexes are needed for query performance?

### Step 2: Spring Boot Project Initialization

Use **Spring Initializr** (https://start.spring.io):

**Settings:**

- **Build Tool:** Maven
- **Language:** Java
- **Java Version:** 17 or 21 (LTS)
- **Spring Boot Version:** 3.2.x (latest stable)
- **Packaging:** Jar
- **Dependencies:**
    - Spring Web
    - Spring Data JPA
    - PostgreSQL Driver
    - Spring Security
    - OAuth2 Client
    - Validation
    - Flyway Migration
    - Lombok (optional, reduces boilerplate)

Download and extract to `backend/` directory.

### Step 3: OpenAPI Specification

Start with ONE endpoint to learn the pattern:

**File:** `backend/src/main/resources/openapi.yaml`

**First endpoint to define:**

```yaml
openapi: 3.0.3
info:
  title: Bookclub API
  version: 1.0.0

paths:
  /api/users/{userId}/library:
    post:
      summary: Add book to user's library
      operationId: addBookToLibrary
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AddBookRequest'
      responses:
        '201':
          description: Book added successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BookResponse'
        '400':
          description: Invalid request
        '401':
          description: Unauthorized

components:
  schemas:
    AddBookRequest:
      type: object
      required:
        - isbn
        - title
        - authors
      properties:
        isbn:
          type: string
          minLength: 10
          maxLength: 13
        title:
          type: string
        # ... more properties

    BookResponse:
      type: object
      properties:
        isbn:
          type: string
        title:
          type: string
        # ... more properties
```

Expand to other endpoints after validating this pattern works.

### Step 4: Configure OpenAPI Codegen

Add to `pom.xml`:

```xml
<plugin>
    <groupId>org.openapitools</groupId>
    <artifactId>openapi-generator-maven-plugin</artifactId>
    <version>7.2.0</version>
    <executions>
        <execution>
            <goals>
                <goal>generate</goal>
            </goals>
            <configuration>
                <inputSpec>${project.basedir}/src/main/resources/openapi.yaml</inputSpec>
                <generatorName>spring</generatorName>
                <apiPackage>com.bookclub.api</apiPackage>
                <modelPackage>com.bookclub.api.model</modelPackage>
                <configOptions>
                    <interfaceOnly>true</interfaceOnly>
                    <useSpringBoot3>true</useSpringBoot3>
                </configOptions>
            </configuration>
        </execution>
    </executions>
</plugin>
```

Run: `mvn generate-sources`

### Step 5: Implement First Endpoint

**Create these files:**

1. **Entity:** `User.java`, `Book.java`, `UserLibrary.java` (JPA entities)
2. **Repository:** `UserRepository.java`, `BookRepository.java`, `UserLibraryRepository.java`
3. **Service:** `LibraryService.java` (business logic)
4. **Controller:** `LibraryController.java` (implements generated API interface)
5. **Test:** `LibraryServiceTest.java`, `LibraryIntegrationTest.java`

**This is your template for all other endpoints.**

### Step 6: Authentication Setup

1. Configure Google OAuth2 client in `application.yml`
2. Create `SecurityConfig.java` with `SecurityFilterChain`
3. Implement JWT token generation utility
4. Create JWT authentication filter
5. Test OAuth flow with Postman before integrating Flutter

### Step 7: Complete All Endpoints

Repeat Step 5 pattern for:

- List library books (GET)
- Add connection (POST)
- List connections (GET)
- Search by ISBN (GET)
- Get user profile (GET)

### Step 8: Docker & Deployment

1. Create `Dockerfile` (multi-stage build)
2. Test locally with Docker Compose (app + PostgreSQL)
3. Deploy to GCP Cloud Run
4. Set up Cloud SQL for PostgreSQL
5. Configure Secret Manager for credentials

---

## Success Criteria

You'll know the backend is ready when:

✅ All OpenAPI endpoints are implemented
✅ Integration tests pass with Testcontainers
✅ OAuth authentication flow works
✅ Database migrations run successfully
✅ Docker image builds and runs locally
✅ Deployed to GCP Cloud Run
✅ Flutter app can authenticate and fetch data

---

## Future Enhancements (Out of Scope for MVP)

- **Rate Limiting:** Protect API from abuse
- **Caching:** Redis for frequently accessed data
- **Pagination:** Limit response sizes for large collections
- **Full-Text Search:** PostgreSQL's `tsvector` for book search
- **File Upload:** Store book cover images (Cloud Storage)
- **WebSockets:** Real-time notifications
- **GraphQL:** Alternative to REST for flexible queries
- **API Versioning:** `/api/v1/` vs `/api/v2/`
- **Monitoring:** Application Performance Monitoring (APM)
- **CI/CD:** Automated testing and deployment pipeline

---

*This document is a living guide. Update it as architectural decisions evolve.*