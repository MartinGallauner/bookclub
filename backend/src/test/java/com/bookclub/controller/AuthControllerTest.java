package com.bookclub.controller;

import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.UserProfile;
import com.bookclub.config.OpenApiConfig;
import com.bookclub.config.SecurityConfig;
import com.bookclub.services.AuthService;
import com.google.auth.oauth2.TokenVerifier;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.webmvc.test.autoconfigure.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultHandlers;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

import java.net.URI;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

@WebMvcTest(AuthController.class)
@Import({SecurityConfig.class, GlobalExceptionHandler.class, OpenApiConfig.class})
public class AuthControllerTest {

    @Autowired
    private AuthController authController;

    @MockitoBean
    private AuthService authService;

    @Autowired
    private MockMvc mockMvc;

    @Test
    public void TestAuthSuccessful() throws Exception {

        UserProfile user = new UserProfile()
                .id("42")
                .displayName("Mock User")
                .photoUrl(URI.create("https://upload.wikimedia.org/wikipedia/en/7/73/Trollface.png"))
                .email("mockuser@mail.com");

        when(authService.authenticate(anyString())).thenReturn(new AuthResponse("token", AuthResponse.TokenTypeEnum.BEARER, 1234, user));

        mockMvc.perform(MockMvcRequestBuilders.post("/api/auth/google")
                        .header("Authorization", "Bearer token")
                        .header("Content-Type", "application/json")
                        .content("{\"idToken\":\"testToken\"}"))
                .andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(MockMvcResultMatchers.jsonPath("$.token").value("token"))
                .andExpect(MockMvcResultMatchers.jsonPath("$.tokenType").value("Bearer"))
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.displayName").value("Mock User"))
                .andExpect(MockMvcResultMatchers.jsonPath("$.user.email").value("mockuser@mail.com"));
    }

    @Test
    public void TestAuthFailed() throws Exception {
        when(authService.authenticate(anyString())).thenThrow(new TokenVerifier.VerificationException("no bueno"));

        mockMvc.perform(MockMvcRequestBuilders.post("/api/auth/google")
                        .header("Authorization", "Bearer token")
                        .header("Content-Type", "application/json")
                        .content("{\"idToken\":\"testToken\"}"))
                .andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isUnauthorized()).andExpect(MockMvcResultMatchers.jsonPath("$.error").value("no bueno"));
    }
}
