#!/bin/bash
echo Building synonyms...
cd ../synonyms
go build -o ../build/synonyms

echo Building tokenize...
cd ../tokenize
go build -o ../build/tokenize

echo Building available...
cd ../available
go build -o ../build/available
cd ../build

echo Done.
