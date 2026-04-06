# Auth Implementation Plan — Option B (Client-Side Google Sign-In + Backend Verification)

## Flow Summary

```
Flutter → google_sign_in SDK → gets Google ID token
       → POST /api/auth/google {idToken: "..."}
       → Backend verifies ID token with Google's public keys
       → Backend upserts user in PostgreSQL
       → Backend issues app JWT
       → Flutter stores JWT, sends as Bearer on all subsequent requests
```

---

## Step 1: Update OpenAPI Spec (`openapi/api.yaml`) ✅ DONE

- `GoogleLoginRequest` schema uses `idToken` field (Google ID token)
- `AuthResponse` schema: `token`, `tokenType`, `expiresIn`, `user` (UserProfile)
- `UserProfile` schema: `id`, `email`, `displayName`, `photoUrl`
- Endpoint description still mentions "authorization code" — update to reflect ID token flow

---

## Step 2: Backend Implementation ✅ DONE

Java backend at `backend/` is fully implemented:

### A) Google ID Token Verification
- Dependency: `google-auth-library-oauth2-http`
- `GoogleIdTokenVerifier` validates the incoming `idToken` against Google's public keys
- Verifies `GOOGLE_CLIENT_ID` is in the token's audience claim
- Extracts user info from token payload: Google UID, email, displayName, photoUrl

### B) JWT Issuance
- Dependency: `com.auth0:java-jwt`
- On valid Google token: upserts user in PostgreSQL, issues signed app JWT (7-day expiry)
- Returns `AuthResponse` with JWT and `UserProfile`

### C) Spring Security Configuration
- `SecurityConfig.java` with `SecurityFilterChain`:
  - `permitAll`: `/api/auth/google`, `/api/health`
  - All other endpoints require valid JWT
- JWT filter (`OncePerRequestFilter`): reads `Authorization: Bearer` header, validates, populates security context

### After any spec change, regenerate:
```bash
cd backend && ./gradlew openApiGenerate
```

---

## Step 3: Flutter Implementation (`frontend/`)

### 3a: Set up OpenAPI codegen

The frontend has no generated API client yet (`lib/generated/api/` is empty).
Use the `dart-dio` generator — produces a Dio-based typed client with model classes.

**`mise.toml`:**
- Add `'npm:@openapitools/openapi-generator-cli' = "latest"` to `[tools]`
- Add a `[tasks.generate]` task:
  ```
  openapi-generator-cli generate \
    -i ../openapi/api.yaml \
    -g dart-dio \
    -o lib/generated/api \
    --additional-properties=pubName=bookclub_api
  ```
- Make `test` depend on `generate` (mirrors `compileJava.dependsOn openApiGenerate` in backend)
- Add `lib/generated/api/` to `.gitignore`

Run: `mise run generate`

