# Simple Todo App

A full-stack Kanban board application built with Go-Gin, Typescript-React, and PostgreSQL.

[https://nathanielpruitt.com]
Note: If there are issues with access, please contact me at [nathaniel.n.pruitt@gmail.com]

## About

This is a simple todo application built by me, Nate Pruitt, for the following purposes:
 - Practice full stack development, with a focus on deployable code
 - Deploy a working application to the internet
 - Have a "clickable" element on my resume to showcase my capabilities as a software engineer.

This project will continue to be updated and will eventually be fully featured.

## Features
 - User Authentication
    - Register and Login via JWTs.
    - Refresh using specialized refresh tokens.
 - Kanban board data structures with (partial) CRUD interfaces.
    - Projects -> Categories -> Cards.
    - There are extra, specialized command endpoints that allow one to move cards and categories.
    - Cards and categories currently have no get endpoints, instead edits are seen through the project aggregate.
 - Simple Frontend.
    - Currenlty only basic styles and functionalities, major focus of future updates.

## Tech Stack
**Frontend:** Typescript, React, Vite
**Backend:** Go, Gin, Postgres
**Infra:** Docker Compose, Cloudflare Tunnels, Nginx.

## Architecture
 - Frontend: Accessed through a cloudflare tunnel container.
    - Tunnel is configured to send requests to the frontend nginx container.
    - Nginx container exposes built React frontend.
    - Nginx container routes /api requests to the go backend container.
 - Backend: Accessed through the go backend container.
    - Routed requests are processed and fufilled by the go/gin backend.
    - Backend interacts with postgres database within the containers' virtual private network.

## Getting Started
### Prerequisites
 - Docker & Docket Compose.

All other tools are included automatically by the dockerfiles. 
If you would like to run the code without containers, you will need
 - Node.js / npm
 - Go
 - Go migrate

### Setup (Dev)
\`\`\`bash
git clone ...
cp apps/backend/.env.example apps/backend/.env.development
cp apps/backend/internal/db/.env.db.development.example apps/backend/internal/db/.env.db.development
cp apps/frontend/.env.example apps/frontend/.env # We don't use a dev env here becuase this is configuration, not secrets.
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
\`\`\`
