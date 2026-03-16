package com.bookclub.services;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
public class AuthServiceTest {

    @InjectMocks
    private AuthService authService;

    @Test
    public void TestNewUserLogin() {
        String idToken = "idToken";

        var authResponse = authService.authenticate(idToken);

    }


}
