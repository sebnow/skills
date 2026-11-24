# Getting Started

This tutorial sends one document to Inkwell and downloads the PDF.

## Get an API Key

Ask the Inkwell team for an API key for your tenant.

## Send a Document

Send the file as a multipart form to `POST /v1/documents`, with your key in
the `Authorization` header:

```sh
curl -H "Authorization: Bearer $INKWELL_KEY" \
  -F file=@report.docx -F output=pdf \
  https://inkwell.example.com/v1/documents
```

The response holds the document's `id` and its `status`.

## Download the Result

Ask for the document until its status is `done`, then follow the `pdf` link
in the response.
