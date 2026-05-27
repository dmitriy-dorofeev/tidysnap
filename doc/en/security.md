# Security

## Safety Mechanisms

### 1. Dry Run by Default

On first launch, `dry_run` is set to `true`. The utility shows which files will be deleted but **does not delete** them. This allows you to verify the settings before actual deletion.

### 2. Extension-Only Search

Scanning is performed strictly by file extension (`.png`, `.mov`, etc.), **not by name**. This prevents false positives on files with similar names (for example, `report.png.txt` will not be found as an image).

### 3. Root System Folders Restriction

The following cannot be selected as target directories:
- `/`
- `/System`
- `/Users`
- `~` (home folder)

### 4. Important Folder Warning

If one of the following folders is selected, a warning appears:
- `Desktop`
- `Downloads`
- `Documents`
- `Movies`
- `Music`
- `Pictures`
- `Public`
- `Library`

The user must explicitly acknowledge the risk (`warning_ack`).

### 5. Logging of All Operations

Every deletion (or dry-run) is logged with a timestamp, path, and file size. This allows you to track what was deleted.

## Recommendations

- **Always start with Dry Run** — make sure scanning finds only the intended files.
- **Check logs** after each cleanup (`[l]` in TUI or `~/Library/Logs/tidysnap/cleanup.log`).
- **Do not disable Dry Run** until you fully understand how the utility works.
- **Create backups** of important data before the first real cleanup.
