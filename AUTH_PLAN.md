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

## Step 1: Update OpenAPI Spec (`openapi/api.yaml`)

- Change `GoogleLoginRequest` schema:
  - Remove: `code` (authorization code) and `redirectUri`
  - Add: `idToken` (Google ID token from Flutter's google_sign_in)
- Update the endpoint description to reflect the new flow
- `AuthResponse` and `UserProfile` schemas are correct as-is

---

## Step 2: Backend Implementation

### After updating spec, run:
```bash
cd backend && ./gradlew openApiGenerate
```

### A) Google ID Token Verification
- Add dependency: `google-auth-library-oauth2-http`
- Use `GoogleIdTokenVerifier` to validate the incoming `idToken`
- Verifies signature using Google's public keys
- Verifies `GOOGLE_CLIENT_ID` is in the token's audience claim
- Extract user info from token payload: Google UID, email, displayName, photoUrl

### B) JWT Issuance
- Add dependency: `com.auth0:java-jwt`
- After verifying Google token:
  1. Upsert user in PostgreSQL `users` table (create on first login, update on return)
  2. Generate app JWT signed with `JWT_SECRET_KEY`
  3. JWT claims: user ID, email, expiration (7 days)
- Return `AuthResponse` with the JWT

### C) Spring Security Configuration
- Create `SecurityConfig.java` with `SecurityFilterChain`:
  - `permitAll`: `/api/auth/google`, `/api/health`
  - All other endpoints require valid JWT
- Create JWT filter (`OncePerRequestFilter`):
  - Reads `Authorization: Bearer {token}` header
  - Validates JWT signature and expiration
  - Populates Spring Security context on success

---

## Step 3: Flutter Implementation

1. Add `google_sign_in` to `pubspec.yaml`
2. In auth service:
   - Call `googleSignIn.signIn()` to trigger Google UI
   - Get `GoogleSignInAuthentication` from the signed-in account
   - Extract `auth.idToken`
3. POST `{idToken}` to `POST /api/auth/google`
4. Store the returned JWT (e.g., in `SharedPreferences` or `FlutterSecureStorage`)
5. Attach JWT to all subsequent requests: `Authorization: Bearer {token}`

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

---

## Environment Variables Needed

```
GOOGLE_CLIENT_ID=your-google-client-id
JWT_SECRET_KEY=your-signing-secret
```

`GOOGLE_CLIENT_ID` is used both by Flutter (to initiate sign-in) and by the backend (to verify token audience).
