#!/bin/bash

# Script to reorganize go-mvc codebase following Clean Architecture best practices
# Author: GitHub Copilot
# Date: 2024

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base directory
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$BASE_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Go-MVC Code Reorganization Script${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to confirm action
confirm() {
    read -p "$(echo -e ${YELLOW}[CONFIRM]${NC}) $1 (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Skipped."
        return 1
    fi
    return 0
}

# Function to create backup
create_backup() {
    print_info "Creating backup..."
    BACKUP_DIR="backup_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    cp -r internal "$BACKUP_DIR/"
    print_success "Backup created at: $BACKUP_DIR"
}

# Function to check if git is clean
check_git_status() {
    if command -v git &> /dev/null; then
        if [[ -n $(git status -s) ]]; then
            print_warning "You have uncommitted changes."
            if ! confirm "Do you want to continue anyway?"; then
                print_error "Aborting. Please commit or stash your changes first."
                exit 1
            fi
        fi
    fi
}

# Phase 1: Rename core/domain to domain
phase1_rename_domain() {
    echo -e "\n${GREEN}=== Phase 1: Rename core/domain to domain ===${NC}\n"
    
    if confirm "Rename internal/core/domain/ to internal/domain/?"; then
        print_info "Renaming directories..."
        
        # Check if internal/domain already exists (might have legacy code)
        if [ -d "internal/domain" ]; then
            print_warning "internal/domain already exists!"
            if confirm "Move existing internal/domain to internal/domain_legacy?"; then
                mv internal/domain internal/domain_legacy
                print_success "Moved to internal/domain_legacy"
            else
                print_error "Cannot proceed with existing internal/domain"
                return 1
            fi
        fi
        
        # Move core/domain to domain
        if [ -d "internal/core/domain" ]; then
            mv internal/core/domain internal/domain
            print_success "Renamed internal/core/domain → internal/domain"
        fi
        
        # Move core/ports to domain
        if [ -d "internal/core/ports" ]; then
            print_info "Merging internal/core/ports into domain..."
            
            # Move repositories
            if [ -d "internal/core/ports/repositories" ]; then
                mkdir -p internal/domain/_legacy_ports
                mv internal/core/ports/repositories internal/domain/_legacy_ports/
            fi
            
            # Move services
            if [ -d "internal/core/ports/services" ]; then
                mkdir -p internal/domain/_legacy_ports
                mv internal/core/ports/services internal/domain/_legacy_ports/
            fi
            
            print_success "Moved ports to internal/domain/_legacy_ports for manual migration"
        fi
        
        # Clean up empty core directory
        if [ -d "internal/core" ] && [ -z "$(ls -A internal/core)" ]; then
            rmdir internal/core
            print_success "Removed empty internal/core directory"
        fi
        
        print_success "Phase 1 completed!"
    fi
}

# Phase 2: Rename handlers to interfaces
phase2_rename_handlers() {
    echo -e "\n${GREEN}=== Phase 2: Rename handlers to interfaces ===${NC}\n"
    
    if confirm "Rename internal/handlers/ to internal/interfaces/?"; then
        print_info "Renaming directories..."
        
        if [ -d "internal/handlers" ]; then
            mv internal/handlers internal/interfaces
            print_success "Renamed internal/handlers → internal/interfaces"
        else
            print_warning "internal/handlers not found, skipping..."
        fi
        
        print_success "Phase 2 completed!"
    fi
}

# Phase 3: Consolidate infrastructure
phase3_consolidate_infrastructure() {
    echo -e "\n${GREEN}=== Phase 3: Consolidate Infrastructure ===${NC}\n"
    
    if confirm "Merge internal/adapters/ into internal/infrastructure/?"; then
        print_info "Consolidating infrastructure..."
        
        if [ -d "internal/adapters" ]; then
            # Create infrastructure directory if it doesn't exist
            mkdir -p internal/infrastructure
            
            # Move each subdirectory
            for dir in internal/adapters/*/; do
                if [ -d "$dir" ]; then
                    dirname=$(basename "$dir")
                    if [ -d "internal/infrastructure/$dirname" ]; then
                        print_warning "internal/infrastructure/$dirname already exists"
                        print_info "Moving to internal/infrastructure/${dirname}_from_adapters"
                        mv "$dir" "internal/infrastructure/${dirname}_from_adapters"
                    else
                        mv "$dir" "internal/infrastructure/"
                        print_success "Moved adapters/$dirname → infrastructure/$dirname"
                    fi
                fi
            done
            
            # Remove empty adapters directory
            if [ -z "$(ls -A internal/adapters)" ]; then
                rmdir internal/adapters
                print_success "Removed empty internal/adapters directory"
            fi
        else
            print_warning "internal/adapters not found, skipping..."
        fi
        
        # Move shared utilities
        if [ -d "internal/shared" ]; then
            print_info "Moving shared utilities to infrastructure..."
            
            for dir in internal/shared/*/; do
                if [ -d "$dir" ]; then
                    dirname=$(basename "$dir")
                    if [ -d "internal/infrastructure/$dirname" ]; then
                        print_warning "internal/infrastructure/$dirname already exists"
                    else
                        mv "$dir" "internal/infrastructure/"
                        print_success "Moved shared/$dirname → infrastructure/$dirname"
                    fi
                fi
            done
            
            # Remove empty shared directory
            if [ -z "$(ls -A internal/shared)" ]; then
                rmdir internal/shared
                print_success "Removed empty internal/shared directory"
            fi
        fi
        
        print_success "Phase 3 completed!"
    fi
}

# Phase 4: Update imports
phase4_update_imports() {
    echo -e "\n${GREEN}=== Phase 4: Update Import Paths ===${NC}\n"
    
    if confirm "Update import paths in Go files?"; then
        print_info "Updating imports (this may take a while)..."
        
        # Get module name from go.mod
        MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}')
        print_info "Module name: $MODULE_NAME"
        
        # Update imports for core/domain → domain
        print_info "Updating internal/core/domain imports..."
        find . -name "*.go" -type f -exec sed -i.bak \
            "s|$MODULE_NAME/internal/core/domain|$MODULE_NAME/internal/domain|g" {} \;
        
        # Update imports for handlers → interfaces
        print_info "Updating internal/handlers imports..."
        find . -name "*.go" -type f -exec sed -i.bak \
            "s|$MODULE_NAME/internal/handlers|$MODULE_NAME/internal/interfaces|g" {} \;
        
        # Update imports for adapters → infrastructure
        print_info "Updating internal/adapters imports..."
        find . -name "*.go" -type f -exec sed -i.bak \
            "s|$MODULE_NAME/internal/adapters|$MODULE_NAME/internal/infrastructure|g" {} \;
        
        # Update imports for shared → infrastructure
        print_info "Updating internal/shared imports..."
        find . -name "*.go" -type f -exec sed -i.bak \
            "s|$MODULE_NAME/internal/shared|$MODULE_NAME/internal/infrastructure|g" {} \;
        
        # Clean up backup files
        find . -name "*.go.bak" -type f -delete
        
        print_success "Import paths updated!"
        print_success "Phase 4 completed!"
    fi
}

