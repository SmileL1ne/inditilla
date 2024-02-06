# INDITILLA

Test backend of user profile page 

## Endpoints

- **POST: /v1/user/signup** - sign up new user (returns registered user's id)
- **POST: /v1/user/login** - log in existing user (returns JWT access token)
- **GET: /v1/user/profile/:id** - get user profile info (returns user profile information)
- **PATCH: /v1/user/profile/:id** - update user info (returns updated user info)

## Usage

1. Navigate to the project directory
2. Install dependecies:
```bash
    go get -v ./...
```
3. Create ".env" file by example ".env.example" file and fill all required fields
4. Run the application with Makefile command:
```bash
    make run
```

> If make command is not installed - run manually:

```bash
    go run ./cmd/app
```

> [!WARNING]
> This project uses postgresql, specifically - 'pgx' package for database connection and management