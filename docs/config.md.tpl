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

check_interval: 1h

...
```

Each catalog entry contains a `name` and a `tag` representing the entry. You can choose those freely but they should be URL-safe. Some examples I'm using are: `alpine:stable`, `google-chrome:dev`, `google-chrome:stable`, `factorio:latest`, …

Additionally you will configure a `fetcher` with its corresponding `fetcher_config` for the catalog entry. In the example above the `html` fetcher is used with two attributes configured. The attributes for each fetcher can be found below.

You can provide your own `links` for each catalog entry which will be added to or override the links returned from the fetcher. If you provide the same `name` as the fetcher uses the link of the fetcher will be overridden. The `icon_class` should consist of `fas` or `fab` and an icon (for example `fa-globe`). You can use all **solid** (`fas`) or **breand** (`fab`) icons within [Font Awesome v5 Free](https://fontawesome.com/v5.15/icons?d=gallery&s=brands,solid&m=free).

## Available Fetchers

{% for module in modules -%}
## Fetcher: `{{ module.type }}`

{{ module.description }}

| Attribute | Req. | Type | Default Value | Description |
| --------- | :--: | ---- | ------------- | ----------- |
{%- for attr in module.attributes %}
| `{{ attr.name }}` | {% if attr.required == 'required' %}✅{% endif %} | {{ attr.type }} | {% if attr.default != "" %}`{{ attr.default }}`{% endif %} | {{ attr.description }} |
{%- endfor %}

{% endfor %}

<!-- vim: set ft=markdown : -->
