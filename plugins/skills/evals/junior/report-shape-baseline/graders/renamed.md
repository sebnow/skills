---
type: regex
target: {source: file, path: internal/export/beancount.go}
pattern: 'holdingAccount\(tx\)[\s\S]*func holdingAccount\(tx Transaction\) string'
match: contains
---
