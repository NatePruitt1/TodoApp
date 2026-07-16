# TodoApp Backend API Endpoints

This document lists the endpoint set discussed for a simple frontend-friendly API.

## Base
- Base path: `/api/v0`
- Authenticated routes require header: `Authorization: Bearer <jwt>`

## Auth Endpoints
These are typically public except where noted.

| Method | Path             | Auth                                | Purpose                          |
| ------ | ---------------- | ----------------------------------- | -------------------------------- |
| POST   | `/auth/register` | No                                  | Create account                   |
| POST   | `/auth/login`    | No                                  | Sign in and return auth token    |
| POST   | `/auth/refresh`  | No (uses refresh cookie/token flow) | Issue new auth token             |
| POST   | `/auth/logout`   | Yes                                 | Invalidate session/refresh token |
| GET    | `/auth/me`       | Yes                                 | Return current user profile      |

## Project Endpoints
| Method | Path                   | Auth | Purpose                                                    |
| ------ | ---------------------- | ---- | ---------------------------------------------------------- |
| GET    | `/projects`            | Yes  | List user projects                                         |
| POST   | `/projects`            | Yes  | Create project                                             |
| GET    | `/projects/:projectId` | Yes  | Get project board aggregate (project + categories + cards) |
| PATCH  | `/projects/:projectId` | Yes  | Update project name/description                            |
| DELETE | `/projects/:projectId` | Yes  | Delete project                                             |

### Suggested payloads
- `POST /projects`
```json
{
  "name": "My Project",
  "description": "Optional description"
}
```

- `PATCH /projects/:projectId`
```json
{
  "name": "Renamed Project",
  "description": "Updated description"
}
```

## Category Endpoints
| Method | Path                               | Auth | Purpose                                   |
| ------ | ---------------------------------- | ---- | ----------------------------------------- |
| POST   | `/projects/:projectId/categories`  | Yes  | Create category in project                |
| PATCH  | `/categories/:categoryId`          | Yes  | Update category fields (for example name) |
| DELETE | `/categories/:categoryId`          | Yes  | Delete category                           |
| PATCH  | `/categories/:categoryId/position` | Yes  | Reorder category                          |

### Suggested payloads
- `POST /projects/:projectId/categories`
```json
{
  "name": "To Do"
}
```

- `PATCH /categories/:categoryId`
```json
{
  "name": "In Progress"
}
```

- `PATCH /categories/:categoryId/position`
```json
{
  "index": 2
}
```

## Card Endpoints
| Method | Path                            | Auth | Purpose                                    |
| ------ | ------------------------------- | ---- | ------------------------------------------ |
| POST   | `/categories/:categoryId/cards` | Yes  | Create card in category                    |
| PATCH  | `/cards/:cardId`                | Yes  | Update card fields (title/content)         |
| DELETE | `/cards/:cardId`                | Yes  | Delete card                                |
| PATCH  | `/cards/:cardId/move`           | Yes  | Move card to another category and/or index |
| PATCH  | `/cards/:cardId/finish`         | Yes  | Toggle finished status                     |

### Suggested payloads
- `POST /categories/:categoryId/cards`
```json
{
  "title": "Implement API",
  "content": "Create handlers and service methods"
}
```

- `PATCH /cards/:cardId`
```json
{
  "title": "Implement API v2",
  "content": "Update endpoints and tests"
}
```

- `PATCH /cards/:cardId/move`
```json
{
  "targetCategoryId": "11111111-1111-1111-1111-111111111111",
  "index": 1
}
```

- `PATCH /cards/:cardId/finish`
```json
{
  "finished": true
}
```

## Note on Move Endpoints vs Pure CRUD
Both are valid:

1. Command style (clear intent):
- `PATCH /categories/:categoryId/position`
- `PATCH /cards/:cardId/move`

2. Pure CRUD-style PATCH (single update endpoint):
- `PATCH /categories/:categoryId` with `{ "index": ... }`
- `PATCH /cards/:cardId` with `{ "categoryId": ..., "index": ... }`

The command style is often easier to validate and maintain for board apps.
