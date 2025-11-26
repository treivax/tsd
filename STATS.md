# 📊 TSD Code Statistics

**⚡ Quick Access to Latest Reports**

This file provides quick links to all code statistics and quality reports for the TSD project.

---

## 🎯 Main Reports

### 📊 Interactive Dashboard
**[docs/reports/DASHBOARD.md](docs/reports/DASHBOARD.md)** - Visual dashboard with ASCII charts
- Real-time metrics overview
- Coverage by package with visual bars
- Priority matrix
- Quick commands

### 📈 Detailed Statistics Report
**[docs/reports/CODE_STATS_2025-11-26.md](docs/reports/CODE_STATS_2025-11-26.md)** - Complete analysis
- Volume metrics (29,434 total lines, 11,614 manual, 12,590 tests)
- Coverage by package (48.7% global)
- Complexity analysis
- Recommendations and next actions

### 🌐 Coverage HTML Report
**[docs/reports/coverage_report.html](docs/reports/coverage_report.html)** - Interactive browser view
- Line-by-line coverage visualization
- Generated with `go tool cover`
- Color-coded coverage indicators

### 📋 Machine-Readable Metrics
**[docs/reports/code_metrics.json](docs/reports/code_metrics.json)** - JSON format
- Structured data for CI/CD
- Automated reporting
- Grafana/dashboard integration

---

## 📊 Current Statistics (2025-11-26)

```
Commit: 68fcd48
Branch: main

Code Volume:
  Total Lines:        29,434
  Manual Code:        11,614 (39.5%)
  Tests:              12,590 (42.8%)
  Generated:           5,230 (17.8%)

Files:
  Total:                  90
  Production:             59
  Tests:                  31

Quality Metrics:
  Test/Code Ratio:     108.4% ✅
  Global Coverage:      48.7% 🟡
  Quality Score:        85/100 ✅
```

---

## 🚀 Quick Commands

```bash
# Update all statistics
./update_stats.sh

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Generate JSON metrics
./generate_metrics.sh

# View dashboard
cat docs/reports/DASHBOARD.md
```

---

## 📈 Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| `rete/pkg/domain` | 100.0% | ✅ Perfect |
| `rete/pkg/network` | 100.0% | ✅ Perfect |
| `constraint/pkg/validator` | 96.5% | ✅ Excellent |
| `constraint/pkg/domain` | 90.0% | ✅ Excellent |
| `rete/pkg/nodes` | 71.6% | 🟢 Good |
| `constraint` | 59.6% | 🟡 Fair |
| `rete` | 39.7% | 🟠 Needs Work |
| `test/integration` | 29.4% | 🟠 Needs Work |
| `cmd/*` | 0.0% | 🔴 Critical |

---

## 🎯 Priorities

### 🔴 High Priority (Immediate)
- [ ] Test `cmd/tsd` (0% → 80%) - 2-3h
- [ ] Test `cmd/universal-rete-runner` (0% → 70%) - 2-3h

### 🟡 Medium Priority (This Sprint)
- [ ] Increase `rete` coverage (39.7% → 70%) - 4-6h
- [ ] Increase `constraint` coverage (59.6% → 75%) - 3-4h
- [ ] Complete `rete/pkg/nodes` (71.6% → 90%) - 2-3h

### 🟢 Low Priority (Next Sprint)
- [ ] Test config packages (0% → 80%) - 2-4h
- [ ] Increase integration tests (29.4% → 60%) - 3-4h

**Goal:** Reach 70%+ global coverage (Est: 20-30 hours)

---

## 📚 Related Documentation

- [Testing Guide](docs/TESTING.md)
- [Development Guidelines](docs/development_guidelines.md)
- [Session Report](docs/SESSION_REPORT_2025-11-26.md)
- [Test Reports](docs/testing/)

---

## 🔄 Update Frequency

- **Statistics:** Run `./update_stats.sh` after major changes
- **Reports:** Updated automatically by the script
- **Manual Review:** Weekly or after adding new tests

---

**Last Updated:** 2025-11-26  
**Maintained By:** Engineering Team  
**Auto-Generated:** Yes (via update_stats.sh)