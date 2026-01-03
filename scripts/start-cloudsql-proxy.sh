#!/bin/bash

echo "Starting Cloud SQL Proxy for planner-prod..."
echo ""
echo "pgAdmin Connection Info:"
echo "  Host:     localhost"
echo "  Port:     5433"
echo "  Database: $(gcloud secrets versions access latest --secret=DB_NAME 2>/dev/null || echo '<check secrets>')"
echo "  Username: $(gcloud secrets versions access latest --secret=DB_USERNAME 2>/dev/null || echo '<check secrets>')"
echo ""
echo "Press Ctrl+C to stop the proxy"
echo ""

cloud-sql-proxy systemic-productivity-planner:us-central1:planner-prod --port=5433