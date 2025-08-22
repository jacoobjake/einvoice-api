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
- /cmd/api # API entry point
- /internal
- /auth # Auth related handlers and services
- /invoice # Handlers, services, repository for invoice
- /user # Handlers, services, repository for user
- /migrations # Database migrations

## 🛠️ Planned Tech Stack
- Go
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- PostgreSQL
- REST API
- Docker (planned)
- Cloud hosting (OCI, later AWS)

## License  
This project is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.  

---
