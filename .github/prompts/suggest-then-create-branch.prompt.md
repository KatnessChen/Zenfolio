# GitHub Copilot Agent Mode Branch Name Prompt

## Goal:

Instruct GitHub Copilot Agent Mode to generate a concise, descriptive Git branch name based on task/issue content, and automatically provide the `git checkout -b` command for a streamlined workflow.

## Desired Branch Name Pattern:

Branch names should follow this structure for easy identification and automation:

## Desired Output:

The final output should be the generated branch name, followed by the `git checkout -b` command using that name, formatted as follows:

Suggested Branch Name: Command to create branch: git checkout -b

## Guidelines:

To ensure consistent and useful branch names, adhere to these guidelines:

- **`<type>`**: Mandatory prefix indicating task type. Choose one:
  - `feat`: New feature.
  - `fix`: Bug fix.
  - `docs`: Documentation changes.
  - `chore`: Maintenance, tool/dependency updates.
  - `refactor`: Code restructuring.
  - `perf`: Performance optimization.
- **`<issue-id>`**: Mandatory task/issue number (e.g., `JIRA-123`, `GH-456`). Prioritize actual IDs; use `core-001` if none.
- **`<short-description>`**: Mandatory brief summary of task content, converted to lowercase, hyphen-separated words (e.g., `implement-user-registration`). Be concise.
- **No Special Characters**: Only lowercase letters, numbers, and hyphens (`-`). No spaces, underscores, or other special characters.
- **Length**: Reasonable length (max 50 chars) for readability.
- **Clarity**: Must clearly represent the task/issue without viewing full content.
- **Command Generation**: The `git checkout -b` command _must_ be generated using the exact branch name produced.

## Example Generated Output:

Here is an example of the desired output based on the above guidelines:

Suggested Branch Name: feat/GH-101/add-user-profile-editingCommand to create branch: git checkout -b feat/GH-101/add-user-profile-editing

## Execution:

Execute the `git checkout -b` in the terminal to create the branch with the suggested name. Ensure the branch name is valid and follows the specified pattern before executing.
