# Changelog

This project provides a tool for making the handling of changelogs easier for software developers.
Having a changelog is important for users of your software. The conveys critical changes and new features
to the user.
It must not be confused with commit logs which are mostly technical in nature and cannot be
changed later without rewriting the history of the repository.
A simple text file does not work well for changelogs. When multiple developers prepare changes
in feature branches, the single changelog file is usually the source of permanent merge conflicts.
To solve this issue, this project provides a tool which makes it very easy.

## Install

Install the changelog util into your development environment: 
```shell
$ !!! TO DO !!!
```

## How to use
Initialize the changelog for your development project. In the project root execute:
```shell
$ changelog init
```
This will create a `.changelog.yaml` configuration file and a subdirectory called `changelog`.
Commit them both to source control.

When you want to add a new changelog entry for your software, execute on your feature branch:
```shell
$ changelog add
```
The changelog util will ask you for a description of your feature and save it in the `changelog` subdirectory.
Each unreleased changelog is saved in its own file with a timestamp in its filename. That way, there are no
merge conflicts when multiple developers make changes in lots of feature branches. The changelog can easily
be cherry picked into other branches, because it is just a file.

When you want to update the changelog for your software, execute on your main branch:
```shell
$ changelog update
```
This will update a `CHANGELOG.md` file in your project root.
Commit the changes to your version control.

When you want to create a new release of your software, execute on your main branch:
```shell
$ changelog release
```
The changelog util will ask you for the version specification you want to use for the release. It will also 
consolidate all unreleased changelog entries into a single release and clean up the `changelog` subdirectory.
Commit the changes to your version control.

## Customization
If the naming conventions of the changelog utils do not suit your needs, you can always edit the `.changelog.yaml`
config file and adjust its behavior. Note that you have to rename files and directories manually when you make
changes to the config file.

You can even customize the template which is used to render the `CHANGELOG.md`. It is a simple Go template.
