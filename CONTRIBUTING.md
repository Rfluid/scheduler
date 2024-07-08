--- Generated with Copilot ---

# Contributing to Scheduler Library

First off, thank you for considering contributing to Scheduler Library. It's people like you that make Scheduler Library such a great tool.

## Getting Started

- Submit a ticket for your issue, assuming one does not already exist.
  - Clearly describe the issue including steps to reproduce when it is a bug.
  - Make sure you fill in the earliest version that you know has the issue.
- Fork the repository on GitHub.

## Making Changes

- Create a topic branch from where you want to base your work.
  - This is usually the main branch.
  - Only target release branches if you are certain your fix must be on that branch.
  - To quickly create a topic branch based on main; `git branch fix/main/my_contribution main`. Then checkout the new branch with `git checkout fix/main/my_contribution`. Please avoid working directly on the `main` branch.
- Make commits of logical units.
- Check for unnecessary whitespace with `git diff --check` before committing.
- Make sure your commit messages are in the proper format.

## Submitting Changes

- Push your changes to a topic branch in your fork of the repository.
- Submit a pull request to the repository in the Rfluid organization.
- The core team looks at Pull Requests on a regular basis.

## Code Style

Follow the same coding style as the rest of the project. You can check the style by running `gofmt -d .` in your terminal.

## Examples

If you are developing functionality, please add examples to the [examples](examples) directory making sure that they are as simple as they can be to exemplify your feature.

## Documentation

We strive for complete documentation. If you make changes to the functionality of the project, please ensure that corresponding changes to the documentation are made.

## Thank you

Thanks again for your contribution, your time and effort are greatly appreciated. Happy coding!
