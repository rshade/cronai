#!/bin/bash
set -eo pipefail

# generate_changelog.sh - Generate a changelog from conventional commits using git-chglog
# Usage: ./scripts/generate_changelog.sh [<from-tag>] [<to-ref>]
#   <from-tag>: Starting tag (default: auto-detect)
#   <to-ref>:   Ending reference (default: HEAD)

COLOR_GREEN="\033[0;32m"
COLOR_YELLOW="\033[0;33m"
COLOR_BLUE="\033[0;34m"
COLOR_PURPLE="\033[0;35m"
COLOR_RESET="\033[0m"

# Ensure git-chglog is installed
if ! command -v git-chglog &> /dev/null; then
    echo -e "${COLOR_YELLOW}git-chglog not found. Please install with: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest${COLOR_RESET}"
    exit 1
fi

# Check if there's a config file
CONFIG_FILE="$(dirname "$0")/git-chglog-config.yml"
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${COLOR_YELLOW}Configuration file not found at $CONFIG_FILE${COLOR_RESET}"
    exit 1
fi

# Get parameters
FROM_TAG=$1
TO_REF=${2:-HEAD}

# Display range information
if [ -n "$FROM_TAG" ]; then
    echo -e "${COLOR_GREEN}Generating changelog from $FROM_TAG to $TO_REF${COLOR_RESET}"
    RANGE="$FROM_TAG..$TO_REF"
else
    echo -e "${COLOR_GREEN}Generating changelog with auto-detected range${COLOR_RESET}"
    RANGE=""
fi

# Determine the future version (assuming semver)
NEXT_VERSION=""
if [ "$TO_REF" = "HEAD" ]; then
    # Get the latest tag or default to v0.0.0 if none exists
    if git tag -l | grep -q .; then
        LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null)
    else
        LATEST_TAG="v0.0.0"
        # If no tags exist, use the first commit as the start
        FIRST_COMMIT=$(git rev-list --max-parents=0 HEAD)
        echo -e "${COLOR_YELLOW}No tags found. Starting from first commit: ${FIRST_COMMIT:0:8}${COLOR_RESET}"
    fi
    
    # Strip the 'v' prefix if it exists
    LATEST_VERSION=${LATEST_TAG#v}
    
    # Parse major, minor, patch
    IFS='.' read -r MAJOR MINOR PATCH <<< "$LATEST_VERSION"
    
    # Determine the commit range for analysis
    if [ -n "$FROM_TAG" ]; then
        ANALYZE_RANGE="$FROM_TAG..$TO_REF"
    else
        if [ "$LATEST_TAG" = "v0.0.0" ]; then
            # If no tags exist, analyze all commits
            ANALYZE_RANGE="$FIRST_COMMIT..$TO_REF"
        else
            ANALYZE_RANGE="$LATEST_TAG..$TO_REF"
        fi
    fi
    
    # Check for breaking changes to bump major
    if git log "$ANALYZE_RANGE" --pretty=format:"%B" 2>/dev/null | grep -q "BREAKING CHANGE"; then
        NEXT_VERSION="v$((MAJOR + 1)).0.0"
    # Check for features to bump minor
    elif git log "$ANALYZE_RANGE" --pretty=format:"%s" 2>/dev/null | grep -q "^feat"; then
        NEXT_VERSION="v$MAJOR.$((MINOR + 1)).0"
    # Otherwise bump patch
    else
        NEXT_VERSION="v$MAJOR.$MINOR.$((PATCH + 1))"
    fi
    
    echo -e "${COLOR_BLUE}Projected next version: ${NEXT_VERSION}${COLOR_RESET}"
    echo
fi

# Generate changelog using git-chglog
echo -e "${COLOR_PURPLE}Generating changelog...${COLOR_RESET}"

# Run git-chglog with the appropriate options
if [ -n "$FROM_TAG" ]; then
    # If a specific range is provided
    git-chglog --config "$CONFIG_FILE" --template "$(dirname "$0")/CHANGELOG.tpl.md" --output CHANGELOG.md "$RANGE"
else
    # Check if there are any tags
    if git tag -l | grep -q .; then
        # Auto-detect range with existing tags
        git-chglog --config "$CONFIG_FILE" --template "$(dirname "$0")/CHANGELOG.tpl.md" --output CHANGELOG.md
    else
        # No tags exist, create changelog from all commits
        FIRST_COMMIT=$(git rev-list --max-parents=0 HEAD)
        echo -e "${COLOR_YELLOW}No tags found. Creating changelog from all commits${COLOR_RESET}"
        
        # For repositories without tags, create a simple changelog manually
        echo "# Changelog" > CHANGELOG.md
        echo "" >> CHANGELOG.md
        echo "## ${NEXT_VERSION:-v0.1.0} ($(date +"%Y-%m-%d"))" >> CHANGELOG.md
        echo "" >> CHANGELOG.md
        
        # Function to process commit types
        process_commit_type() {
            local type=$1
            local title=$2
            
            echo "### $title" >> CHANGELOG.md
            local commits=$(git log --pretty=format:"- %s (%h)" | grep -E "^$type" || true)
            if [ -n "$commits" ]; then
                echo "$commits" | sed -E "s/^$type(\([^)]*\))?:\s*/- /" >> CHANGELOG.md
            fi
            echo "" >> CHANGELOG.md
        }
        
        # Process each commit type
        process_commit_type "feat" "Features"
        process_commit_type "fix" "Bug Fixes"
        process_commit_type "docs" "Documentation"
        process_commit_type "refactor" "Code Refactoring"
        process_commit_type "test" "Tests"
        process_commit_type "perf" "Performance Improvements"
        process_commit_type "build" "Build System"
        process_commit_type "ci" "Continuous Integration"
        process_commit_type "chore" "Chores"
        
        # Add special handling for breaking changes
        BREAKING_CHANGES=$(git log --pretty=format:"%B" | grep -A 5 "BREAKING CHANGE" | grep -v "^$" | grep -v "^--$" || true)
        if [ -n "$BREAKING_CHANGES" ]; then
            echo "### BREAKING CHANGES" >> CHANGELOG.md
            echo "$BREAKING_CHANGES" >> CHANGELOG.md
            echo "" >> CHANGELOG.md
        fi
    fi
fi

echo -e "${COLOR_GREEN}Changelog generated at CHANGELOG.md${COLOR_RESET}"