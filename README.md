# Notification Service

A real-time notification service built with Go that handles both in-app and email notifications.

## Features
- Real-time notifications via WebSocket
- Email notification support
- Message queuing with RabbitMQ
- MongoDB for persistence
- REST API for notification history

## Tech Stack
- Go 1.22+
- MongoDB
- RabbitMQ
- WebSocket
- Docker

## Setup
1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `docker-compose up` to start required services
4. Run `go run cmd/server/main.go` to start the application

## Architecture
[Add your flowchart or architecture diagram here]