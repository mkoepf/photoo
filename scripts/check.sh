#!/bin/bash

# Photoo Code Quality and Automated Test Script
# Inspired by mkoepf/ghcrctl

set -e

# --- Configuration & Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}--- Starting Photoo Code Quality Checks ---${NC}"

# --- 1. Frontend Check & Build ---
# MUST BE FIRST because Go embedding (main.go) requires frontend/dist
if [ -d "frontend" ]; then
    echo -e "1. Checking Frontend (React/TS)..."
    cd frontend
    
    echo -n "   - Syncing dependencies... "
    if [ -f "package-lock.json" ]; then
        npm ci --silent
    else
        npm install --silent
    fi
    echo -e "${GREEN}DONE${NC}"
    
    echo -n "   - Building frontend... "
    npm run build --silent
    echo -e "${GREEN}DONE${NC}"

    echo -n "   - Running frontend type-check (tsc)... "
    npm run type-check || npx tsc --noEmit
    echo -e "${GREEN}PASSED${NC}"

    # If a lint script exists in package.json
    if grep -q "lint" package.json; then
        echo -n "   - Running frontend lint... "
        npm run lint
        echo -e "${GREEN}PASSED${NC}"
    fi
    cd ..
else
    echo -e "${YELLOW}1. Frontend directory not found. Skipping build/checks.${NC}"
fi

# --- 2. Go Formatting Check ---
echo -n "2. Checking Go formatting (gofmt)... "
UNFORMATTED=$(gofmt -l . | grep -v "wailsjs" || true)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}FAILED${NC}"
    echo -e "The following files are not formatted correctly:"
    echo "$UNFORMATTED"
    echo -e "${YELLOW}Run 'go fmt ./...' to fix this.${NC}"
    exit 1
fi
echo -e "${GREEN}PASSED${NC}"

# --- 3. Go Vet & Static Analysis ---
echo -n "3. Running go vet... "
go vet ./...
echo -e "${GREEN}PASSED${NC}"

# --- 4. Go Build Check ---
echo -n "4. Verifying Go build... "
go build -o /dev/null .
echo -e "${GREEN}PASSED${NC}"

# --- 5. Go Tests ---
echo -e "5. Running Go tests..."
# Using -race for concurrency checks
go test -v -race ./...
echo -e "${GREEN}Go tests PASSED${NC}"

# --- 6. Vulnerability Scan (if govulncheck is installed) ---
if command -v govulncheck &> /dev/null; then
    echo -e "6. Running vulnerability scan (govulncheck)..."
    govulncheck ./... || echo -e "${YELLOW}Vulnerability scan found issues, but continuing...${NC}"
    echo -e "${GREEN}Scan complete.${NC}"
else
    echo -e "${YELLOW}6. govulncheck not found. Skipping scan.${NC}"
fi

echo -e "${BLUE}--- All Photoo Quality Checks PASSED ---${NC}"
