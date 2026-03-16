package com.bookclub.config;

import com.google.api.client.util.Value;
import com.google.auth.oauth2.TokenVerifier;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.stereotype.Component;

@Configuration
public class GoogleAuthConfig {

    @Value("${google.client-id}")
    private String clientId;

    @Bean
    public TokenVerifier buildTokenVerifier() {
        return TokenVerifier.newBuilder().setAudience(clientId).build();
    }




}
