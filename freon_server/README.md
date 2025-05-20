# Freon

Freon is the server pendant of [frigoligo](https://github.com/casimir/frigoligo). It is a wallabag proxy server made to experiment on features too heavy to be done in the client but too experimental to be included in wallabag as-is. Like a lab.

## Why not contribute to wallabag directly?

Maintainers of wallabag are reluctant to include features that are not part of the core features, for [good reasons](https://wallabag.org/news/wallabag-wont-accept-pr-feature-request). While it makes a lot of sense for the project and the stability of the software, it also make feature improvements harder.

The codebase as a lot of history and while clean and well written, working with it requires a PHP track record which I don't have. What I have though is a lot of ideas on what features would be nice to have in frigoligo (and other wallabag clients).
The big part of these experiments is to evaluate and see if these features are worth the effort to implement. Once they are, I intend to contribute them to wallabag.

## Experiments (as of now)

### API authentication

The current authentication flow of wallabag is a partial OAuth2 implementation that is annoying to use and work with. It comes with OAuth2 complexity but without most of the advantages.
For example refresh tokens are not always working and only password grant is reliable making storing username and password locally a semi-requirement which defeat the usage of OAuth2 in the first place.
It makes implementing a client more complex than it should be but also prevent the usage of one-shot scripts/commands (e.g. to save content) without storing credentials locally or implementing a stateful flow.

The goal of this experiment is to try a new authentication system based on bearer tokens. The main goal is convenience and simplicity but security is not forgotten and is addressed by adding an optional lifetime and scopes to the tokens.

- [x] Bearer token authentication
- [x] API token scope

### Continuity

_Yep, the one [from Apple](https://www.apple.com/macos/continuity/)._

Conceptually, frigoligo is a client that you should find on any of your devices. This ubiquity is nice but frustrating in practice when switching from one device to another. That's the exact use case for the principle of continuity.

In order to make this work, frigoligo needs to have access to the state of things, with the server as the source of truth. The state does not need to be specific to frigoligo, on a more theoretical level, the server just need to maintain arbitrary data at different level (e.g. account, entry, ...)

- [ ] User preferences sync (for frigoligo: notification badge, reader theme, ...)
- [ ] Article (entry) read progression

### Service announcements

By its experimental nature, a freon service won't necessarily have strong QoS guarantees. Some experiments might also require users actions for them to work. This is why it is important to have a simple way to make announcements to the users of the service (e.g. wallabag credentials reset, planned maintenance, feature deprecation, ...)

- [ ] Service announcements endpoint

### API inconsistencies (in wallabag)

- [ ] More precise HTTP status codes instead of only 200 and _not-200_
- [ ] Remove all tags of an entry
- [ ] Real incremental sync (including deletions)

### Post-consume steps

This one is inspired by the awesome [paperless](https://github.com/the-paperless-project/paperless) project which is a kind of wallabag but for documents.
One of the killer features of paperless is the ability to automatically tag documents. This is done by using light machine learning based on good old text classifiers and implicit user feedback, as cheap on resources as it is effective.

- [ ] Post ingestion webhook
- [ ] Tags suggestion
- [ ] Summarization (still need to think if it would be useful?)
- [ ] Push notifications (ping clients to trigger a refresh)

### Miscellaneous

- [ ] Save an article using only a GET request (for a simpler intergration in all kind of tools)
- [ ] Fetch article from [archive.org](https://archive.org/help/wayback_api.php)
