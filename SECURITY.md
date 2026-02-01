# Security

## Reporting a vulnerability

Please report security issues privately. Do not open a public GitHub issue for security-sensitive bugs.

- Email or otherwise contact the maintainers to report a vulnerability.
- Include a clear description and steps to reproduce if possible.
- We will acknowledge and work on a fix; we may request more detail.

## Security design

- **No telemetry**: Codeez does not send usage data or analytics.
- **Local-first**: Ollama runs entirely on your machine; API keys are stored only in your config.
- **Approvals**: Tool actions (file system, git, commands) require explicit user approval.
- **Secret redaction**: API keys and detected secrets are redacted in logs and audit output.
