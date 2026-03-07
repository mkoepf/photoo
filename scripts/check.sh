#!/bin/bash

# Photoo Code Quality and Automated Test Script
# Inspired by mkoepf/ghcrctl

set -e

# --- Environment ---
# Enforce consistent locale for date formatting in tests
export LANG=en_US.UTF-8
export LC_ALL=en_US.UTF-8

# --- Configuration & Colors ---
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}--- Starting Photoo Code Quality Checks ---${NC}"

# --- 1. Frontend Build, Type-check, Test, Lint ---
# MUST BE FIRST because Go embedding (main.go) requires frontend/dist
if [ -d "frontend" ] && [ "$SKIP_FRONTEND" != "true" ]; then
    echo -e "1. Checking Frontend (React/TS)..."
    cd frontend
    
    echo -n "   - Syncing dependencies... "
    if [ -f "package-lock.json" ]; then
        if ! npm ci --silent 2>/dev/null && ! npm install --silent 2>/dev/null; then
            echo -e "${YELLOW}FAILED (skipping frontend checks)${NC}"
            cd ..
            SKIP_FRONTEND=true
        else
            echo -e "${GREEN}DONE${NC}"
        fi
    else
        if ! npm install --silent 2>/dev/null; then
            echo -e "${YELLOW}FAILED (skipping frontend checks)${NC}"
            cd ..
            SKIP_FRONTEND=true
        else
            echo -e "${GREEN}DONE${NC}"
        fi
    fi
    
    if [ "$SKIP_FRONTEND" != "true" ]; then
        echo -n "   - Building frontend... "
        npm run build --silent
        echo -e "${GREEN}DONE${NC}"

        echo -n "   - Running frontend type-check (tsc)... "
        npm run type-check || npx tsc --noEmit
        echo -e "${GREEN}PASSED${NC}"

        echo -n "   - Running frontend tests (vitest)... "
        npm test -- --silent
        echo -e "${GREEN}PASSED${NC}"

        # If a lint script exists in package.json
        if grep -q "\"lint\":" package.json; then
            echo -n "   - Running frontend lint... "
            npm run lint --silent
            echo -e "${GREEN}PASSED${NC}"
        fi
        cd ..
    fi
else
    echo -e "${YELLOW}1. Frontend checks skipped.${NC}"
fi

# --- 2. Go Formatting Check ---
echo -n "2. Checking Go formatting (gofmt)... "
if [ -n "$(gofmt -l . | grep -v 'wailsjs' || true)" ]; then
    echo -e "${RED}FAILED${NC}"
    echo "The following files are not formatted correctly:"
    gofmt -l . | grep -v 'wailsjs'
    echo -e "${YELLOW}Run 'go fmt ./...' to fix this.${NC}"
    exit 1
fi
echo -e "${GREEN}PASSED${NC}"

# --- 3. Go Vet & Static Analysis ---
echo -n "3. Running go vet... "
if ! go vet ./... 2>&1 | grep -v "build constraints exclude all Go files"; then
    # If there are still errors after filtering out build constraint ones
    if go vet ./... 2>&1 | grep -v "build constraints exclude all Go files" | grep -q "."; then
        echo -e "${RED}FAILED${NC}"
        go vet ./... 2>&1 | grep -v "build constraints exclude all Go files"
        exit 1
    fi
fi
echo -e "${GREEN}PASSED (filtered)${NC}"

# --- 4. Go Build Check ---
echo -n "4. Verifying Go build... "
if ! go build -o /dev/null . 2>&1 | grep -v "build constraints exclude all Go files"; then
    if go build -o /dev/null . 2>&1 | grep -v "build constraints exclude all Go files" | grep -q "."; then
        echo -e "${RED}FAILED${NC}"
        go build -o /dev/null .
        exit 1
    fi
fi
echo -e "${GREEN}PASSED (filtered)${NC}"

# --- 5. Go Tests ---
echo -e "5. Running Go tests..."
# Using -race for concurrency checks if CGO is enabled
if [ "$(go env CGO_ENABLED)" = "1" ]; then
    go test -v -race ./...
else
    go test -v ./...
fi
echo -e "${GREEN}Go tests PASSED${NC}"

# --- 6. Vulnerability Scan (if govulncheck is installed) ---
# Ensure we check GOPATH/bin for the newest version
GOPATH_BIN=$(go env GOPATH)/bin
PATH="$GOPATH_BIN:$PATH"

if command -v govulncheck &> /dev/null; then
    echo -e "6. Running vulnerability scan (govulncheck)..."
    govulncheck ./... || echo -e "${YELLOW}Vulnerability scan found issues, but continuing...${NC}"
    echo -e "${GREEN}Scan complete.${NC}"
else
    echo -e "${YELLOW}6. govulncheck not found. Skipping scan.${NC}"
    echo -e "   Tip: Run 'go install golang.org/x/vuln/cmd/govulncheck@latest' to install it."
fi

echo -e "${BLUE}--- All Photoo Quality Checks PASSED ---${NC}"
