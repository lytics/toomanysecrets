toomanysecrets
==============

Deletes all private gists.

First go to https://github.com/settings/tokens and generate a token with just
`Gist` access.

Then if you have Go installed:

```
go get github.com/lytics/toomanysecrets

# Doesn't delete anything by default!
GH_TOKEN=... toomanysecrets

# When you're sure you're ready:
GH_TOKEN=... toomanysecrets -dryrun=false
```
