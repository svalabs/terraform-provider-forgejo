#!/bin/sh

# Change the contents of this output to get the environment variables
# of interest. The output must be valid JSON, with strings for both
# keys and values.
cat <<EOF
{
  "FORGEJO_ADMIN_USERNAME": "$FORGEJO_ADMIN_USERNAME",
  "FORGEJO_ADMIN_PASSWORD": "$FORGEJO_ADMIN_PASSWORD"
}
EOF
