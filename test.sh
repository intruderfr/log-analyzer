#!/bin/bash

echo "🧪 Testing Log Analyzer"
echo "======================="

# Build the application
echo "Building application..."
go build -o log-analyzer main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Test basic analysis
echo ""
echo "📊 Testing basic log analysis..."
./log-analyzer -file=examples/sample.log -verbose

# Test with custom patterns
echo ""
echo "🎯 Testing with custom patterns..."
./log-analyzer -file=examples/sample.log -config=examples/custom-patterns.json

# Test JSON output
echo ""
echo "💾 Testing JSON output..."
./log-analyzer -file=examples/sample.log -output=test-report.json
echo "Report saved to test-report.json"

echo ""
echo "✅ All tests completed successfully!"