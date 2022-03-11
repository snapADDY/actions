# detect-changes

Action to detect changes between the current and last job with a custom key. Can be used to keep track of incremental releases to many environments.

## Usage

```yaml
    # [...]
    steps:
      - name: Change Detection
          uses: snapADDY/actions/detect-changes@v1
          id: detect-changes
          with:
            persist-run: ${{ steps.img-tag.outputs.publish }}

      - name: Incremental Build
        run: make build
        env:
          INCR_LAST_COMMIT: ${{ steps.detect-changes.outputs.last-commit-sha }}

```

## Input
| Name | Type | Default | Description |
| - | - | - | - |
| `change-key` | string | `${{ github.event_name }}-${{ github.ref }}` | Key prefix to differentiate multiple target environments  |
| `persist-run` | string(boolean) | `"true"` | Should this run be persisted on success |


## Output
| Name | Type | Example | Description |
| - | - | - | - |
| `last-commit-sha` | string | `1f4e5a738c8a63cc07efa5abd89e3f3f8289b59e` | Commit SHA of the last persisted run |
| `changed-files` | string(multiline) | `"detect-changes/dist/post-run/index.js\ndetect-changes/dist/run/index.js\ndetect-changes/src/action.ts\n"` | List of changes files between the restored and current commit |
