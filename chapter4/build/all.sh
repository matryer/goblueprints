#!/bin/bash
echo Building synonyms...
cd ../synonyms
go build -o ../build/synonyms

echo Building available...
cd ../available
go build -o ../build/available
cd ../build

echo Building sprinkle...
cd ../sprinkle
go build -o ../build/sprinkle
cd ../build

echo Building coolify...
cd ../coolify
go build -o ../build/coolify
cd ../build

echo Building domainify...
cd ../domainify
go build -o ../build/domainify
cd ../build

echo Done.
