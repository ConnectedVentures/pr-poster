# pr-poster

Post content to a Github PR easily from CI

## How to Use

You can either pipe content into `pr-poster`, or if you have content in a file, just use the file as an argument:

```yaml
  build:
    docker:
      - image: connectedventures/pr-poster:v1.0
    steps:
      - run: echo "Test" | pr-poster
```

```yaml
  build:
    docker:
      - image: connectedventures/pr-poster:v1.0
    steps:
      - run: echo "Test" > output.txt
      - run: pr-poster output.txt
```

## Required Environment Variables

### Project/Context Variables

* `GITHUB_API_USERNAME` - Username for the authenticated user posting to Github
* `GITHUB_API_TOKEN` - Token to authenticate the user posting to Github

### CircleCI Generated

* `CIRCLE_PULL_REQUEST`
* `CIRCLE_PROJECT_USERNAME`
* `CIRCLE_PROJECT_REPONAME`
