# fork

`fork` forks a repository and set a remote for it.

## Usage

```
fork # in a git repository
```

## Github

### Environment

You need to set these environment variables to create fork or access private repo.

```
FORK_GITHUB_USER # your github id
FORK_GITHUB_AUTH # your github authentication token
```

### Create an Authentication token

Login into your account.

Goto `Setting` > `Developer Settings` > `Personal access token`

Select `Generate new token`. Turn on the `repo` checkbox.
When you generate the token, you should see your token value.

Save it in your `KEEP_GITHUB_AUTH` environment variable.

