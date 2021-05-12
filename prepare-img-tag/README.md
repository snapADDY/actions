# prepare-img-tag
Prepares a tag for Docker images and decide if they should be published

Tag names have multiple formats. For images that will never be published (e.g Pull Requests) the format `run-${RunNumber}` is used. Git tags containing only valid characters for docker tags are used as is. Images from push events which could be published use the `${branch|overwrite}-${shortSHA}-${RunNumber}` format.

## Ruleset

### Push - Tag
Git tags using only allowed characters will always be published using the tag as tag-name. For all other tags see [Other Triggers](other-triggers)

### Push - Branch
Push actions depend on the branch and git commit message to determine the output.
Pushes to `master`, `main`, or `staging` will always be published.
For all other branches the following rules apply:
- Will not be published unless commit message contains a `[img]` flag
- Uses the branch name of the current branch replacing all `/` with `-` and dropping all other illegal characters
- Branch name part of the image tag can be overwritten with a commit message flag `[img:mytestenv]`

Commit message flags(e.g. `[img]`) are case insensitive and can be anywhere in the commit message. Examples for valid flags:
- `[img]` - Shortest valid flag;
- `[img:testenv]` - Override the branch name with "testenv"
- `[img::dev]` - Dev build enabled (not used by gh-actions jobs yet)
- `[img:testenv:dev]`

### Other Triggers
All other events don't trigger an image publish.


## Usage

```yaml
    # [...]
    steps:
      - name: Checkout snapADDY/actions
        uses: actions/checkout@v2
        with:
          repository: snapADDY/actions
          path: .github/actions
          persist-credentials: false
          ssh-key: ${{ secrets.ACTIONS_REPO_SSH_KEY }}
      - name: Prepare Image Tag
        uses: ./.github/actions/prepare-img-tag
        id: img-tag

      # [...]

      - name: Publish Images
        if: success() && steps.img-tag.outputs.publish == 'true'
        run: make docker-push
```
## Output
| Name | Type | Example | Description |
| - | - | - | - |
| `tag-name` | string | `fix_123-1a2b3c4-26` | Docker image tag |
| `publish` | string(boolean) | `"true"` | Decides if images should be published |

