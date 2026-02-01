# Installation

## Quick install (curl)

```bash
curl -sSL https://raw.githubusercontent.com/Europroiect-Estate/Codeez-AI/main/scripts/install.sh | bash
```

Or with explicit version:

```bash
VERSION=v0.1.0 bash -c "$(curl -sSL https://raw.githubusercontent.com/Europroiect-Estate/Codeez-AI/main/scripts/install.sh)"
```

## Homebrew

Tap and install:

```bash
brew tap Europroiect-Estate/Codeez-AI
brew install codeez
```

Or install from the repo formula:

```bash
brew install Europroiect-Estate/Codeez-AI/codeez
```

## npm (wrapper)

Installs the appropriate binary for your platform from GitHub releases:

```bash
npm install -g codeez
```

Then run:

```bash
codeez version
```

## Build from source

Requires Go 1.23+.

```bash
git clone https://github.com/Europroiect-Estate/Codeez-AI.git
cd Codeez-AI
go build -o codeez ./cmd/codeez
./codeez version
```

For a versioned binary:

```bash
go build -ldflags "-X github.com/Europroiect-Estate/Codeez-AI/internal/cli.Version=v0.1.0" -o codeez ./cmd/codeez
```
