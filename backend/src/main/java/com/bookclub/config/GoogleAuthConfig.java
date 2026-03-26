package com.bookclub.config;

import com.google.auth.oauth2.TokenVerifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GoogleAuthConfig {

    @Value("${google.client-id}")
    private String clientId;

    @Bean
    public TokenVerifier buildTokenVerifier() {
        return TokenVerifier.newBuilder().setAudience(clientId).build();
    }




}
