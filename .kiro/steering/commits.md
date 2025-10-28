# Commit Guidelines

## Developer Certificate of Origin (DCO) Signoff

All commits to this repository MUST include a Developer Certificate of Origin (DCO) signoff. This is a requirement for contributing to the Valkey project and ensures legal compliance.

### Required Signoff Format

Every commit message must end with a signoff line in the following format:

```
Signed-off-by: Your Name <your.email@example.com>
```

### How to Add Signoffs

#### Automatic Signoff (Recommended)

Always use the `--signoff` or `-s` flag when committing:

```bash
git commit -s -m "Your commit message"
```

#### Configure Git to Always Signoff

Set up Git to automatically add signoffs to all commits:

```bash
git config --global format.signOff true
```

#### Amend Missing Signoffs

If you forget to add a signoff, amend the last commit:

```bash
git commit --amend --signoff --no-edit
```

For multiple commits without signoffs, use interactive rebase:

```bash
git rebase -i HEAD~n --signoff
```

### What the Signoff Means

By adding the signoff, you certify that:

1. The contribution was created in whole or in part by you and you have the right to submit it under the open source license indicated in the file
2. The contribution is based upon previous work that, to the best of your knowledge, is covered under an appropriate open source license and you have the right under that license to submit that work with modifications
3. The contribution was provided directly to you by some other person who certified (1), (2) or (3) and you have not modified it
4. You understand and agree that this project and the contribution are public and that a record of the contribution is maintained indefinitely

### Enforcement

- All pull requests will be checked for proper signoffs
- Commits without signoffs will be rejected
- This applies to all contributors, including maintainers

### Example Commit Message

```
feat: Add ReadFrom parsing support to ConfigurationOptions

- Implement ParseReadFromStrategy method
- Add validation for ReadFrom strategy and AZ parameter combinations
- Extend DoParse method to handle readFrom and az parameters

Addresses GitHub issue #26

Signed-off-by: Joe Brinkman <joe.brinkman@improving.com>
```

## Commit Message Format

Follow conventional commit format for consistency:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]

Signed-off-by: Your Name <your.email@example.com>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(config): add ReadFrom parsing support
fix(client): resolve connection timeout issue
docs(readme): update installation instructions
test(config): add comprehensive ReadFrom tests
```
