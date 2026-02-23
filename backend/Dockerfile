# Stage 1 - Build
FROM eclipse-temurin:25-jdk AS builder
WORKDIR /app
COPY gradlew gradlew.bat gradle ./
RUN chmod +x ./gradlew
COPY build.gradle settings.gradle ./
RUN ./gradlew dependencies --no-daemon
COPY src src
RUN ./gradlew build  -x test --no-daemon

# Stage 2 - Runtime
FROM eclipse-temurin:25-jre
WORKDIR /app
COPY --from=builder /app/build/libs/boookclub-backened-0.0.1-SNAPSHOT.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]