#!/bin/bash

# Get the committer email of the latest commit
committer_email=$(git log -1 --pretty=format:'%ae')

# Define the allowed email domain
allowed_domain="@company.com"

# Check if the committer's email ends with the allowed domain
if [[ "$committer_email" == *"$allowed_domain" ]]; then
    echo "Committer email '$committer_email' is valid."
    exit 0
else
    echo "Error: Committer email '$committer_email' does not end with '$allowed_domain'."
    echo "Please configure your Git email to use a company email address."
    exit 1
fi
