# 🔐 Korun.io: Open-Source Secret Management for Modern DevOps

## THIS IS NOT COMPLETE!

## 🛡️ Korun.io - Secured. Integrated. Open.

`korun.io` (derived from "Korunmak", Turkish for "to be secured/protected") is an open-source, self-hostable secret management platform designed for the modern DevOps landscape. Built from the ground up with performance, security, and developer experience in mind, Korun.io aims to provide a robust, easy-to-use alternative to existing solutions, empowering teams to manage their environment variables and sensitive credentials with confidence.

## ✨ Why Korun.io?

The digital world demands robust secret management, but existing solutions often come with significant operational overhead, steep learning curves, or prohibitive costs for smaller teams and open-source projects.

Korun.io addresses these challenges by offering:

- **Self-Hostable by Design**: Full control over your data, ideal for organizations with strict compliance requirements or those preferring on-premise solutions.
- **DevOps-Friendly**: Built with GitOps, Kubernetes, and CLI-first workflows in mind for seamless integration into your existing CI/CD pipelines.
- **Performance & Scalability**: Engineered with Go, Kafka, and a cloud-native architecture to handle high-throughput demands from microservices at scale.
- **Transparent & Secure**: Open-source codebase allows for full auditing and community-driven security enhancements.
- **Cost-Effective**: Free to use and self-host, with future plans for optional managed hosting for those who prefer convenience over self-management.

## 🚀 Key Features (Planned & In Progress)

### Core Secret Management

- Hierarchical Organization:
  - Workspaces: Top-level containers for your entire organization or business units.
  - Projects: Group related applications or services within a Workspace.
  - Environments: Manage environment-specific secrets (e.g., development, staging, production, test).
  - Secrets: Key-value pairs for environment variables (DB_URL, API_KEY, etc.).
- Versioning & Rollback: Track changes to secrets with a full audit trail and easily revert to previous versions.
- Dynamic Secrets (Future): Integrate with cloud providers (AWS, Azure, GCP) to generate short-lived credentials.
- Secret Rotation (Future): Automated rotation of database passwords, API keys, and other credentials.

### Developer & DevOps Experience

- Powerful CLI: Fetch, set, and inject secrets directly into your shell or applications.
- Kubernetes Operator: Seamlessly inject secrets into Kubernetes Pods as environment variables or mounted files, removing the need for manual kubectl commands.
- Next.js Web Dashboard: Intuitive UI for secret management, user roles, and audit log exploration.
- GitOps Integration: Define secrets and policies as code, enabling declarative secret management via your Git repository.
- Webhooks: Trigger custom actions or notify external systems on secret changes.

### Security & Compliance

- End-to-End Encryption: Secrets are encrypted at rest and in transit.
- Role-Based Access Control (RBAC): Fine-grained permissions for users, teams, workspaces, projects, and environments.
- Audit Logging: Comprehensive, immutable logs of all secret access and modification events.
- Secret Scanners: Integrate with tools to prevent accidental secret exposure in codebases.

## 🎯 Architecture & Technologies

Korun.io is built on a modern, cloud-native architecture designed for resilience, scalability, and maintainability.

- Backend Microservices (Go):
  - API Gateway: Central entry point, handles authentication, rate limiting, and request routing.
  - Auth Service: User authentication, authorization, and RBAC.
  - Secrets Service: Core logic for secret CRUD, encryption, and versioning.
  - Audit Service: Ingests, processes, and exposes audit logs.
  - Notification Service: Handles external communications (e.g., email alerts).
  - Webhook Service: Manages and dispatches webhooks.
- Event Streaming (Apache Kafka):
  - The central nervous system for asynchronous communication between microservices.
  - Enables real-time audit logging, change notifications, and eventual consistency.
  - Provides durable storage for events, supporting stream processing and replay for analytics and compliance.
- Databases:
  - PostgreSQL: Primary data store for metadata (users, projects, environments, secret versions). Selected for its ACID compliance, extensibility, and proven reliability.
  - Redis: Caching layer for performance optimization and session management.
  - Elasticsearch (Future): For efficient, full-text search and analysis of audit logs.
- Frontend (Next.js & React):
  - Responsive and intuitive web interface for managing your Korun.io instance.
- CLI (Go & Cobra):
  - Cross-platform command-line tool for seamless integration with CI/CD and developer workflows.
- Containerization (Docker):
  - All services are containerized for consistent deployment across environments.
- Orchestration (Kubernetes):
  - Cloud-native deployment and management of microservices.
  - Leverages Kubernetes Operators for simplified secret injection.
- Infrastructure as Code (Terraform):
  - Define and provision infrastructure resources (databases, Kafka clusters, etc.) declaratively.

## 🛣️ Roadmap & Future Vision

### Phase 1: Core Functionality (Current Focus)

- Workspaces, Projects, Environments, Secrets CRUD: Establish the core hierarchy.
- Basic User Authentication & Authorization: Local authentication.
- Web Dashboard: Basic secret listing, creation, and editing.
- CLI: get, set, list, and run commands for environment variable injection.
- Containerized Deployment: Docker Compose for local development, basic Kubernetes manifests.
- Kafka Integration: Initial event publishing for audit trails.

### Phase 2: Enhanced Features & DevOps Integration

- Advanced RBAC: Granular permissions, team management.
- Secret Versioning & Rollback: UI/CLI support for history.
- Kubernetes Operator: Automated secret injection.
- Audit Service: Dedicated API for querying historical events.
- SSO Integration: OAuth2, OpenID Connect support.
- Webhooks: Configurable event notifications.

### Phase 3: Enterprise & Advanced Capabilities

- Dynamic Secrets: Integration with cloud provider key management services (KMS) and identity providers.
- Secret Rotation: Automated, scheduled secret updates.
- Disaster Recovery & High Availability: Cross-region deployment strategies.
- Cloud Provider Managed Hosting: Offer Korun.io as a hosted service.
- Client Libraries: Language-specific SDKs (Python, Node.js, Java) for programmatic access.

## 🤝 Contributing

Korun.io is an open-source project, and we welcome contributions from the community! Whether you're a Go backend enthusiast, a React frontend wizard, a Kubernetes operator expert, or a documentation guru, there's a place for you. Please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to get started, set up your development environment, and submit pull requests.

## 📄 License

Korun.io is released under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## ⭐ Star Us!

If you find Korun.io useful, please consider giving us a star on GitHub! It helps a lot.
