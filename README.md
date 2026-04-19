# vaultpipe

> CLI tool to inject secrets from Vault or env files into subprocess environments

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

## Usage

Inject secrets from HashiCorp Vault into a subprocess:

```bash
vaultpipe --vault-path secret/myapp -- ./myapp serve
```

Inject secrets from a `.env` file:

```bash
vaultpipe --env-file .env -- npm start
```

Combine both sources:

```bash
vaultpipe --vault-path secret/myapp --env-file .env.local -- python app.py
```

### Environment Variables

| Variable | Description |
|---|---|
| `VAULT_ADDR` | Address of the Vault server |
| `VAULT_TOKEN` | Authentication token for Vault |

### Flags

| Flag | Description |
|---|---|
| `--vault-path` | Vault secret path to read from |
| `--env-file` | Path to a `.env` file |
| `--no-inherit` | Do not inherit the parent process environment |

## How It Works

`vaultpipe` fetches secrets from the specified sources, merges them into an environment, and executes the given command as a subprocess with that environment applied. Secrets never touch disk and are scoped to the subprocess lifetime.

## License

MIT © [yourusername](https://github.com/yourusername)