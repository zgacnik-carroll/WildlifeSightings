# Wildlife Sightings

Wildlife Sightings is a server-rendered Go web application for tracking wildlife sightings.
It provides user registration/login, protected routes, and CRUD operations for sightings, with search and profile views.

---

## Features

- User authentication with session cookies
- Register, login, and logout flows
- Protected routes for authenticated users
- Create, edit, and delete wildlife sightings
- Search sightings by animal or location
- Profile page showing a user's sightings
- Basic aggregate stats on the dashboard

---

## Tech Stack

- Go 1.26
- [Gin](https://github.com/gin-gonic/gin) for routing and HTTP handling
- [GORM](https://gorm.io/) as ORM
- [glebarez/sqlite](https://github.com/glebarez/sqlite) SQLite driver
- [gorilla/sessions](https://github.com/gorilla/sessions) for cookie-based sessions
- HTML templates + static CSS/JS assets

---

## Project Structure

```text
Go-Web/
  db/
    db.go              # Database initialization, models, migrations
  handlers/
    auth.go            # Register/login/logout handlers
    handlers.go        # Sighting and profile handlers
  middleware/
    auth.go            # Auth-required middleware and session store
  static/
    css/style.css      # Styles
    js/main.js         # Client-side JS
  templates/           # HTML templates used by Gin
  main.go              # Application entrypoint and route setup
  go.mod               # Module + Go version
  go.sum               # Dependency lock file
  wildlife.db          # Local SQLite database file
```

---

## Requirements

- Go 1.26+
- Git
- Internet access on first run (to download Go modules)
- Port `8080` available locally

---

## Getting Started

1. Clone the repository:

```bash
git clone <your-repo-url>
cd WildlifeSightings
```

2. Install dependencies:

```bash
go mod download
```

3. Run the app:

```bash
go run main.go
```

4. Open in browser:

```text
http://localhost:8080
```

On startup, the app will initialize SQLite and auto-migrate the `users` and `sightings` tables.

---

## Routes

### Public

- `GET /register` - registration form
- `POST /register` - create account
- `GET /login` - login form
- `POST /login` - authenticate user
- `POST /logout` - logout user

### Protected (requires login)

- `GET /` - dashboard + aggregate stats
- `GET /sightings/new` - new sighting form
- `POST /sightings` - create sighting
- `GET /sightings/search` - search sightings
- `GET /profile` - current user profile and sightings
- `GET /sightings/:id/edit` - edit form
- `POST /sightings/:id/edit` - update sighting
- `POST /sightings/:id/delete` - delete sighting

---

## Database

The app uses a local SQLite file:

- `wildlife.db` in the project root

Schema is managed via GORM AutoMigrate in `db.Init()`, and includes:

- `User` (`id`, `username`, `password`, `created_at`)
- `Sighting` (`id`, `animal`, `location`, `notes`, `user_id`, `created_at`)

---

## Configuration Notes

- The session secret is currently hardcoded in `middleware/auth.go`:
  - `wildlife-secret-key`
- For production use, move this to an environment variable and set secure cookie options.

---

## Security and Production Considerations

- Replace the hardcoded session key with an environment variable.
- Configure secure cookies (`HttpOnly`, `Secure`, `SameSite`) appropriately.
- Add CSRF protection for state-changing routes.
- Add input validation and stricter error handling.
- Put the app behind TLS in production.

---

## Closing Remarks

WildlifeSightings was created to enhance my experience in the Go programming language,
while bringing in unique project perspectives from my own personal life.