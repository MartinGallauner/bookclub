package com.bookclub.services;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.Date;

// Provides functionality to create JWT for the client.
@Service
public class TokenService {

    @Value("${JWT_SECRET}")
    private String secret;

    @Value("${JWT_EXPIRATION_MS}")
    private long expirationInMs;

    /**
     * Creates a signed JWT containing the user's identity claims, for use as a Bearer token in subsequent API requests.
     * @param userId the users unique identifier.
     * @param email the users email
     * @return a signed jwt token
     */
    public String createToken(String userId, String email) {
        return JWT.create()
                .withSubject(userId)
                .withClaim("email", email)
                .withExpiresAt(new Date(System.currentTimeMillis() + expirationInMs))
                .sign(Algorithm.HMAC256(secret));
    }
}
