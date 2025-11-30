# AWS Config Switcher

A simple Go tool to switch between different AWS configurations stored in separate folders. This tool uses fuzzy finding to select a configuration profile and updates your `~/.aws/config` and `~/.aws/credentials` files accordingly.

## Prerequisites

- Go installed on your system.
- AWS CLI installed and configured (for verification).
- Fuzzy finder library: `github.com/ktr0731/go-fuzzyfinder`

## Installation

1. Clone or download the repository.
2. Navigate to the `aws-config-switcher` directory.
3. Run `go mod tidy` to install dependencies.
4. Build the binary: `go build -o aws-config-switcher aws-config-switcher.go`

## Usage

1. **Prepare your AWS configurations:**
   - **Create a directory named `.awsconfigs` in your home directory.** (Note: The tool uses `.awsconfigs` as the constant directory name, defined as `const awsDir = ".awsconfigs"` in the code.)
   - Inside `~/.awsconfigs`, create subdirectories for each AWS profile (e.g., `personal-profile-ro`, `personal-profile-admin`).
   - In each profile directory, place `config` and `credentials` files as you would in `~/.aws/`.

   Example structure:
   ```
   ~/.awsconfigs/
   ├── personal-profile-ro/
   │   ├── config
   │   └── credentials
   ├── personal-profile-admin/
   │   ├── config
   │   └── credentials
   ├── work-profile-ro/
   │   ├── config
   │   └── credentials
   └── work-profile-admin/
       ├── config
       └── credentials
   ```

2. **Run the tool:**
   - Execute `./aws-config-switcher` (or the built binary name).
   - Use the fuzzy finder interface to select the desired profile.
   - The tool will:
     - Backup your current `~/.aws/config` and `~/.aws/credentials` to `*.backup` files.
     - Copy the selected profile's `config` to `~/.aws/config`.
     - Rewrite the selected profile's `credentials` under the `[default]` section in `~/.aws/credentials`.
     - Display the current AWS caller identity to confirm the switch.

3. **Verify the switch:**
   - The tool outputs the switched profile name and AWS identity details.
   - The ARN is highlighted in orange for easy identification.

## Example Output

```
AWS configuration switched successfully!
Switched AWS configuration to profile: personal-profile-admin
####################
Current AWS Identity:
{
    "UserId": "AIDACKCEVSQ6C2EXAMPLE",
    "Account": "123456789012",
    "Arn": "arn:aws:iam::123456789012:user/example-user"
}
```

## Notes

- The tool assumes AWS configurations are stored in `~/.awsconfigs` with subdirectories containing `config` and `credentials` files.
- Backups are created as `config.backup` and `credentials.backup` in `~/.aws/`.
- If the `~/.awsconfigs` directory does not exist, the tool will create it.
- Ensure your AWS credentials are valid and have the necessary permissions to call `aws sts get-caller-identity`.