Generated output includes:
- `AuthApi` class with `loginWithGoogle(GoogleLoginRequest)` method
- `GoogleLoginRequest`, `AuthResponse`, `UserProfile` model classes (use these everywhere — don't create manual equivalents)

### 3b: Update `pubspec.yaml`

Add dependencies:
- `google_sign_in` — replaces Firebase Auth for the Google popup
- `flutter_secure_storage` — stores JWT in OS keychain/keystore
- `dio` — required by dart-dio generated client
- `json_annotation` — required by generated models

Add dev dependencies:
- `build_runner`
- `json_serializable`

Remove:
- `firebase_auth`
- `cloud_firestore`

Remove from dev dependencies:
- `fake_cloud_firestore`

Keep `firebase_core` only if other Firebase services still need it (remove when fully migrated).

### 3c: Create `lib/services/api_client.dart`

Configure a shared `Dio` instance — this is the HTTP client factory for all API calls:
- Base URL: `http://localhost:8080`
- Add an `InterceptorsWrapper` that reads the JWT from `flutter_secure_storage` and injects `Authorization: Bearer <token>` on every request
- Handle 401 responses globally (force logout / clear stored token)

Pass this `Dio` instance to all generated API classes (e.g., `AuthApi(apiClient)`).

### 3d: Rewrite `lib/services/auth_service.dart`

Remove all Firebase imports. New implementation:
- Constructor takes `AuthApi` (generated) and `FlutterSecureStorage`
- `signInWithGoogle()`:
  1. `GoogleSignIn().signIn()` → triggers Google popup
  2. `account.authentication` → `GoogleSignInAuthentication`
  3. Extract `auth.idToken`
  4. Call `authApi.loginWithGoogle(GoogleLoginRequest(idToken: idToken))`
  5. Store `response.data.token` in secure storage
  6. Return `response.data.user` (generated `UserProfile`)
- `signOut()` — deletes token from secure storage
- Remove: `authStateChanges()`, `fetchProfile()` — backend handles profile creation now

### 3e: Update `AppState` in `main.dart`

- Replace `User?` (Firebase type) with the generated `UserProfile?`
- Remove `StreamSubscription` / `authStateChanges()` listener — no Firebase stream
- Add explicit `login(UserProfile user)` and `logout()` methods that call `notifyListeners()`
- On app startup: check secure storage for an existing token; if absent, stay on LoginPage (simple approach — no token validation on startup for now)

### 3f: Update provider tree in `main.dart`

- `AuthService` no longer takes `FirebaseAuth` or `FirebaseFirestore` — inject `AuthApi` and `FlutterSecureStorage`
- Remove `Firebase.initializeApp()` from `main()` once all Firebase packages are removed
- `LibraryService` and `ContactService` still use Firestore — leave them unchanged; they will be migrated in the next phase

### 3g: `lib/pages/login_page.dart`

No changes expected — already calls `signInWithGoogle` on the service.

### Execution order

```
3a: mise run generate        (produces Dart client from OpenAPI spec)
3b: flutter pub get          (installs new dependencies)
3c: create api_client.dart   (configured Dio instance with JWT interceptor)
3d: rewrite auth_service.dart
3e: update AppState
3f: update provider tree in main.dart
3g: verify login_page.dart
```

---

## Responsibility Split

| Concern                        | Owner   |
|-------------------------------|---------|
| Triggering Google Sign-In UI  | Flutter |
| Getting the Google ID token   | Flutter |
| Verifying the ID token        | Backend |
| Managing users in PostgreSQL  | Backend |
| Issuing the app JWT           | Backend |
| Storing & sending the JWT     | Flutter |
| Injecting JWT on requests     | Flutter (Dio interceptor in api_client.dart) |

---

## Environment Variables Needed

```
GOOGLE_CLIENT_ID=your-google-client-id
JWT_SECRET_KEY=your-signing-secret
```

`GOOGLE_CLIENT_ID` is used both by Flutter (to initiate sign-in) and by the backend (to verify token audience).

---

## Step 4: GET /api/users/me — Session Restoration

### 4a: Update OpenAPI Spec (`openapi/api.yaml`)

Add `GET /api/users/me` path:
- `operationId`: `getCurrentUser`
- Tags: `Users`
- No `security: []` override — JWT required
- `200` response: `UserProfile`
- `401` response: `ErrorResponse`

Regenerate after: `./gradlew openApiGenerate`

### 4b: Backend Implementation

- Controller method implementing the generated `getCurrentUser` interface
- Read authenticated user ID from Spring Security context (`SecurityContextHolder`)
- Call user service/repository to fetch `UserProfile` by ID
- Return `UserProfile` DTO

### 4c: Flutter — Startup Session Check

In `AppState` (or `main.dart`), on app startup:
1. Read JWT from `FlutterSecureStorage`
2. If token exists, call `usersApi.getCurrentUser()` (generated client)
3. On success: call `appState.login(response.data.user)` — restores session silently
4. On 401: token expired — delete it from storage, stay on `LoginPage`
5. If no token: stay on `LoginPage`

This replaces the `authStateChanges()` stream that Firebase provided.

---

## What's Left After This (Next Phase)

- Migrate `LibraryService` from Firestore to REST API (`/api/users/{userId}/library`)
- Migrate `ContactService` from Firestore to REST API (`/api/connections`)
- Both will use the same `api_client.dart` Dio instance (JWT injection comes for free)
