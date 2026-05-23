# cronlens

> Parse and visualize crontab schedules with human-readable next-run times and conflict detection.

---

## Installation

```bash
go install github.com/yourusername/cronlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronlens.git
cd cronlens && go build -o cronlens .
```

---

## Usage

Parse a crontab expression and see the next scheduled runs:

```bash
$ cronlens "*/15 9-17 * * 1-5"

Expression : */15 9-17 * * 1-5
Description: Every 15 minutes, between 09:00 and 17:00, Monday through Friday

Next runs:
  1. Mon, 02 Jun 2025 09:00:00
  2. Mon, 02 Jun 2025 09:15:00
  3. Mon, 02 Jun 2025 09:30:00
  4. Mon, 02 Jun 2025 09:45:00
  5. Mon, 02 Jun 2025 10:00:00
```

Scan a full crontab file for conflicts:

```bash
$ cronlens --file /etc/crontab --conflicts

[CONFLICT] Jobs overlap at 2025-06-02 10:00:00:
  → "*/10 * * * *"  (backup-job)
  → "0,10,20 * * * *"  (sync-job)
```

### Flags

| Flag | Description |
|------|-------------|
| `--file` | Path to a crontab file to analyze |
| `--conflicts` | Detect overlapping job schedules |
| `--count` | Number of next-run times to display (default: 5) |
| `--tz` | Timezone for output (default: local) |

---

## License

[MIT](LICENSE) © 2025 yourusername