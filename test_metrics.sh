#!/bin/bash
echo "Testing /metrics endpoint..."
curl -s http://localhost:8080/metrics | jq '.'
