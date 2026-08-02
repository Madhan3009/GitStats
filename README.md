# GitStats
GitStats is a Go-based CLI that scans local folders for .git repositories, stores them in ~/.gitlocalstats, and analyzes commit history from the last 183 days to visualize contributions. It uses the go-git library for both plumbing and porcelain operations and works fully offline without GitHub APIs.

## Testing
Run `./gogit -email madhankumaar96@gmail.com` to see your local git commit statistics in a contribution calendar format.
