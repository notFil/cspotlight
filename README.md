# 🌐 cspotlight – Content Security Policy (CSP) Violation Collector & Analyzer

**cspotlight** is an open-source CSP report collector and analyzer designed to help security engineers and developers monitor, analyze, and respond to browser-enforced Content Security Policy violations.

---

## 🚀 Features

- 📥 **CSP Report Ingestion**: Receive and store incoming CSP violation reports via a secure API.
- 📊 **Analytics Dashboard**: Visualize violations by source, type, page, and frequency.
- 🎯 **Severity Scoring**: Highlight suspicious or potentially dangerous violations (e.g., `eval`, inline scripts, third-party JS).
- 🕵️ **Enrichment**: Optional GeoIP lookup, user-agent parsing, and referrer analysis.
- 🔔 **Alerts & Integrations**: Configurable Slack/email alerting for high-risk reports.
- 🐳 **Dockerized & Deployable**: Run locally or deploy via Docker, Kubernetes, or serverless platforms.

---

## 📸 Screenshots

> _Coming soon: screenshots of the dashboard and sample alert notifications._

---

## 🛠️ Getting Started

### Prerequisites

- Python 3.13+
- Docker (optional for container-based deployment)

### Installation (Local)

```bash
git clone https://github.com/notFil/cspotlight.git
cd cspotlight
<TODO>