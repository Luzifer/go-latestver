#!/usr/bin/env python3

import jinja2
import sys
import re


def main(args):
    modules = []

    for filename in args:
        with open(filename, 'r') as codefile:
            mod = {
                'attributes': [],
            }

            for line in codefile:
                match = re.search(r'@module (.*)', line)
                if match is not None:
                    mod['type'] = match[1]

                match = re.search(r'@module_desc (.*)', line)
                if match is not None:
                    mod['description'] = match[1]

                match = re.search(
                    r'@attr ([^\s]+) ([^\s]+) ([^\s]+) "([^"]*)" (.*)', line)
                if match is not None:
                    mod['attributes'].append({
                        'name': match[1],
                        'required': match[2],
                        'type': match[3],
                        'default': match[4],
                        'description': match[5],
                    })

        mod['attributes'] = sorted(
            mod['attributes'], key=lambda a: ('0' if a['required'] == 'required' else '1') + ':' + a['name'])

        modules.append(mod)

    modules = sorted(modules, key=lambda m: m['type'])

    with open('docs/config.md.tpl', 'r') as f:
        tpl = jinja2.Template(f.read())
        print(tpl.render(modules=modules))

    return 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
