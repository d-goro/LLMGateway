#!/bin/bash
echo "Testing /health endpoint..."
curl -s http://localhost:8080/health | jq '.'
