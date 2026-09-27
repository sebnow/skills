# 2. One namespace per team

Date: 2025-06-20

## Status

Accepted

## Context

Services from different teams share the `apps` namespace. Resource quotas,
RBAC and network policies cannot be scoped to a team, and a noisy workload
from one team can starve another.

## Decision

Each team gets its own namespace, named after the team's short code (`pay`,
`cat`, `ship`, `id`). Quotas, RoleBindings and default-deny NetworkPolicies
are applied per namespace by the `team-namespace` Helm chart.

## Consequences

Teams own their namespace's quota. Cross-team traffic needs an explicit
NetworkPolicy, which platform reviews.
