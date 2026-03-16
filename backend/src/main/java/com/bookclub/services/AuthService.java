package com.bookclub.services;

import com.bookclub.api.model.AuthResponse;
import org.springframework.stereotype.Service;

@Service
public class AuthService {


    public AuthResponse authenticate(String idToken) {
        return new AuthResponse();

    }
}
