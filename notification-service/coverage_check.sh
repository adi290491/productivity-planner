#!/bin/bash

cd /Users/adityasawant/Documents/Projects/golang/productivity-planner/productivity-planner/notification-service

echo "🧪 Running tests with coverage analysis..."
echo "========================================"

# Run tests with coverage
go test -coverprofile=coverage.out ./cmd ./notification ./models

echo ""
echo "📊 Function-level Coverage Details:"
echo "====================================="

# Show detailed function coverage
go tool cover -func=coverage.out

echo ""
echo "📈 Package Coverage Summary:"
echo "============================"

# Extract and display package summaries
go tool cover -func=coverage.out | grep -E "total|cmd/|notification/|models/" | tail -10

echo ""
echo "🎯 Overall Coverage:"
echo "===================="

# Calculate overall coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "Total Coverage: $COVERAGE"

# Determine if coverage meets target
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
if (( $(echo "$COVERAGE_NUM >= 70" | bc -l) )); then
    echo "✅ Coverage target met (≥70%)"
else 
    echo "❌ Coverage below target (<70%)"
fi

echo ""
echo "📄 Generating HTML report..."
go tool cover -html=coverage.out -o coverage.html
echo "HTML coverage report: coverage.html"
