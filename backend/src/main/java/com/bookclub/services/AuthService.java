package com.bookclub.services;

import com.bookclub.api.model.AuthResponse;
import com.bookclub.api.model.UserProfile;
import com.bookclub.persistence.UserEntity;
import com.bookclub.persistence.UserRepository;
import com.google.api.client.json.webtoken.JsonWebSignature;
import com.google.api.client.json.webtoken.JsonWebToken;
import com.google.auth.oauth2.TokenVerifier;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private final TokenVerifier tokenVerifier;
    private final UserRepository userRepository;

    public AuthService(TokenVerifier tokenVerifier, UserRepository userRepository) {
        this.tokenVerifier = tokenVerifier;
        this.userRepository = userRepository;
    }

    /**
     * This method verifies the idToken, upserts the user and creates the JWT for further use.
     *
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

        UserEntity user = UserEntity.fromJws(jws);
        user = userRepository.save(user);

        return new AuthResponse().user(new UserProfile()
                        .id(user.id)
                        .email(user.email)
                        .displayName(user.displayName))
                        .token("fancyjwttoken");
    }
}