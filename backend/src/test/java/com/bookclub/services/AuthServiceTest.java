package com.bookclub.services;

import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.UserProfile;
import com.bookclub.persistence.UserEntity;
import com.bookclub.persistence.UserRepository;
import com.google.api.client.json.webtoken.JsonWebSignature;
import com.google.api.client.json.webtoken.JsonWebToken;
import com.google.auth.oauth2.TokenVerifier;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import org.testcontainers.postgresql.PostgreSQLContainer;


import java.util.Base64;
import java.util.Optional;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@SpringBootTest
@Testcontainers
public class AuthServiceTest {

    @Container
    static PostgreSQLContainer postgres = new PostgreSQLContainer("postgres:15-alpine");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        registry.add("JWT_SECRET", () -> "testsecret");
        registry.add("JWT_EXPIRATION_MS", () -> 604800000);
    }

    @Autowired
    private AuthService authService;

    @Autowired
    private UserRepository userRepository;

    @MockitoBean
    private TokenVerifier tokenVerifier;

    @MockitoBean
    private JsonWebSignature jws;

    @DisplayName("New user logs in, is persisted in the database and a jwt is created.")
    @Test
    public void TestNewUserLogin() throws TokenVerifier.VerificationException {
        //given
        String idToken = "idToken";

        JsonWebToken.Payload payload = new JsonWebToken.Payload();
        payload.setSubject("google-uid-123");
        payload.set("email", "test@example.com");
        payload.set("name", "Test User");
        when(jws.getPayload()).thenReturn(payload);
        when(tokenVerifier.verify(anyString())).thenReturn(jws);

        //when
        AuthResponse authResponse = authService.authenticate(idToken);


        //then
        UserProfile user = authResponse.getUser();
        Assertions.assertEquals("Test User", user.getDisplayName().get());

        verify(tokenVerifier, times(1)).verify(anyString());

        //verify token
        String[] token = authResponse.getToken().split("\\.");
        Assertions.assertEquals(3, token.length);
        String tokenPayload = new String(Base64.getUrlDecoder().decode(token[1]));
        Assertions.assertTrue(tokenPayload.contains("google-uid-123"));

        Optional<UserEntity> persistedUser = userRepository.findById(user.getId());
        Assertions.assertTrue(persistedUser.isPresent());



    }


}
