# 🧠 Systemic Productivity Planner

A cloud-native productivity tracking platform inspired by the Pomodoro technique, designed to help users understand how they spend their time, identify patterns in focus and distraction, and gradually improve productivity through data-driven insights.

This project is intentionally backend- and infrastructure-focused, showcasing scalable service design, cloud automation, and production-ready patterns, with room for future AI-powered insights.

---

## ✨ Key Features

### Core Functionality

**⏱️ Session Tracking**
- Track Focus, Meeting, and Break sessions
- Start/stop sessions with validation to prevent overlaps

**📊 Daily & Weekly Trend Analysis**
- Background workers compute aggregated productivity trends
- Daily summaries (per day)
- Weekly summaries (per week, with daily breakdown)

**🧩 API Gateway**
- Central entry point for frontend communication
- JWT-based authentication
- Service-to-service routing

**🔐 Authentication**
- Secure login & signup
- JWT authentication

**🌐 Frontend**
- React + Vite SPA
- Hosted on Firebase Hosting
- Consumes backend APIs via Gateway

---

## 🏗️ Architecture

### High-Level System Design

```
[ Frontend (Firebase Hosting) ]
              |
              v
        [ API Gateway ]
              |
  --------------------------------
  |        |        |            |
  v        v        v            v
User   Session   Summary     Trend
Service Service  Service     Service
              |
        [ Cloud SQL (Postgres) ]
              |
        [ Trend Workers ]
        (Daily / Weekly Jobs)
```

### Background Processing

- **Daily Trend Analysis Worker**
- **Weekly Trend Analysis Worker**
- Executed as Cloud Run Jobs
- Triggered via Cloud Run Job Schedulers

---

## ☁️ Cloud Infrastructure

**Platform:** Google Cloud Platform (GCP)

### Services Used

- **Cloud Run** - Services & Jobs
- **Cloud SQL** - PostgreSQL database
- **Artifact Registry** - Container images
- **Cloud Build** - CI/CD automation
- **Secret Manager** - Secure credential storage
- **Firebase Hosting** - Frontend deployment

### Networking

- Direct VPC Egress for Cloud SQL access

### Security

- IAM-based service accounts
- Secrets injected via Secret Manager
- No secrets committed to source control

---

## 📁 Repository Structure

```
.
├── gateway/
├── user-service/
├── session-service/
├── summary-service/
├── trend-service/
├── daily-trend-analysis-worker/
├── weekly-trend-analysis-worker/
├── productivity-frontend/
├── cloudbuild/
│   ├── gateway.yaml
│   ├── user-service.yaml
│   ├── daily-worker.yaml
│   └── weekly-worker.yaml
├── cloudrun/
│   ├── gateway.yaml
│   ├── user-service.yaml
│   └── ...
└── README.md
```

---

## ⚙️ Configuration & Environment Management

### Philosophy

- No complex config files in production
- Services read directly from environment variables
- `godotenv` is used only for local development

---

## 🔁 CI/CD Strategy

### Deployment Pipeline

- **One Cloud Build trigger per service**
- Builds and deploys only when relevant files change

### Each Service Includes

- Its own `cloudbuild/*.yaml`
- Its own `cloudrun/*.yaml`

### Example Flow

1. Push to `main`
2. Cloud Build trigger detects service changes
3. Build Docker image
4. Push to Artifact Registry
5. Deploy/update Cloud Run service or job

---

## 📊 Trend Analysis Model

### Daily Summary

- Total productive time per day
- Breakdown by session type (Focus / Meeting / Break)

### Weekly Summary

- Aggregated weekly total
- Includes daily summaries for context
- Enables trend comparison week-over-week

---

## 🔔 Notifications (Current Approach)

To keep the system simple and production-ready:

- **No email notifications yet**
- Frontend polls for availability of new trends
- Backend exposes lightweight endpoints to check trend freshness
- Failure cases are logged centrally for admin review

This avoids premature complexity while keeping the system extensible.

---

## 🚀 Running Locally

### Backend (example: user-service)

```bash
cd user-service
cp .env.example .env
go run ./cmd
```

### Frontend

```bash
cd productivity-frontend
npm install
npm run dev
```

---

## 🧪 Testing

- Unit tests for core business logic
- `go test ./...` for backend services
- CI enforces tests before deployment

---

## 🧠 Future Roadmap

### Planned Features

**AI-Assisted Productivity Insights:**
- Personalized focus recommendations
- Consistency nudges
- Goal breakdown suggestions

**Additional Enhancements:**
- Advanced notification channels (email, Slack)
- Analytics dashboards
- Role-based access (admin vs user)

---

## 🎯 Design Philosophy

- **Backend & infrastructure first**
- **Production-grade patterns** over quick hacks
- **Minimal features, done well**
- **Scalable foundations** before AI augmentation

---

## 📌 Status

🟢 **Actively developed**  
🧱 **Core backend & infrastructure complete**  
🚧 **AI features planned next**