package com.bookclub.controller;

import com.bookclub.api.AuthApi;
import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.GoogleLoginRequest;
import com.bookclub.services.AuthService;
import com.google.auth.oauth2.TokenVerifier;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class AuthController implements AuthApi {

    private AuthService authService;

    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    @Override
    public ResponseEntity<AuthResponse> loginWithGoogle(GoogleLoginRequest googleLoginRequest) {
        try {
            return ResponseEntity.ok(authService.authenticate(googleLoginRequest.getIdToken()));
        } catch (TokenVerifier.VerificationException e) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        }
    }
}
