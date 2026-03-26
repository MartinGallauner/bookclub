package com.bookclub.persistence;

import com.google.api.client.json.webtoken.JsonWebToken;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

@Builder
@Entity
@Table(name = "users")
@NoArgsConstructor
@AllArgsConstructor
public class UserEntity {

    @Id
    public String id;
    public String displayName;
    public String photoUrl;
    public String email;

    public static UserEntity fromJws(JsonWebToken jws) {
        JsonWebToken.Payload payload = jws.getPayload();
        String photoUrl = payload.get("picture") != null ? payload.get("picture").toString() : null;
        String email = payload.get("email") != null ? payload.get("email").toString() : null;
        String displayName = payload.get("name") != null ? payload.get("name").toString() : null;
        return UserEntity.builder().id(payload.getSubject()).displayName(displayName).email(email).photoUrl(photoUrl).build();
    }
}