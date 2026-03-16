package com.bookclub.services;

import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.UserProfile;
import com.google.api.client.json.webtoken.JsonWebSignature;
import com.google.api.client.json.webtoken.JsonWebToken;
import com.google.auth.oauth2.TokenVerifier;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;



import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
public class AuthServiceTest {

    @InjectMocks
    private AuthService authService;

    @Mock
    private TokenVerifier verifier;

    @Mock
    private JsonWebSignature jws;

    @DisplayName("New user logs in, is persisted in the database and a jwt is created.")
    @Test
    public void TestNewUserLogin() throws TokenVerifier.VerificationException {
        String idToken = "idToken";

        JsonWebToken.Payload payload = new JsonWebToken.Payload();
        payload.setSubject("google-uid-123");
        payload.set("email", "test@example.com");
        payload.set("name", "Test User");
        when(verifier.verify(anyString())).thenReturn(jws);
        when(jws.getPayload()).thenReturn(payload);

        AuthResponse authResponse = authService.authenticate(idToken);

        UserProfile user = authResponse.getUser();
        Assertions.assertEquals("Test User", user.getDisplayName().get());
        Assertions.assertEquals("fancyjwttoken", authResponse.getToken());

        verify(verifier, times(1)).verify(anyString());



    }


}
