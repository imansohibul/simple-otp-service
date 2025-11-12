# Backend Engineering Interview Assignment (Golang)

## Author

**Name:** Mochamad Sohibul Iman  
**Email:** [iman@imansohibul.my.id](mailto:iman@imansohibul.my.id)  
**LinkedIn:** [www.linkedin.com/in/imansohibul](https://www.linkedin.com/in/imansohibul)

## 🏗️ Project Structure & Software Architecture

This project follows a **modular clean architecture** pattern. It ensures high maintainability, testability, and clear separation of concerns.


### 👨‍💻 Candidate (Author)

**Mochamad Sohibul Iman**  
📧 [iman@imansohibul.my.id](mailto:iman@imansohibul.my.id) 
💼 Candidate: **Backend Engineer** 
--

### 🧱 Architectural Layers

| Layer        | Responsibility |
|--------------|----------------|
| **Entity**   | Core domain logic: models and business rules. No framework or external dependency here. |
| **Usecase**  | Orchestrates application flow: how data moves and is transformed. Calls repositories and domain logic. |
| **Repository** | Data persistence and third-party integration. Implements storage logic (PostgreSQL, etc.). |
| **Delivery (REST)** | Handles HTTP requests and responses using Echo. Maps JSON ↔ DTO ↔ Entities. |
| **Config**   | Dependency wiring (DI), configuration loading, and server setup. |
| **DB Migrate** | Database version control using SQL migrations. |

---

