## fly-github-auth-task

Reusable Concourse task that gets the ATC bearer token that results from logging in a GitHub user.

Useful for automating `fly` operations on a Concourse team that uses GitHub auth.

## Getting a Token

The following Concourse pipeline will get you a fresh ATC bearer token for a GitHub-authed team once every 23 hours:

```
---
jobs:
- name: get-github-token
  plan:
  - get: every-23h
    trigger: true
  - get : fly-github-auth-task
  - task: get-atc-token
    file: fly-github-auth-task/task.yml
    timeout: 10m
    attempts: 5
    params:
      GITHUB_USERNAME: {{ci_github_username}}
      GITHUB_PASSWORD: {{ci_github_password}}
      ATC_URL: {{atc_url}}
      TEAM_NAME: team-name
  - put: github-token
    params:
      file: bearer-token/bearer-token

resources:
- name: every-23h
  type: time
  source: {interval: 23h}

- name: github-token
  type: s3
  source:
    versioned_file: resources/bearer-token
    bucket: {{bucket_name}}
    region_name: eu-central-1
    access_key_id: {{concourse_ci_s3_access_key}}
    secret_access_key: {{concourse_ci_s3_secret_key}}
    server_side_encryption: AES256

- name: fly-github-auth-task
  type: git
  source:
    uri: https://github.com/EngineerBetter/fly-github-auth-task.git
    branch: master
```

## Using a Token

Write the token into a `./flyrc` file, and you'll be automagically auth'd.

Here's a BASH snippet to do just that:

```
# Load contents of file into env var
export ATC_BEARER_TOKEN=$(<bearer-token/bearer-token)

cat <<ENDOFSCRIPT > ~/.flyrc
targets:
  ${CONCOURSE_TARGET_NAME}:
    api: ${CONCOURSE_ATC_API}
    team: ${CONCOURSE_TEAM_NAME}
    insecure: false
    token:
      type: Bearer
      value: ${ATC_BEARER_TOKEN}
ENDOFSCRIPT
```

## How The Token Is Acquired

One can't get a user's OAuth token from GitHub using the API - it's a bit nonsensical to offer an API endpoint that gets a token for a human-based grant.

We use Agouti to drive a headless web browser to do the 'manual' login for a user. It was easiest to use Agouti from within Ginkgo, so we're being awful and using tests as an entrypoint to our app.