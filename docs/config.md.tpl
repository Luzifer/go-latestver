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

check_interval: 1h

...
```

Each catalog entry contains a `name` and a `tag` representing the entry. You can choose those freely but they should be URL-safe. Some examples I'm using are: `alpine:stable`, `google-chrome:dev`, `google-chrome:stable`, `factorio:latest`, …

Additionally you will configure a `fetcher` with its corresponding `fetcher_config` for the catalog entry. In the example above the `html` fetcher is used with two attributes configured. The attributes for each fetcher can be found below.

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
