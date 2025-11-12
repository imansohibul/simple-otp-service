# Backend Engineering Interview Assignment (Golang)

### 👨‍💻 Candidate (Author)

**Mochamad Sohibul Iman**  
📧 [iman@imansohibul.my.id](mailto:iman@imansohibul.my.id) 
💼 Candidate: **Backend Engineer** 
--
## 🏗️ Software Architecture

This project follows a **modular clean architecture** pattern. It ensures high maintainability, testability, and clear separation of concerns.

### 🧱 Architectural Layers

| Layer        | Responsibility |
|--------------|----------------|
| **Entity**   | Core domain logic: models and business rules. No framework or external dependency here. |
| **Usecase**  | Orchestrates application flow: how data moves and is transformed. Calls repositories and domain logic. |
| **Repository** | Data persistence and third-party integration. Implements storage logic (PostgreSQL, etc.). |
| **Delivery (REST)** | Handles HTTP requests and responses using Echo. |
---
