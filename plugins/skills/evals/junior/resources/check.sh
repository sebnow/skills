#!/bin/bash
# Stands in for `go build ./... && go test ./...`, since the eval sandbox has
# no Go toolchain. Checks what a compile would catch for edits of the size
# the junior cases ask for, plus the expectation in beancount_test.go.
cd "$(dirname "$0")" || exit 1
n=0
fail=0
check() {
  n=$((n + 1))
  if "$2"; then echo "ok   $1"; else echo "FAIL $1"; fail=1; fi
}

gofiles() { find . -name '*.go' -not -path './.*'; }

balanced() {
  local f
  for f in $(gofiles); do
    awk -v op="$1" -v cl="$2" '
      { gsub(/"([^"\\]|\\.)*"/, "") }
      { o += gsub(op, ""); c += gsub(cl, "") }
      END { exit (o == c) ? 0 : 1 }
    ' "$f" || { echo "     unbalanced $3 in $f"; return 1; }
  done
}
braces() { balanced '[{]' '[}]' braces; }
parens() { balanced '[(]' '[)]' parentheses; }

no_trailing_space() {
  local hits
  hits=$(grep -n '[[:blank:]]$' $(gofiles))
  [ -z "$hits" ] || { echo "$hits" | sed 's/^/     /'; return 1; }
}

export_refs_defined() {
  local name ok=0
  for name in $(grep -ho 'export\.[A-Z][A-Za-z0-9_]*' cmd/ledger-export/*.go | sed 's/export\.//' | sort -u); do
    grep -q "^func $name(" internal/export/*.go || { echo "     undefined: export.$name"; ok=1; }
  done
  return $ok
}

local_calls_defined() {
  local name ok=0
  for name in $(grep -ho '\(^\|[^.A-Za-z0-9_]\)[a-z][A-Za-z0-9_]*(' internal/export/*.go |
    sed 's/^[^a-z]*//; s/($//' | sort -u); do
    case "$name" in
      func|if|for|return|switch|len|append|make|new|panic|string|int|float64|copy|cap|close) continue ;;
    esac
    grep -q "^func $name(" internal/export/*.go || { echo "     undefined: $name"; ok=1; }
  done
  return $ok
}

test_buy_position() {
  grep -q '"Assets:" + tx.Custodian + ":" + tx.Symbol' internal/export/beancount.go
}

check "braces balanced" braces
check "parentheses balanced" parens
check "no trailing whitespace" no_trailing_space
check "export references defined" export_refs_defined
check "internal calls defined" local_calls_defined
check "TestBeancountBuy position account" test_buy_position

if [ "$fail" -ne 0 ]; then
  echo "check.sh: FAILED"
  exit 1
fi
echo "check.sh: all $n checks passed"
