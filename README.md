# 📊 Log Analyzer

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

A powerful CLI tool for parsing and analyzing log files with pattern detection, real-time monitoring, and comprehensive reporting. Built for developers, DevOps engineers, and system administrators who need quick insights into application logs.

## ✨ Features

- **🔍 Pattern Detection** - Built-in patterns for common issues (HTTP errors, memory leaks, auth failures)
- **📈 Log Level Analysis** - Automatic categorization and statistics by log level
- **⚡ Real-time Monitoring** - Watch files for changes with live analysis
- **📊 Comprehensive Reports** - JSON export with detailed metrics
- **🎯 Smart Alerts** - Configurable thresholds for critical patterns
- **🔧 Custom Patterns** - Define your own regex patterns via config files
- **⏰ Time Range Analysis** - Automatic timestamp extraction and duration calculation
- **🚨 Top Error Tracking** - Identify most frequent error messages

## 🚀 Installation

### Binary Release
```bash
# Download latest release
wget https://github.com/intruderfr/log-analyzer/releases/latest/download/log-analyzer-linux-amd64
chmod +x log-analyzer-linux-amd64
sudo mv log-analyzer-linux-amd64 /usr/local/bin/log-analyzer
```

### From Source
```bash
git clone https://github.com/intruderfr/log-analyzer.git
cd log-analyzer
go build -o log-analyzer main.go
```

### Go Install
```bash
go install github.com/intruderfr/log-analyzer@latest
```

## 📖 Usage

### Basic Analysis
```bash
# Analyze a log file
log-analyzer -file=/var/log/app.log

# Verbose output with progress
log-analyzer -file=/var/log/app.log -verbose

# Export results to JSON
log-analyzer -file=/var/log/app.log -output=report.json
```

### Real-time Monitoring
```bash
# Watch file for changes
log-analyzer -file=/var/log/app.log -watch

# Set custom alert threshold
log-analyzer -file=/var/log/app.log -watch -threshold=50
```

### Custom Patterns
```bash
# Use custom pattern config
log-analyzer -file=/var/log/app.log -config=patterns.json
```

## 🔧 Configuration

### Custom Patterns File (`patterns.json`)
```json
[
  {
    "name": "PAYMENT_ERROR",
    "pattern": "(?i)(payment (failed|declined|error)|transaction (failed|timeout))",
    "description": "Payment processing failures",
    "severity": "critical"
  },
  {
    "name": "SLOW_QUERY",
    "pattern": "slow query: [0-9]+\\.[0-9]+s",
    "description": "Database slow queries",
    "severity": "medium"
  }
]
```

### Built-in Patterns

| Pattern | Description | Severity |
|---------|-------------|----------|
| `HTTP_ERROR` | HTTP 4xx/5xx responses | High |
| `MEMORY_ERROR` | Out of memory, heap issues | Critical |
| `CONNECTION_ERROR` | Network connectivity problems | Medium |
| `AUTH_FAILURE` | Authentication failures | High |
| `DATABASE_ERROR` | Database connection/query errors | High |

## 📊 Output Examples

### Console Output
```
📊 Log Analysis Results
═══════════════════════

📈 Summary: Analyzed 15,420 log lines spanning 2h15m30s. Error rate: 2.3% (354 errors). 3 critical patterns detected

⏰ Time Range:
   Start: 2026-03-20 09:15:30
   End:   2026-03-20 11:31:00
   Duration: 2h15m30s

📋 Logs by Level:
   INFO: 12,890 (83.6%)
   ERROR: 354 (2.3%)
   WARN: 2,176 (14.1%)

🔍 Pattern Matches:
   ✅ HTTP_ERROR: 12 matches
   ⚠️  MEMORY_ERROR: 45 matches
   ✅ AUTH_FAILURE: 8 matches

🚨 Top Errors:
   1. Database connection timeout after 30s (count: 23)
   2. Failed to parse JSON response from API (count: 18)
   3. Memory allocation failed: cannot allocate... (count: 15)
```

### JSON Report Structure
```json
{
  "total_lines": 15420,
  "logs_by_level": {
    "INFO": 12890,
    "ERROR": 354,
    "WARN": 2176
  },
  "pattern_matches": {
    "HTTP_ERROR": 12,
    "MEMORY_ERROR": 45,
    "AUTH_FAILURE": 8
  },
  "top_errors": [
    "Database connection timeout after 30s (count: 23)",
    "Failed to parse JSON response from API (count: 18)"
  ],
  "time_range": {
    "start": "2026-03-20 09:15:30",
    "end": "2026-03-20 11:31:00",
    "duration": "2h15m30s"
  },
  "summary": "Analyzed 15,420 log lines spanning 2h15m30s..."
}
```

## 🎯 Use Cases

- **Production Monitoring** - Real-time log analysis with alerting
- **Incident Response** - Quickly identify error patterns during outages
- **Performance Analysis** - Track error rates and patterns over time
- **Compliance Auditing** - Generate reports for security and audit teams
- **CI/CD Integration** - Automated log analysis in deployment pipelines

## 🛠️ Advanced Features

### Pipeline Integration
```bash
# Use with other tools
tail -f /var/log/app.log | log-analyzer -file=/dev/stdin -watch

# Process rotated logs
find /var/log -name "app*.log" -exec log-analyzer -file={} \;
```

### Docker Usage
```bash
# Analyze container logs
docker logs myapp 2>&1 | log-analyzer -file=/dev/stdin

# Monitor in real-time
docker logs -f myapp 2>&1 | log-analyzer -file=/dev/stdin -watch
```

## 🔄 Supported Log Formats

- **Standard formats**: Apache, Nginx, syslog
- **Application logs**: JSON structured logs, plain text
- **Timestamps**: ISO 8601, RFC 3339, custom formats
- **Log levels**: DEBUG, INFO, WARN, ERROR, FATAL, TRACE

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Author

**Aslam Ahamed**
- GitHub: [@intruderfr](https://github.com/intruderfr)
- LinkedIn: [aslam-ahamed](https://linkedin.com/in/aslam-ahamed)

## ⭐ Star History

If this project helped you, please consider giving it a star! ⭐