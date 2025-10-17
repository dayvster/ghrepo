<div align="center">

# ghprofile 

A beautiful, customizable GitHub profile viewer for your terminal.

---

<p>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <img src="https://img.shields.io/badge/go-1.21+-00ADD8.svg" alt="Go Version">
  <img src="https://img.shields.io/github/stars/dayvster/ghrepo?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/github/downloads/dayvster/ghrepo/total.svg" alt="Downloads">
  <img src="https://img.shields.io/github/issues/dayvster/ghrepo.svg" alt="Issues">
  <img src="https://img.shields.io/github/forks/dayvster/ghrepo.svg" alt="Forks">
  <img src="https://img.shields.io/github/repo-size/dayvster/ghrepo.svg" alt="Repo size">
  <img src="https://img.shields.io/github/last-commit/dayvster/ghrepo.svg" alt="Last commit">
</p>

</div>

---

## Features
- ✨ Fetch and display GitHub user profiles and top repositories
- 📊 Shows language stats, repo stars, forks, and more
- 🎨 Icon-rich output (with options for plain text)
- ⚡ Caching and demo mode for offline/limited API use
- 🛠️ CLI flags for customization

---

## Usage
```sh
./ghprofile [flags]
```

### Flags
- `-u`, `--user`        GitHub username to fetch (default: dayvster)
- `-n`                  How many top repos to show (default: 5)
- `--no-icons`          Disable icons in the output
- `--no-border`         Remove card border from output
- `--no-style`          Remove all styles from output
- `--no-demo`           Do not fall back to demo data on fetch error; exit instead
- `--demo`              Force demo data (skip network and cache)
- `-h`, `--help`        Show help message

---

## Screenshots
| Default | No Border | No Icons |
|---------|-----------|----------|
| ![ghprofile demo](screenshots/basic.png) | ![No border](screenshots/no-border.png) | ![No icons](screenshots/no-icons.png) |

---

## Build
```sh
go build -o ghprofile ./cmd/main.go
```

## Example
```sh
./ghprofile --user dayvster --topN 5
```

---

## License
See [LICENSE](LICENSE).

## Contributing
See [CONTRIBUTING.md](CONTRIBUTING.md).
