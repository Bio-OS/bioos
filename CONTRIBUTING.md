<!-- omit in toc -->

# Contributing Guidelines

First off, thanks for taking the time to contribute!

Many types of contributions are encouraged (see [Ways to Contribute](#ways-to-contribute)). Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions.

<!-- omit in toc -->

## Table of Contents

- [Contributing Guidelines](#contributing-guidelines)
  - [Code of Conduct](#code-of-conduct)
  - [Ways to Contribute](#ways-to-contribute)
  - [How to Contribute](#how-to-contribute)
    - [Report Bugs](#report-bugs)
      - [Before Submitting a Bug Report](#before-submitting-a-bug-report)
      - [How Do I Submit a Good Bug Report?](#how-do-i-submit-a-good-bug-report)
    - [Suggest Enhancements](#suggest-enhancements)
      - [Before Submitting an Enhancement](#before-submitting-an-enhancement)
      - [How Do I Submit a Good Enhancement Suggestion?](#how-do-i-submit-a-good-enhancement-suggestion)
    - [Your First Code Contribution](#your-first-code-contribution)
    - [Improve Documentation](#improve-documentation)
  - [Styleguides](#styleguides)
    - [Code Formatting](#code-formatting)
    - [Commit Messages](#commit-messages)

## Code of Conduct

This project and everyone participating in it is governed by the
[Code of Conduct](https://github.com/Bio-OS/bioos/blob/main/CODE_OF_CONDUCT.md).

## Ways to Contribute

We welcome many different types of contributions including:

-   New features
-   Builds, CI/CD
-   Bug fixes
-   Documentation
-   Web design
-   ...

## How to Contribute

> ### Legal Notice <!-- omit in toc -->
>
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

### Report Bugs

<!-- omit in toc -->

#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

-   Make sure that you are using the latest version.
-   Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [documentation](https://bio-os.gitbook.io/userguide/).).
-   To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/Bio-OS/bioos/issues?q=label%3Abug).
-   Also make sure to search the internet (including Stack Overflow) to see if users outside of the GitHub community have discussed the issue.
-   Collect information about the bug:
    -   Stack trace (Traceback)
    -   OS, Platform and Version (Windows, Linux, macOS, x86, ARM)
    -   Version of the interpreter, compiler, SDK, runtime environment, package manager, depending on what seems relevant.
    -   Possibly your input and the output
    -   Can you reliably reproduce the issue? And can you also reproduce it with older versions?

<!-- omit in toc -->

#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to **[email needed to be filled from core maintainers]**.

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

-   Open an [Issue](https://github.com/Bio-OS/bioos/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
-   Explain the behavior you would expect and the actual behavior.
-   Please provide as much context as possible and describe the _reproduction steps_ that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
-   Provide the information you collected in the previous section.

Once it's filed:

-   The project team will label the issue accordingly.
-   A team member will try to reproduce the issue with your provided steps. If there are no reproduction steps or no obvious way to reproduce the issue, the team will ask you for those steps and mark the issue as `needs-repro`. Bugs with the `needs-repro` tag will not be addressed until they are reproduced.
-   If the team is able to reproduce the issue, it will be marked `needs-fix`, as well as possibly other tags (such as `critical`), and the issue will be left to be [implemented by someone](#your-first-code-contribution).


### Suggest Enhancements

This section guides you through submitting an enhancement suggestion for Bio-OS, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->

#### Before Submitting an Enhancement

-   Make sure that you are using the latest version.
-   Read the [documentation](https://bio-os.gitbook.io/userguide/) carefully and find out if the functionality is already covered, maybe by an individual configuration.
-   Perform a [search](https://github.com/Bio-OS/bioos/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
-   Find out whether your idea fits with the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset. If you're just targeting a minority of users, consider writing an add-on/plugin library.

<!-- omit in toc -->

#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/Bio-OS/bioos/issues).

-   Use a **clear and descriptive title** for the issue to identify the suggestion.
-   Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
-   **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
-   You may want to **include screenshots and animated GIFs** which help you demonstrate the steps or point out the part which the suggestion is related to.
-   **Explain why this enhancement would be useful** to most Bio-OS users. You may also want to point out the other projects that solved it better and which could serve as inspiration.

<!-- You might want to create an issue template for enhancement suggestions that can be used as a guide and that defines the structure of the information to be included. If you do so, reference it here in the description. -->

### Your First Code Contribution

Here's a quick rundown for contributing to an open source project:

1. Find an issue that you are interested in addressing or a feature that you would like to add.
2. Fork the repository associated with the issue to your local GitHub organization. This means that you will have a copy of the repository under `Bio-OS/bioos`.
3. Clone the repository to your local machine using `git clone https://github.com/Bio-OS/bioos.git`.
4. Create a new branch for your fix using `git checkout -b fix-xxx`.
5. Make the appropriate changes for the issue you are trying to address or the feature that you want to add.
6. Use `git add <file>` to add the file contents of the changed files to the "snapshot" git uses to manage the state of the project, also known as the index.
7. Use `git commit -s -m "Some description"` to store the contents of the index with a descriptive message. Here "Some description" should follow our [commit message convention](#commit-messages).
8. Push the changes to the remote repository using `git push origin fix-xxx`.
9.  Submit a pull request to the upstream repository.
10. Title the pull request with a short description of the changes made and the issue or bug number associated with your change. For example, you can title an issue like so "Added more log outputting to resolve #4352".
11. In the description of the pull request, explain the changes that you made, any issues you think exist with the pull request you made, and any questions you have for the maintainer. It's OK if your pull request is not perfect (no pull request is), the reviewer will be able to help you fix any problems and improve it!
12. Wait for the pull request to be reviewed by a maintainer.
13. Make changes to the pull request if the reviewing maintainer recommends them.
14. Celebrate your success after your pull request is merged!

### Improve Documentation

TODO: guidelines for updating, improving and correcting the docs

## Styleguides

### Code Formatting

You can simply utilize the basic Go formatter for go snippets, prettier and its plugins for web frontend. Variables are named in camelCase, full word in most scenarios.

### Commit Messages

Commit messages should follow the [AngularJS convention](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.greljkmo14y0).


<!-- omit in toc -->
## Attribution
This guide is based on the [contributing-gen](https://github.com/bttger/contributing-gen).