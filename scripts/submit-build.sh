#!/bin/bash

gcloud builds submit user-service/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/user-service:manual-v1

gcloud builds submit trend-service/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/trend-service:manual-v1

gcloud builds submit summary-service/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/summary-service:manual-v1

gcloud builds submit session-service/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/session-service:manual-v1

gcloud builds submit gateway/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/gateway:manual-v1

gcloud builds submit daily-trend-analysis-worker/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/daily-trend-analysis-worker:manual-v1

gcloud builds submit weekly-trend-analysis-worker/  \     
--tag us-central1-docker.pkg.dev/systemic-productivity-planner/prod-planner-repo/weekly-trend-analysis-worker:manual-v1

projects/systemic-productivity-planner/locations/us-central1/connections/MyGithub/repositories/adi290491-productivity-planner