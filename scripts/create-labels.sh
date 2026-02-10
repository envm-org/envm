#!/bin/bash

# Define colors
COLOR_BACKEND="0E8A16" # Green
COLOR_FRONTEND="1D76DB" # Blue
COLOR_CLI="5319E7" # Purple
COLOR_DB="BFDC5C" # Yellow-Green
COLOR_DOCS="0075CA" # Blue
COLOR_DX="C5DEF5" # Light Blue
COLOR_SEC="B60205" # Red
COLOR_PRIORITY_HIGH="B60205" # Red
COLOR_PRIORITY_MED="FBCA04" # Yellow
COLOR_PRIORITY_LOW="0E8A16" # Green
COLOR_INFRA="795DA3" # Purple
COLOR_TESTING="FBCA04" # Yellow

# Check if we are in a git repo
if [ ! -d ".git" ]; then
    if [ -d "envm/.git" ]; then
        echo "Found git repo in envm directory. Switching to it..."
        cd envm
    else
        echo "Error: Not a git repository and 'envm' git repo not found."
        exit 1
    fi
fi

declare -A labels=(
    ["backend"]="$COLOR_BACKEND"
    ["ci-cd"]="0052CC"
    ["cli"]="$COLOR_CLI"
    ["database"]="$COLOR_DB"
    ["design"]="FBCA04"
    ["documentation"]="$COLOR_DOCS"
    ["dx"]="$COLOR_DX"
    ["events"]="E99695"
    ["frontend"]="$COLOR_FRONTEND"
    ["grpc"]="0052CC"
    ["gui"]="006B75"
    ["ide"]="C2E0C6"
    ["infra"]="$COLOR_INFRA"
    ["npm"]="CB3837"
    ["observability"]="D4C5F9"
    ["operations"]="006B75"
    ["performance"]="D93F0B"
    ["post-mvp"]="F9D0C4"
    ["priority-high"]="$COLOR_PRIORITY_HIGH"
    ["priority-medium"]="$COLOR_PRIORITY_MED"
    ["priority-low"]="$COLOR_PRIORITY_LOW"
    ["rest"]="1D76DB"
    ["security"]="$COLOR_SEC"
    ["setup"]="$COLOR_DX"
    ["testing"]="$COLOR_TESTING"
)

echo "Creating/Updating labels..."

for label in "${!labels[@]}"; do
    color="${labels[$label]}"
    echo "Processing label: $label (Color: $color)"
    # --force updates the label if it already exists
    gh label create "$label" --color "$color" --description "" --force
done

echo "Done."