# Phase 5: Format and verify
phase5_verify() {
    echo -e "\n${GREEN}=== Phase 5: Format & Verify ===${NC}\n"
    
    print_info "Running gofmt..."
    gofmt -w .
    print_success "Code formatted!"
    
    print_info "Running go mod tidy..."
    go mod tidy
    print_success "Dependencies cleaned!"
    
    print_info "Building project..."
    if go build ./...; then
        print_success "✓ Build successful!"
    else
        print_error "✗ Build failed! Please check the errors above."
        return 1
    fi
    
    print_info "Running tests..."
    if go test ./... -short; then
        print_success "✓ Tests passed!"
    else
        print_warning "⚠ Some tests failed. Please review."
    fi
    
    print_success "Phase 5 completed!"
}

# Main execution
main() {
    echo -e "${BLUE}This script will reorganize your codebase.${NC}"
    echo -e "${YELLOW}It's recommended to commit your changes first!${NC}\n"
    
    check_git_status
    
    if confirm "Do you want to create a backup?"; then
        create_backup
    fi
    
    echo -e "\n${BLUE}Available phases:${NC}"
    echo "  1. Rename core/domain to domain"
    echo "  2. Rename handlers to interfaces"
    echo "  3. Consolidate infrastructure"
    echo "  4. Update import paths"
    echo "  5. Format and verify"
    echo "  6. Run all phases"
    echo ""
    
    read -p "$(echo -e ${BLUE}[SELECT]${NC}) Which phase do you want to run? (1-6): " phase
    
    case $phase in
        1)
            phase1_rename_domain
            ;;
        2)
            phase2_rename_handlers
            ;;
        3)
            phase3_consolidate_infrastructure
            ;;
        4)
            phase4_update_imports
            ;;
        5)
            phase5_verify
            ;;
        6)
            print_info "Running all phases..."
            phase1_rename_domain
            phase2_rename_handlers
            phase3_consolidate_infrastructure
            phase4_update_imports
            phase5_verify
            ;;
        *)
            print_error "Invalid option!"
            exit 1
            ;;
    esac
    
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}  Reorganization completed!${NC}"
    echo -e "${GREEN}========================================${NC}\n"
    
    echo -e "${BLUE}Next steps:${NC}"
    echo "  1. Review changes: git status"
    echo "  2. Check if everything builds: go build ./..."
    echo "  3. Run tests: go test ./..."
    echo "  4. Review documentation: docs/CURRENT_STRUCTURE_ANALYSIS.md"
    echo "  5. Commit changes: git add . && git commit -m 'refactor: reorganize codebase structure'"
    echo ""
}

# Run main function
main "$@"
