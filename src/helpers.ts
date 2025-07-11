function iconClassesToIcon(ic: string): string {
  let namespace = 'fas'
  let icon = ''

  for (const c of ic.split(' ')) {
    if (c === 'fa-fw') {
      continue
    }

    if (['fab', 'fas'].includes(c)) {
      namespace = c
    }

    if (c.startsWith('fa-')) {
      icon = c
    }
  }

  return [namespace, 'fa-fw', icon].join(' ')
}

export { iconClassesToIcon }
