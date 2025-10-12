#!/bin/bash

set -e

# Configuration
PROJECT_ID="${GCP_PROJECT_ID:-your-project-id}"
REGION="us-central1"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check arguments
if [ $# -lt 1 ]; then
    echo -e "${RED}Usage: $0 <service-name>${NC}"
    echo ""
    echo "Available services:"
    echo "  - daily-trend-analysis-worker"
    echo "  - gateway"
    echo "  - session-service"
    echo "  - summary-service"
    echo "  - trend-service"
    echo "  - user-service"
    echo "  - weekly-trend-analysis-worker"
    exit 1
fi

SERVICE=$1
BUILD_CONFIG="cloudbuild/${SERVICE}.yaml"

# Check if build config exists
if [ ! -f "$BUILD_CONFIG" ]; then
    echo -e "${RED}Error: Build config not found: $BUILD_CONFIG${NC}"
    exit 1
fi

echo -e "${GREEN}Triggering manual build for: ${SERVICE}${NC}"
echo "Project: $PROJECT_ID"
echo "Config: $BUILD_CONFIG"
echo ""

# Submit build
gcloud builds submit \
    --config="$BUILD_CONFIG" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    .

echo -e "${GREEN}✓ Build submitted successfully!${NC}"