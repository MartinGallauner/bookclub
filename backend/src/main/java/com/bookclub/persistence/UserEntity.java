package com.bookclub.persistence;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(name = "users")
public class UserEntity {

    @Id
    public String id;
    public String displayName;
    public String photoUrl;
    public String email;
}
