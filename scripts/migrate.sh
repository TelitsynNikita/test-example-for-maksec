#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo -e "${RED} Error: Database configuration not set in .env${NC}"
    exit 1
fi

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE:-disable}"

echo -e "${GREEN} Database URL: postgres://${DB_USER}:****@${DB_HOST}:${DB_PORT}/${DB_NAME}${NC}"

if ! command -v migrate &> /dev/null; then
    echo -e "${YELLOW}  migrate not found. Installing...${NC}"
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

COMMAND=${1:-up}

case $COMMAND in
    up)
        echo -e "${GREEN} Applying migrations...${NC}"
        migrate -path migrations -database "$DB_URL" up
        echo -e "${GREEN} Migrations applied successfully${NC}"
        ;;
    down)
        echo -e "${YELLOW} Rolling back migration...${NC}"
        migrate -path migrations -database "$DB_URL" down 1
        echo -e "${GREEN} Migration rolled back${NC}"
        ;;
    down-all)
        echo -e "${RED} Rolling back ALL migrations...${NC}"
        migrate -path migrations -database "$DB_URL" down -all
        echo -e "${GREEN} All migrations rolled back${NC}"
        ;;
    status|version)
        echo -e "${GREEN} Migration status:${NC}"
        migrate -path migrations -database "$DB_URL" version
        ;;
    create)
        if [ -z "$2" ]; then
            echo -e "${RED} Please provide migration name: ./scripts/migrate.sh create <name>${NC}"
            exit 1
        fi
        echo -e "${GREEN} Creating migration: $2${NC}"
        migrate create -ext sql -dir migrations -seq "$2"
        ;;
    *)
        echo -e "${RED} Unknown command: $COMMAND${NC}"
        echo "Usage: ./scripts/migrate.sh [up|down|down-all|status|create <name>]"
        exit 1
        ;;
esac