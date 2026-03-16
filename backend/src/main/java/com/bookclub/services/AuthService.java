package com.bookclub.services;

import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.UserProfile;
import com.google.api.client.json.webtoken.JsonWebSignature;
import com.google.auth.oauth2.TokenVerifier;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private final TokenVerifier tokenVerifier;

    public AuthService(TokenVerifier tokenVerifier) {
        this.tokenVerifier = tokenVerifier;
    }

    /**
     * This method verifies the idToken, upserts the user and creates the JWT for further use.
     * @param idToken the token auth providers token passed from the frontend.
     * @return The complete AuthResponse including the jwt
     */
    public AuthResponse authenticate(String idToken) throws TokenVerifier.VerificationException {
        JsonWebSignature jws;
        try {
            jws = tokenVerifier.verify(idToken);
        } catch (TokenVerifier.VerificationException e) {
            throw new TokenVerifier.VerificationException(e.getMessage());
        }

        String name = jws.getPayload().get("name").toString();


        return new AuthResponse().user(new UserProfile().displayName(name)).token("fancyjwttoken");

    }
}
