#!/bin/bash

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

INPUT_FILE="../docs/mvp2.txt"
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: $INPUT_FILE not found"
    exit 1
fi

title=""
desc=""
labels=""

echo "Reading from $INPUT_FILE..."

while IFS= read -r line || [ -n "$line" ]; do
    # Trim whitespace
    line=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    
    if [[ "$line" == "---" ]]; then
        if [[ -n "$title" ]]; then
            echo "Creating issue: $title"
            
            cmd=(gh issue create --title "$title" --body "$desc")
            
            IFS=',' read -ra LABEL_ARRAY <<< "$labels"
            for i in "${LABEL_ARRAY[@]}"; do
                # trim whitespace
                lbl=$(echo "$i" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
                if [[ -n "$lbl" ]]; then
                    cmd+=(--label "$lbl")
                fi
            done
            
            # Execute command
            "${cmd[@]}"
            
            if [ $? -eq 0 ]; then
                echo "Success"
            else
                echo "Failed to create issue: $title"
            fi
            
            # Reset variables
            title=""
            desc=""
            labels=""
            
            # Sleep to avoid rate limiting
            sleep 2
        fi
        continue
    fi
    
    if [[ "$line" == title:* ]]; then
        title="${line#title: }"
    elif [[ "$line" == desc:* ]]; then
        desc="${line#desc: }"
    elif [[ "$line" == label:* ]]; then
        labels="${line#label: }"
    fi
done < "$INPUT_FILE"

# Handle the last entry if file doesn't end with ---
if [[ -n "$title" ]]; then
    echo "Creating issue: $title"
    cmd=(gh issue create --title "$title" --body "$desc")
    IFS=',' read -ra LABEL_ARRAY <<< "$labels"
    for i in "${LABEL_ARRAY[@]}"; do
        lbl=$(echo "$i" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [[ -n "$lbl" ]]; then
            cmd+=(--label "$lbl")
        fi
    done
    "${cmd[@]}"
fi
