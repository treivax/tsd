# TSD Documentation

Welcome to the TSD (Type System Development) documentation. This guide will help you get started and master all features of the TSD rule engine.

## Documentation Overview

### 🚀 Getting Started

Start here if you're new to TSD:

1. **[Installation Guide](INSTALLATION.md)** - Install TSD on your system
2. **[Quick Start](QUICK_START.md)** - Get up and running in 5 minutes
3. **[Tutorial](TUTORIAL.md)** - Step-by-step learning path

### 📚 Core Documentation

Complete reference documentation:

- **[User Guide](USER_GUIDE.md)** - Comprehensive guide covering all features:
  - Language syntax and structure
  - Type system and pattern matching
  - Conditions and operators
  - Actions and rule execution
  - Type casting
  - String operations
  - Arithmetic operations
  - Configuration and best practices

- **[Grammar Guide](GRAMMAR_GUIDE.md)** - Complete language syntax reference
  - Lexical structure
  - Type definitions
  - Rule syntax
  - Expressions and operators
  - Comments and identifiers

- **[API Reference](API_REFERENCE.md)** - HTTP API documentation
  - Server endpoints
  - Request/response formats
  - Authentication
  - Error handling

### 🔐 Security & Operations

- **[Authentication Guide](AUTHENTICATION.md)** - Complete authentication documentation
  - API key authentication
  - JWT tokens
  - TLS/SSL configuration
  - Key management and rotation

- **[Logging Guide](LOGGING_GUIDE.md)** - Logging configuration and reference
  - Log levels and formats
  - Output destinations
  - Structured logging
  - Performance considerations

### 🏗️ Advanced Topics

- **[Architecture](ARCHITECTURE.md)** - Technical architecture and design
  - RETE algorithm implementation
  - Node types and network structure
  - Memory management and optimization
  - Performance characteristics
  - Transaction system
  - Concurrency model

- **[Contributing](CONTRIBUTING.md)** - Contribution guidelines
  - Development setup and workflow
  - Coding standards and style guide
  - Testing requirements and best practices
  - Pull request process
  - Performance guidelines

## Quick Links

### By Use Case

**I want to...**

- **Start using TSD quickly** → [Quick Start](QUICK_START.md)
- **Learn TSD step by step** → [Tutorial](TUTORIAL.md)
- **Understand all features** → [User Guide](USER_GUIDE.md)
- **Look up syntax** → [Grammar Guide](GRAMMAR_GUIDE.md)
- **Use the HTTP API** → [API Reference](API_REFERENCE.md)
- **Secure my deployment** → [Authentication](AUTHENTICATION.md)
- **Deploy to production** → [Installation](INSTALLATION.md) + [Authentication](AUTHENTICATION.md)
- **Debug issues** → [User Guide - Troubleshooting](USER_GUIDE.md#troubleshooting)

### By Topic

**Core Language Features:**
- [Types](USER_GUIDE.md#type-system)
- [Facts](USER_GUIDE.md#pattern-matching)
- [Rules](USER_GUIDE.md#pattern-matching)
- [Actions](USER_GUIDE.md#actions)
- [Conditions](USER_GUIDE.md#conditions)
- [Operators](USER_GUIDE.md#conditions)

**Advanced Features:**
- [Type Casting](USER_GUIDE.md#type-casting)
- [String Operations](USER_GUIDE.md#string-operations)
- [Arithmetic](USER_GUIDE.md#arithmetic-operations)
- [Pattern Matching](USER_GUIDE.md#string-operations) (LIKE, MATCHES, CONTAINS, IN)

**Deployment:**
- [Installation](INSTALLATION.md)
- [Configuration](USER_GUIDE.md#configuration)
- [Server Mode](USER_GUIDE.md#server-mode)
- [Authentication](AUTHENTICATION.md)

## Examples

See the [examples/](../examples/) directory for complete, working examples:

```bash
# List all examples
ls ../examples/

# Run an example
tsd ../examples/basic-rules.tsd
tsd ../examples/type-casting.tsd
tsd ../examples/string-operations.tsd
```

## Getting Help

### Documentation Issues

If you find errors or gaps in the documentation:
- [Report an issue](https://github.com/treivax/tsd/issues)
- Suggest improvements
- Submit documentation pull requests

### Usage Questions

For help using TSD:
1. Check the [User Guide](USER_GUIDE.md)
2. Review [examples/](../examples/)
3. Enable debug logging: `TSD_LOG_LEVEL=debug tsd program.tsd`
4. Ask in GitHub Discussions

### Bug Reports

To report bugs:
1. Check [existing issues](https://github.com/treivax/tsd/issues)
2. Create a new issue with:
   - TSD version (`tsd --version`)
   - Operating system
   - Minimal reproduction case
   - Expected vs actual behavior

## Document Status

| Document | Status | Last Updated |
|----------|--------|--------------|
| Installation | ✅ Complete | 2024-12-07 |
| Quick Start | ✅ Complete | 2024-12-07 |
| User Guide | ✅ Complete | 2024-12-07 |
| Tutorial | ✅ Complete | - |
| Grammar Guide | ✅ Complete | - |
| API Reference | ✅ Complete | - |
| Authentication | ✅ Complete | - |
| Logging Guide | ✅ Complete | - |
| Architecture | ✅ Complete | 2024-12-07 |
| Contributing | ✅ Complete | 2024-12-07 |

## Contributing to Documentation

We welcome documentation improvements! To contribute:

1. Fork the repository
2. Edit or create markdown files in `docs/`
3. Test examples and code snippets
4. Submit a pull request

**Documentation standards:**
- Clear, concise language
- Complete, working examples
- Proper markdown formatting
- Cross-references between documents
- Code examples with expected output

## License

This documentation is part of the TSD project and is licensed under the MIT License.
See [LICENSE](../LICENSE) for details.