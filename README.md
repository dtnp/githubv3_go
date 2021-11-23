# Github Repo Parser

Defining what a repo does, who owns it, and if it is still active can be difficult.  As an organization grows
it becomes difficult to keep track of all work done, ensure correct documentation, and avoid effort duplication.
This project is a simple parser that could help automate efforts to solve those issues.

# Current Functionality

1. Get all repos
2. Get last commit message
3. Calculate days since their last commit

## Github Token
To get PRIVATE repository information, you must supply a github token

## How to Run

Pull down the required modules/libraries
```sh
go mod tidy
```

Run the script and pass in your github token
```sh
GITHUB_TOKEN=xxx go run main.go
```

## Last Commits
If the repo is not archived, the last commit information will be fetched.  This was done to try and figure out
ownership of repos that have not been worked on for a long time

## Future Ideas

1. CODEOWNERS
    - [Official Docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
    - Check for file (CODEOWNERS, OWNERS) in the following standard directories:
        + `/` (root)
        + `.docs/`
        + `.github/`
    - If file exists, validate and parse
        + Set is "COExists" flag in the CSV output to: true or false
        + Validate to make sure it is correctly formatted
        + Set is "COValid" flag in the CSV output to: true or false
        + Put owners into csv (could be many, needs some thought)
2. README.md format
   - First line should be the "clean title" for the repo
   - A business summary should be next explaining the goal and point of this repo
   - A "How to Run" section must exist to explain how to run the repo locally
