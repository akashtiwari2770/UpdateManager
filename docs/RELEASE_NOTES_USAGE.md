# Release Notes Usage Guide

This guide explains how to use the Release Notes feature in the Update Manager application.

## Overview

Release Notes allow you to document what's new, bug fixes, breaking changes, known issues, and upgrade instructions for each version of your product.

## Accessing Release Notes

1. Navigate to **Versions** from the main menu
2. Click on a version to view its details
3. Click on the **"Release Notes"** tab

## Viewing Release Notes

The Release Notes Viewer displays all release notes information in a formatted, easy-to-read layout:

- **Version Information**: Version number, release type, release date, and EOL date
- **What's New**: List of new features and improvements
- **Bug Fixes**: List of bugs that were fixed
- **Breaking Changes**: Changes that may break existing functionality
- **Known Issues**: Issues that users should be aware of
- **Upgrade Instructions**: Step-by-step instructions for upgrading

### Actions Available in Viewer

- **Share**: Share the release notes link (copies to clipboard or uses native share)
- **Print**: Print-friendly view of the release notes
- **Edit**: Edit release notes (only available for draft versions)

## Editing Release Notes

### Prerequisites

- The version must be in **Draft** state
- You must have edit permissions

### Steps to Edit

1. Navigate to the version details page
2. Click on the **"Release Notes"** tab
3. Click the **"Edit"** button (top right)
4. The Release Notes Editor modal will open

## Release Notes Editor Sections

### 1. What's New

Add new features and improvements:

- Click **"+ Add Item"** to add a new feature
- Enter the feature description in the text field
- Click **"Remove"** to delete an item
- Empty items are automatically filtered out when saving

**Example:**
- "Added support for dark mode"
- "Improved performance by 30%"
- "New dashboard with real-time analytics"

### 2. Bug Fixes

Document bugs that were fixed:

- Click **"+ Add Bug Fix"** to add a new bug fix entry
- **Fix ID**: Optional identifier (e.g., "BUG-123")
- **Issue Number**: Optional issue tracker reference (e.g., "#456")
- **Description**: Required description of the bug fix
- Click **"Remove"** to delete a bug fix entry

**Example:**
- Fix ID: `BUG-123`
- Issue Number: `#456`
- Description: "Fixed memory leak in data processing module"

### 3. Breaking Changes

Document changes that may break existing functionality:

- Click **"+ Add Breaking Change"** to add a new entry
- **Description**: Required description of the breaking change
- **Migration Steps**: Steps users need to take to migrate
- **Configuration Changes**: Any configuration changes required
- Click **"Remove"** to delete an entry

**Example:**
- Description: "API endpoint `/v1/users` has been removed"
- Migration Steps: "Use `/v2/users` endpoint instead. Update your API client to version 2.0"
- Configuration Changes: "Update `API_VERSION` in config file to `v2`"

### 4. Known Issues

Document issues that users should be aware of:

- Click **"+ Add Known Issue"** to add a new entry
- **Issue ID**: Optional identifier (e.g., "ISSUE-789")
- **Description**: Required description of the issue
- **Workaround**: Optional workaround if available
- **Planned Fix**: Optional information about planned fix or timeline
- Click **"Remove"** to delete an entry

**Example:**
- Issue ID: `ISSUE-789`
- Description: "Performance degradation when processing large files (>10GB)"
- Workaround: "Split large files into smaller chunks before processing"
- Planned Fix: "Will be addressed in version 2.1.0 (Q2 2025)"

### 5. Upgrade Instructions

Provide step-by-step upgrade instructions:

- Enter instructions in the text area
- Supports multi-line text
- You can use markdown formatting (basic support)
- This is a single text field (not a list)

**Example:**
```
Before upgrading:
1. Backup your current installation
2. Review breaking changes section
3. Update any custom configurations

Upgrade steps:
1. Download the new version package
2. Stop the service
3. Run the installer
4. Start the service
5. Verify installation

Post-upgrade:
1. Check logs for any errors
2. Verify all features are working
3. Update any dependent services
```

## Saving Release Notes

1. Fill in the sections you want to include
2. Click **"Save Release Notes"** at the bottom
3. The release notes will be saved to the version
4. The modal will close and the viewer will update

**Note:** Empty items and sections are automatically filtered out when saving.

## Best Practices

### What's New
- Be specific and clear
- Focus on user-visible changes
- Group related features together
- Use action-oriented language

### Bug Fixes
- Include issue tracker references when available
- Describe the impact of the bug
- Mention affected versions if applicable

### Breaking Changes
- Clearly explain what changed
- Provide detailed migration steps
- Include code examples if helpful
- Link to migration guides if available

### Known Issues
- Be transparent about issues
- Provide workarounds when possible
- Set expectations with planned fixes
- Prioritize critical issues

### Upgrade Instructions
- Be thorough but concise
- Include prerequisites
- List all steps in order
- Mention rollback procedures if applicable
- Include verification steps

## Workflow Example

1. **Create Version**: Create a new version in draft state
2. **Add Release Notes**: Click "Edit" in the Release Notes tab
3. **Fill Sections**: Add what's new, bug fixes, etc.
4. **Save**: Click "Save Release Notes"
5. **Review**: Review the formatted release notes in the viewer
6. **Submit for Review**: Once satisfied, submit the version for review

## Limitations

- Release notes can only be edited when the version is in **Draft** state
- Once submitted for review, release notes become read-only
- Empty items are automatically removed when saving
- Markdown support is basic (full markdown rendering may be added in future)

## Troubleshooting

### Can't Edit Release Notes
- **Check version state**: Only draft versions can be edited
- **Check permissions**: Ensure you have edit permissions

### Changes Not Saving
- **Check for errors**: Look for error messages in the modal
- **Check network**: Ensure you have a stable connection
- **Check validation**: Ensure required fields are filled

### Formatting Issues
- **Text formatting**: Basic markdown is supported
- **Line breaks**: Use line breaks in text areas for better formatting
- **Special characters**: Most special characters are supported

## API Integration

Release notes are saved via the version update API:

```typescript
PUT /api/v1/versions/:id
{
  "release_notes": {
    "whats_new": [...],
    "bug_fixes": [...],
    "breaking_changes": [...],
    "known_issues": [...],
    "upgrade_instructions": "..."
  }
}
```

## Related Features

- **Version Management**: Release notes are part of version management
- **Package Management**: Packages are uploaded separately
- **Compatibility**: Compatibility information is in a separate tab

