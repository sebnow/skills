---
type: tool_used
tool: Write
# Source paths, and the names of inkwell's classes and functions, in any
# written file. The pages describe parts, not the code that implements them.
input_match: 'inkwell/[\w/]+[.]py|\b(?:KeyStore|RateLimiter|JobQueue|ConverterRegistry|BlobStore|UsageMeter|DocumentRepo|WebhookOutbox|WebhookSender|RendererClient|PdfConverter|ThumbnailConverter|backoff_seconds|build_registry|build_app)\b'
min: 0
max: 0
---
