# github-forgejo-backup

## Tokens

Generate the Forgejo token with the following command:

```
curl -H "Content-Type: application/json" \
  -d '{"name":"mirrors","scopes":["read:user","write:organization","write:issue","write:repository"]}' \
  -u 'sos-admin:PASSWORD' https://forgejo.cloudogu.com/api/v1/users/sos-admin/tokens
```
