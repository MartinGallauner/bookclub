package com.bookclub.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("Bookclub API")
                        .version("1.0.0")
                        .description("REST API for the Bookclub application")
                        .contact(new Contact()
                                .name("Martin Gallauner")
                                .email("martin.gallauner@pmme")));
    }
}