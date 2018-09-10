# pr-poster

Post content to a Github PR easily from CI

## Required Environment Variables

### Project/Context Variables

* `GITHUB_API_USERNAME` - Username for the authenticated user posting to Github
* `GITHUB_API_TOKEN` - Token to authenticate the user posting to Github

### CircleCI Generated

* `CIRCLE_PULL_REQUEST`
* `CIRCLE_PROJECT_USERNAME`
* `CIRCLE_PROJECT_REPONAME`
