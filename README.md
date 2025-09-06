# e-Invoice API

This is a **learning-in-progress project** where I am building an API service in **Go** using the [Gin](https://github.com/gin-gonic/gin) framework, with plans to integrate with Malaysia's **LHDN e-Invoice** system.  

The main purpose of this repo is to **practice and document my learning journey** in backend development with Go.  

While this is currently a personal learning project, my **long-term goal** is to shape it into a **production-ready application**.

## 🚀 Project Goals
- Learn Go by building a real-world style API with **Gin**.
- Experiment with folder structures and clean architecture.
- Integrate with LHDN e-Invoice.
- Practice database usage, migrations, and background jobs.
- Explore deployment options (OCI, later AWS).

## ⚠️ Note
This repository is **not production-ready yet** and I will not be accepting external contributions at this stage.  
I plan to evolve this project step by step towards a production-grade service.

## 📂 Current Structure (Work in Progress)
```
├── cmd
│   ├── api
│   └── migrate
├── config
├── internal
│   ├── database
│   │   ├── dberrors
│   │   ├── dbinfo
│   │   ├── factory
│   │   ├── migrations
│   │   └── models
│   ├── handlers
│   └── services
├── pkg
├── scripts
└── test
```

## 🛠️ Planned Tech Stack
- Go
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- PostgreSQL - ([golang-migrate](https://github.com/golang-migrate/migrate) + [bob](https://github.com/stephenafamo/bob))
- REST API
- Cache - Redis
- Queue - [asnyq](https://github.com/hibiken/asynq)
- Docker (planned)
- Cloud hosting (OCI, later AWS)

## 📝 Running the Project
1. Install Dependencies
```bash
go mod tidy
```
2. Configure environment variables. 
```bash
cp ./.env.example ./.env
```
3. Run migration
```bash
go run ./cmd/migrate
```
4. Run Bob code generation
```bash
# The package is used for code generation from database schema
# When there is an update in schema, this command shall be executed again
PSQL_DSN=postgres://user:pass@host:port/dbname go run github.com/stephenafamo/bob/gen/bobgen-psql@latest
```
5. Running the application
```bash
go run ./cmd/api
```

## License  
This project is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.  

---
