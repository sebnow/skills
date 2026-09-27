# 1. Record architecture decisions

Date: 2025-06-02

## Status

Accepted

## Context

The platform team makes decisions about cluster layout, naming and
networking that every service team depends on. Those decisions are currently
spread across Confluence pages that drift out of date.

## Decision

We record architecture decisions in `docs/adr/` of platform-infra, using the
format described by Michael Nygard.

## Consequences

Service teams have one place to check platform conventions. Decisions are
reviewed as pull requests to this repository.
