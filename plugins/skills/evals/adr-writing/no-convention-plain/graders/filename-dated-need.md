---
type: regex
target: last_message
pattern: '(?<![\w-])\d{4}-\d{2}-\d{2}-(?![a-z0-9-]*(?:grpc|protobuf))[a-z0-9-]+\.md\b'
match: contains
---
