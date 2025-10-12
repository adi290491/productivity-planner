#!/bin/bash

set -x
set -e

# Configuration
PROJECT_ID="${GCP_PROJECT_ID:-systemic-productivity-planner}"
REGION="us-central1"
REPO_NAME="productivity-planner"
GITHUB_CONNECTION="MyGithub"
GITHUB_REPO_RESOURCE="adi290491-productivity-planner"
BRANCH_PATTERN="${BRANCH_PATTERN:-^master$}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Setting up Cloud Build triggers for productivity-planner...${NC}"
echo "Project ID: $PROJECT_ID"
echo "Region: $REGION"
echo "Repository: $GITHUB_OWNER/$GITHUB_REPO"
echo "Branch Pattern: $BRANCH_PATTERN"
echo "GitHub Connection: $GITHUB_CONNECTION"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}Error: gcloud CLI is not installed${NC}"
    exit 1
fi

# Set project
echo -e "${YELLOW}Setting GCP project...${NC}"
gcloud config set project "$PROJECT_ID"

# Enable required APIs
echo -e "${YELLOW}Enabling required APIs...${NC}"
gcloud services enable cloudbuild.googleapis.com \
    artifactregistry.googleapis.com \
    run.googleapis.com \
    secretmanager.googleapis.com

# Verify GitHub connection exists
echo -e "${YELLOW}Verifying GitHub connection...${NC}"
if ! gcloud builds connections describe "$GITHUB_CONNECTION" --region="$REGION" &>/dev/null; then
    echo -e "${RED}Error: GitHub connection '$GITHUB_CONNECTION' not found${NC}"
    echo -e "${YELLOW}Please create a connection first:${NC}"
    echo "  1. Visit: https://console.cloud.google.com/cloud-build/triggers/connect?project=$PROJECT_ID"
    echo "  2. Connect your GitHub repository"
    echo "  3. Note the connection name and update this script"
    echo ""
    echo "Or list existing connections:"
    echo "  gcloud builds connections list --region=$REGION"
    exit 1
fi

echo -e "${GREEN}✓ GitHub connection verified${NC}"

# Create Artifact Registry repository if it doesn't exist
echo -e "${YELLOW}Creating Artifact Registry repository...${NC}"
gcloud artifacts repositories describe "$REPO_NAME" \
    --location="$REGION" &>/dev/null || \
gcloud artifacts repositories create "$REPO_NAME" \
    --repository-format=docker \
    --location="$REGION" \
    --description="Docker repository for productivity-planner services"

# Grant Cloud Build service account permissions
echo -e "${YELLOW}Granting Cloud Build permissions...${NC}"
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")
CLOUD_BUILD_SA="${PROJECT_NUMBER}@cloudbuild.gserviceaccount.com"

gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${CLOUD_BUILD_SA}" \
    --role="roles/artifactregistry.writer" \
    --condition=None

# Services to create triggers for
declare -a SERVICES=(
    "daily-trend-analysis-worker"
    "gateway"
    "session-service"
    "summary-service"
    "trend-service"
    "user-service"
    "weekly-trend-analysis-worker"
)

# Create triggers for each service
echo -e "${GREEN}Creating Cloud Build triggers...${NC}"
echo ""

for SERVICE in "${SERVICES[@]}"; do
    TRIGGER_NAME="${SERVICE}-trigger"
    
    echo -e "${YELLOW}Creating trigger for: ${SERVICE}${NC}"
    
    # Delete existing trigger if it exists
    gcloud builds triggers delete "$TRIGGER_NAME" --region="$REGION" --quiet 2>/dev/null || true
    
    # Create new trigger using the 2nd generation (repository resource)
    gcloud builds triggers create github \
        --name="$TRIGGER_NAME" \
        --region="$REGION" \
        --repository="projects/systemic-productivity-planner/locations/us-central1/connections/MyGithub/repositories/adi290491-productivity-planner" \
        --branch-pattern="$BRANCH_PATTERN" \
        --build-config="cloudbuild/${SERVICE}.yaml" \
        --included-files="${SERVICE}/**,cloudbuild/${SERVICE}.yaml" \
        --include-logs-with-status \
        --substitutions="_SERVICE_NAME=$SERVICE" \
        --description="Automated build for ${SERVICE} on push to branch"
    # gcloud builds triggers create github \
    #     --name="$TRIGGER_NAME" \
    #     --region="$REGION" \
    #     --repo-name="$GITHUB_REPO" \
    #     --repo-owner="$GITHUB_OWNER" \
    #     --branch-pattern="$BRANCH_PATTERN" \
    #     --build-config="cloudbuild/${SERVICE}.yaml" \
    #     --include-logs-with-status \
    #     --substitutions="_SERVICE_NAME=$SERVICE" \
    #     --description="Automated build for ${SERVICE} on push to branch"
    
    echo -e "${GREEN}✓ Created trigger: ${TRIGGER_NAME}${NC}"
    echo ""
done

echo -e "${GREEN}All triggers created successfully!${NC}"
echo ""
echo -e "${YELLOW}View your triggers:${NC}"
echo "  https://console.cloud.google.com/cloud-build/triggers?project=$PROJECT_ID"
echo ""
echo -e "${YELLOW}Test a trigger:${NC}"
echo "  git add ."
echo "  git commit -m 'Test trigger'"
echo "  git push origin master"
echo ""
echo -e "${GREEN}Setup complete!${NC}"