# 🦆 DuckParty Back-end

> RESTful API backend for DuckParty - where ducks come to party and show off their style.

A robust, scalable backend service powering the DuckParty platform. Handles duck creation, customization, user authentication, leaderboards, and interactions with a modern tech stack.

## 🔗 Front-end

This backend powers the DuckParty frontend application. For frontend setup and documentation, visit:

**[duckparty-frontend](https://github.com/omidnikrah/duckparty-frontend)**

## ✨ Features

- **User Authentication** - JWT-based auth with email OTP verification
- **Duck Management** - Create, customize, and manage duck collections
- **Leaderboard System** - Ranking based on reactions
- **Reaction System** - Like/dislike ducks with rate limiting
- **Image Storage** - AWS S3 integration for duck image hosting
- **Email Service** - AWS SES for OTP delivery
- **API Documentation** - Swagger/OpenAPI documentation
- **Scheduled Tasks** - Cron jobs for automated operations

## 🛠️ Tech Stack

- **[Go](https://go.dev/)** - High-performance backend language
- **[Gin](https://gin-gonic.com/)** - Fast HTTP web framework
- **[GORM](https://gorm.io/)** - ORM for database operations
- **[PostgreSQL](https://www.postgresql.org/)** - Relational database
- **[Redis](https://redis.io/)** - Caching and rate limiting
- **[AWS S3](https://aws.amazon.com/s3/)** - Object storage for images
- **[AWS SES](https://aws.amazon.com/ses/)** - Email delivery service
- **[JWT](https://jwt.io/)** - Token-based authentication
- **[Swagger](https://swagger.io/)** - API documentation

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.24.4 or higher
- [PostgreSQL](https://www.postgresql.org/download/) 16 or higher
- [Redis](https://redis.io/download) 7 or higher
- AWS account with S3 and SES configured
- Docker and Docker Compose (optional, for containerized setup)

### Installation

```bash
# Clone the repository
git clone https://github.com/omidnikrah/duckparty-backend.git
cd duckparty-backend

# Install dependencies
go mod download

# Set up environment variables (see Environment Variables section)
cp .env.example .env
# Edit .env with your configuration

# Start the server
air
```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
APP_PORT=4030
API_PREFIX=/api

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=duckparty

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# AWS
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
S3_BUCKET=your_s3_bucket
S3_REGION=your_s3_region
S3_BASE_URL=your_s3_base_url

# JWT
JWT_SECRET=your_jwt_secret_key

# Email
AUTH_SENDER_EMAIL=your_verified_ses_email
```

### Docker Setup

```bash
# Start all services (PostgreSQL, Redis, and App)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📚 API Documentation

Once the server is running, access the interactive API documentation at:

**http://localhost:4030/swagger/index.html**

## 📝 Project Structure

```
duckparty-backend/
├── cmd/
│   └── server/          # Server setup and initialization
├── internal/
│   ├── client/          # External service clients (Redis, SES, Cron)
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── dto/             # Data transfer objects
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, rate limiting, validation)
│   ├── model/           # Database models
│   ├── routes/          # API route definitions
│   ├── service/         # Business logic layer
│   ├── storage/         # Storage abstractions (S3)
│   ├── templates/       # Email templates
│   ├── types/           # Type definitions
│   └── utils/           # Utility functions
├── docs/                # Swagger documentation
├── docker-compose.yml   # Docker services configuration
├── Dockerfile           # Container build configuration
└── main.go              # Application entry point
```

## 🚢 Deployment

This project is deployed using [Coolify](https://coolify.io/).

---

Made with ❤️‍🔥 for the duck community
