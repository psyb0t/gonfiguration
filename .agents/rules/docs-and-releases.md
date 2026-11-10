# Docs and releases

Prefix: DOC.

- **DOC1:** `README.md` is the front door: what it is, install, a short example, and links. Reference material lives on the site in `docs/`, published by GitHub Pages from `master` at https://psyb0t.github.io/gonfiguration/. Each fact has one home.
- **DOC2:** Every page in `docs/` has front matter with a unique `title` and a `description` of 50 to 170 characters, and one `#` heading. Page-to-page links use the relative `.md` path. A link to a file outside `docs/` uses the full `https://github.com/psyb0t/gonfiguration/blob/master/...` URL. README links to the site use full URLs.
- **DOC3:** `make docs` builds the site into `.testing/docs-site` with GitHub's own build image. Build it after changing `docs/` and fix broken links before committing.
- **DOC4:** Code examples in the docs must compile against the current API. Error output shown in the docs is copied from a real run, not typed by hand.
- **DOC5:** A release is a commit on `master` that passed CI, a `CHANGELOG.md` section for the version, and an annotated `vX.Y.Z` tag on that commit whose message is the release note. The tag run of `pipeline.yml` waits for the `master` run of the same commit, then cuts the GitHub release. There is no version file.
- **DOC6:** Prose follows the README voice: direct, first person singular when a person speaks, no em dashes, "directory" and not "folder".
