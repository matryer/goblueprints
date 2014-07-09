#!/bin/bash
echo "Rebuilding everything..."
./all.sh
echo "Type in your starting word..."
echo -n "  > "
./synonyms | ./sprinkle | ./coolify | ./domainify | ./available
