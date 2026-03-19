# Stage 1 - Build
FROM eclipse-temurin:25-jdk AS builder
WORKDIR /app
COPY openapi/api.yaml /openapi/api.yaml
COPY backend/gradlew backend/gradlew.bat ./
COPY backend/gradle ./gradle
COPY backend/build.gradle backend/settings.gradle ./
COPY backend/src src
RUN chmod +x ./gradlew
RUN ./gradlew dependencies --no-daemon
RUN ./gradlew build  -x test --no-daemon

# Stage 2 - Runtime
FROM eclipse-temurin:25-jre
WORKDIR /app
COPY --from=builder /app/build/libs/*.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]