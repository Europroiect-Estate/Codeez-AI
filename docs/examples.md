# Examples

## First run

```bash
codeez init
codeez provider set ollama
codeez provider set-model ollama llama3.2
codeez doctor
codeez config set palette cyber
```

## Chat (streaming)

```bash
codeez chat --no-tui
# Then type a message and press Enter.
```

Or with a one-liner:

```bash
echo "Explain recursion in Go" | codeez chat --no-tui
```

## Agentic task

```bash
codeez run "add a new cobra command called hello that prints Hello World"
```

This will create a session, stream a plan and suggested changes, and print a summary. Apply changes manually or via `codeez apply` (with approval).

## Git workflow

```bash
codeez git status
codeez git add internal/cli/hello.go
codeez git commit -m "feat: add hello command"
codeez git log --oneline --max 5
codeez git push   # prompts for network approval
```

## Index and repo map

```bash
codeez index
# Then repo.map() in the agent uses .codeez/repo_map.json
```

## Palettes

- **original**: Neutral modern (default)
- **corporate**: Clean, subdued, enterprise
- **cyber**: High-contrast, neon

```bash
codeez config set palette cyber
codeez config show
```
