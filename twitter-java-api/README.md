# Twitter Java API (Legacy/Alternative)

Spring Boot backend kept in the monorepo as an alternative implementation.

Primary backend for active development is `twitter-go-api`.

## Stack

- Java 21
- Spring Boot 4
- Spring Security
- Spring Data JPA
- Flyway
- PostgreSQL
- Maven

## Local Setup

### Prerequisites

- JDK 21+
- Maven 3.8+
- PostgreSQL

### Configure

Edit:
- `src/main/resources/application.yml`

Set DB connection and auth/storage settings for your environment.

### Run

```bash
mvn spring-boot:run
```

or

```bash
mvn clean package
java -jar target/twitter-java-api-*.jar
```

API base URL: `http://localhost:8080/api/v1`

## Testing

```bash
mvn test
```

## Scope Note

This service remains in repository for reference/comparison.
New architecture and behavior changes are implemented first in `twitter-go-api`.
