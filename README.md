# github-forgejo-backup

## Tokens

Generate the Forgejo token with the following command:

```
curl -H "Content-Type: application/json" \
  -d '{"name":"mirrors","scopes":["read:user","write:organization","write:issue","write:repository"]}' \
  -u 'sos-admin:PASSWORD' https://forgejo.cloudogu.com/api/v1/users/sos-admin/tokens
```

---

Forgejo mascot by [David Revoy](https://www.peppercarrot.com/en/viewer/misc-src__2022-11-27_Forgejo_by-David-Revoy.html) (CC-BY 4.0)
