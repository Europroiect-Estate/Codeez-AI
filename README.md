<p align="center">
  <strong>Codeez</strong>
</p>

<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=flat-square" alt="License"></a>
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-555?style=flat-square" alt="Platform">
  <a href="https://europroiect.org"><img src="https://img.shields.io/badge/Made%20by-Europroiect%20Estate-0d9488?style=flat-square" alt="Europroiect Estate"></a>
</p>

<p align="center">
  <sub>⚡ Agentic Coding CLI — rapid, sigur, multi-provider</sub>
</p>
<p align="center">
  <a href="https://europroiect.org">Europroiect Estate</a>
</p>

---

## Despre

**Codeez** este un CLI agentic pentru codare: rapid, curat în terminal, cu streaming, TUI modern și suport multi-provider (OpenAI, Anthropic, Ollama local). Rulează pe mașina ta, lucrează pe repository-urile reale și cere aprobări explicite pentru acțiuni sensibile (fișiere, git, comenzi).

Proiect open-source by **[Europroiect Estate](https://europroiect.org)**.

---

## Caracteristici

| Caracteristică | Descriere |
|----------------|------------|
| **Multi-provider** | Ollama (local), OpenAI, Anthropic — comutare ușoară, chei în config |
| **Sigur din start** | Sandbox, aprobări pe acțiuni, redactare secrete, fără telemetrie |
| **Terminal de calitate** | Palete (original, corporate, cyber), streaming, TUI opțional (Bubble Tea) |
| **Workflow agentic** | Plan → cod → review → git → test, cu aprobări explicite (o dată / sesiune / mereu pentru repo) |

---

## Instalare

### curl (recomandat)

```bash
curl -sSL https://raw.githubusercontent.com/Europroiect-Estate/Codeez-AI/main/scripts/install.sh | bash
```

### Homebrew

```bash
brew tap Europroiect-Estate/Codeez-AI
brew install codeez
```

### npm (wrapper)

```bash
npm install -g codeez
```

### Compilare din sursă

Necesită Go 1.23+.

```bash
git clone https://github.com/Europroiect-Estate/Codeez-AI.git
cd Codeez-AI
go build -o codeez ./cmd/codeez
./codeez version
```

#### Director de instalare

Scriptul respectă (în ordine): `CODEEZ_INSTALL_DIR` → `XDG_BIN_DIR` → `/usr/local/bin` → `~/.local/bin` → `~/bin`.

```bash
CODEEZ_INSTALL_DIR=/usr/local/bin curl -sSL .../install.sh | bash
XDG_BIN_DIR=$HOME/.local/bin curl -sSL .../install.sh | bash
```

**Sfat:** Completare în shell: `codeez completion bash` (sau `zsh` / `fish`) — adaugă în `.bashrc` / `.zshrc` conform instrucțiunilor.

Mai multe variante: [docs/install.md](docs/install.md).

---

## Început rapid

```bash
codeez init
codeez provider set ollama
codeez doctor
codeez config show
codeez config set palette corporate
codeez chat --no-tui   # chat streaming (fără --no-tui pornește TUI)
codeez run "adaugă o comandă hello"
```

---

## Comenzi

| Comandă | Descriere |
|---------|-----------|
| `codeez init` | Inițializează proiectul (`.codeez/`) |
| `codeez chat` | Chat interactiv TUI sau streaming |
| `codeez run "<sarcină>"` | Execuție agentică: plan, aprobări, rezumat |
| `codeez provider list` / `set` / `set-key` / `set-model` / `test` | Gestionează providerii LLM |
| `codeez doctor` | Verifică mediul (git, rg, ollama, node) |
| `codeez index` | Construiește harta repo (fișiere cheie, limbaje) |
| `codeez config show` / `set` | Afișează sau editează config (paletă, provider etc.) |
| `codeez git init` / `status` / `diff` / `add` / `commit` / `branch` / `checkout` / `log` / `remote` / `push` / `pull` | Operații git (push/pull cer aprobare) |
| `codeez version` | Afișează versiunea |

---

## Config și palete

- **Config global**: `~/.config/codeez/config.toml`
- **Config proiect**: `.codeez/config.toml` (suprascrie globalul)
- **Paletă**: `codeez config set palette original | corporate | cyber`

---

## Model de securitate

- **Sandbox**: Operațiile pe fișiere sunt limitate la repo (sau cwd); căi sensibile (~/.ssh, /etc etc.) sunt blocate.
- **Aprobări**: Acțiunile (fs, patch, git, cmd) cer aprobare explicită: Refuz / O dată / Sesiune / Mereu pentru acest repo.
- **Secrete**: Cheile API și secretele detectate sunt redactate în loguri și UI; commit-urile cu secrete detectate sunt blocate decât dacă confirmi explicit.
- **Fără telemetrie**: Codeez nu trimite date de utilizare sau analitică.

---

## Cum se compară cu OpenCode / Claude Code?

Inspirat de [OpenCode](https://github.com/anomalyco/opencode) și fluxurile tip Claude Code. Diferențe relevante pentru Codeez:

| Aspect | Codeez |
|--------|--------|
| **Stack** | Go 1.23+, Cobra, Bubble Tea, SQLite — un singur binar, fără runtime JS/Node |
| **Provider** | Ollama (local), OpenAI, Anthropic — comutare din config, fără vendor lock-in |
| **Aprobări** | O dată / Sesiune / Mereu pentru acest repo, persistate în `.codeez/permissions.toml` |
| **TUI** | Bubble Tea + Lip Gloss, palete (original, corporate, cyber) |
| **Telemetrie** | Zero — nu trimite date nicăieri |

Potrivit pentru: mediu local-first, repo-uri reale pe mașina ta, aprobări explicite și un singur executabil ușor de distribuit.

---

## Licență

MIT — vezi [LICENSE](LICENSE).

---

<p align="center">
  <a href="https://europroiect.org">Europroiect Estate</a> · <a href="CONTRIBUTING.md">Contribuie</a> · <a href="SECURITY.md">Securitate</a>
</p>
