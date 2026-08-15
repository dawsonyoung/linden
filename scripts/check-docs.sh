#!/usr/bin/env sh
# Documentation integrity checks for docs/product.
#
# mdBook silently ignores files absent from SUMMARY.md, so an unpublished page
# looks identical to a published one in the repository. It also does not verify
# relative links. Both are checked here.
#
# Usage: ./scripts/check-docs.sh

set -eu

ROOT=docs/product
SUMMARY="$ROOT/SUMMARY.md"

# Files that are deliberately not published:
#   README.md        contributor guidance about this directory
#   prd/TEMPLATE.md  scaffold for new PRDs
#   SUMMARY.md       the manifest itself
UNPUBLISHED="README.md prd/TEMPLATE.md SUMMARY.md"

if [ ! -f "$SUMMARY" ]; then
  echo "ERROR: $SUMMARY not found"
  exit 1
fi

errors=$(mktemp)
trap 'rm -f "$errors"' EXIT

# 1. Every markdown file is reachable from SUMMARY.md.
for path in $(cd "$ROOT" && find . -name '*.md' | sed 's|^\./||' | sort); do
  case " $UNPUBLISHED " in
    *" $path "*) continue ;;
  esac
  if ! grep -q "($path)" "$SUMMARY"; then
    echo "$ROOT/$path is not listed in SUMMARY.md and would not be published" >> "$errors"
  fi
done

# 2. Every relative markdown link resolves to a file that exists.
for path in $(cd "$ROOT" && find . -name '*.md' | sed 's|^\./||' | sort); do
  file="$ROOT/$path"
  dir=$(dirname "$file")
  grep -o ']([^)]*)' "$file" 2>/dev/null | sed 's|^](||; s|)$||' | while read -r link; do
    case "$link" in
      http://*|https://*|mailto:*|'#'*) continue ;;
    esac
    target=${link%%#*}
    [ -z "$target" ] && continue
    if [ ! -e "$dir/$target" ]; then
      echo "$file links to missing target: $target" >> "$errors"
    fi
  done
done

if [ -s "$errors" ]; then
  echo "Documentation check failed:"
  sed 's/^/  /' "$errors"
  exit 1
fi

echo "Documentation check passed: all pages published, all relative links resolve."
