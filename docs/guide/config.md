# Config

Flamigo provides a flexible **configuration system** that helps you manage application settings across environments and components. It supports reading config files from any `fs.FS` source, aggregating them, and enriching them with environment variables — all while keeping configuration structured, modular, and testable.

::: warning
The config module can only be enabled when initializing a new project.
:::

## Features

- Load configuration from embedded files or external sources  
- Support for environment-based overrides using `APP_ENV`  
- Built-in integration with Flamigo’s Dependency Injection system  
- Simple file structure with extensibility in mind

---

## Loading Configurations

You can load configuration files in two ways:

### 1. Loading a Directory

This reads and merges multiple config files from a folder:

```go
LoadDirectory(fileSys fs.FS, dirPath string) (map[string]*Config, error)
```

Useful for managing multiple configs (e.g., one per domain or feature) in a centralized way.

### 2. Loading a Single File

To load a single configuration file:

```go
LoadConfigFile(fileSys fs.FS, path string) (*Config, error)
```

This is ideal for smaller setups or when working with isolated config values.

---

## Default Setup

When the config module is added to your project, the configuration structure lives under:

```bash
/internal/config
```

The default behavior is:

- Configuration files are stored as **embedded files** inside:  
  ```
  /internal/config/configs
  ```
- The system automatically selects the config file based on the current `APP_ENV` environment variable.
- The resulting `*config.Config` instance is registered with Flamigo's **Dependency Manager**, making it easily injectable wherever needed.

---

## Example Structure

```
internal/
└── config/
    ├── configs/
    │   ├── development.yaml
    │   ├── staging.yaml
    │   └── production.yaml
    └── loader.go
```

::: tip
Use `APP_ENV=development` to control which config file gets loaded at runtime.
:::
## Using Injected Config

Once loaded, the `*config.Config` struct is automatically available for injection:

```go
func InitService(cfg *config.Config) error {
  dbURL := cfg.Get("database.url")
  // ...
}
```

This makes it easy to access environment-specific values without scattering logic throughout your codebase.