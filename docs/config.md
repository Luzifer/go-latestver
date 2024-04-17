# Configuration file format

```yaml
---

catalog:

  - name: alpine
    tag: stable

    fetcher: html
    fetcher_config:
      url: https://alpinelinux.org/downloads/
      xpath: '//div[@class="l-box"]/p/strong'

    links:
      - icon_class: 'fas fa-globe'
        name: 'Website'
        url: 'https://alpinelinux.org'

    version_constraint:
      allow_downgrade: false
      allow_prerelease: false
      type: semver

check_interval: 1h

...
```

Each catalog entry contains a `name` and a `tag` representing the entry. You can choose those freely but they should be URL-safe. Some examples I'm using are: `alpine:stable`, `google-chrome:dev`, `google-chrome:stable`, `factorio:latest`, …

Additionally you will configure a `fetcher` with its corresponding `fetcher_config` for the catalog entry. In the example above the `html` fetcher is used with two attributes configured. The attributes for each fetcher can be found below.

You can provide your own `links` for each catalog entry which will be added to or override the links returned from the fetcher. If you provide the same `name` as the fetcher uses the link of the fetcher will be overridden. The `icon_class` should consist of `fas` or `fab` and an icon (for example `fa-globe`). You can use all **solid** (`fas`) or **breand** (`fab`) icons within [Font Awesome v5 Free](https://fontawesome.com/v5.15/icons?d=gallery&s=brands,solid&m=free).

For the `version_constraint`s you must specify a `type` to parse the version returned by the fetcher (currently `semver` and `numeric_dot` (`104.0.5112.79`) are supported) and then can allow downgrades and pre-releases in the version. If no constraint is present, versions are neither parsed nor checked for downgrade / pre-release. If only the `type` is specified, downgrades and pre-releases are forbidden.

## Available Fetchers

## Fetcher: `atlassian`

Fetches latest version of an Atlassian product

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `product` | ✅ | string |  | Lowercase name of the product to fetch (e.g. confluence, crowd, jira-software, ...) |
| `edition` |  | string |  | Filter down the versions according to its edition (e.g. "Enterprise" or "Standard" for Confluence) |
| `search` |  | string | `TAR.GZ` | What to search in the download description: default is to search for the standalone .tar.gz file |

## Fetcher: `git_tag`

Reads git tags (annotated and leightweight) from a remote repository and returns the newest one

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `remote` | ✅ | string |  | Repository remote to fetch the tags from (should accept everything you can use in `git remote set-url` command) |

## Fetcher: `github_release`

Fetches the latest release from Github for a given repository not marked as pre-release

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `repository` | ✅ | string |  | Repository to fetch in form `owner/repo` |

## Fetcher: `helm`

Fetches the index file of a Helm Repo and yields the latest Helm-Chart version

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `chart` | ✅ | string |  | Chart to fetch the version of (i.e. "grafana") |
| `repo` | ✅ | string |  | URL of the repo (i.e. "https://grafana.github.io/helm-charts") |

## Fetcher: `html`

Downloads website, selects text-node using XPath and optionally applies custom regular expression

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `url` | ✅ | string |  | URL to fetch the HTML from |
| `xpath` | ✅ | string |  | XPath expression leading to the text-node containing the version |
| `regex` |  | string | `(v?(?:[0-9]+\.?){2,})` | Regular expression to apply to the text from the XPath expression |

## Fetcher: `json`

Fetches a JSON / JSONP file from remote source and traverses it using XPath expression

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `url` | ✅ | string |  | URL to fetch the HTML from |
| `xpath` | ✅ | string |  | XPath expression leading to the text-node containing the version |
| `jsonp` |  | boolean | `false` | File contains JSONP function, strip it to get the raw JSON |

## Fetcher: `regex`

Fetches URL and applies a regular expression to extract a version from it

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
| `regex` | ✅ | string |  | Regular expression (RE2) to apply to the text fetched from the URL. The regex MUST have exactly one submatch containing the version. |
| `url` | ✅ | string |  | URL to fetch the content from |



<!-- vim: set ft=markdown : -->
