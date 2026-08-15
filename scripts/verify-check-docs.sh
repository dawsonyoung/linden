#!/usr/bin/env sh
# Verifies that scripts/check-docs.sh actually fails on the conditions it claims
# to detect. Run manually; not part of the gate.
set -u

pass=0
fail=0

report() {
  if [ "$1" = "$2" ]; then
    echo "PASS: $3 (exit $1)"
    pass=$((pass + 1))
  else
    echo "FAIL: $3 (exit $1, expected $2)"
    fail=$((fail + 1))
  fi
}

echo "--- baseline: should pass ---"
sh scripts/check-docs.sh > /dev/null 2>&1
report $? 0 "clean tree passes"

echo "--- unlisted page: should fail ---"
echo "# Orphan" > docs/product/spec/99-orphan.md
sh scripts/check-docs.sh > /tmp/docs-orphan.log 2>&1
code=$?
rm -f docs/product/spec/99-orphan.md
report $code 1 "page absent from SUMMARY is rejected"
grep -q "99-orphan" /tmp/docs-orphan.log && echo "      named the file" || echo "      WARNING: did not name the file"

echo "--- broken link: should fail ---"
cp docs/product/spec/90-glossary.md /tmp/glossary.bak
printf '\n[broken](./does-not-exist.md)\n' >> docs/product/spec/90-glossary.md
sh scripts/check-docs.sh > /tmp/docs-link.log 2>&1
code=$?
cp /tmp/glossary.bak docs/product/spec/90-glossary.md
report $code 1 "broken relative link is rejected"
grep -q "does-not-exist" /tmp/docs-link.log && echo "      named the target" || echo "      WARNING: did not name the target"

echo
echo "passed: $pass  failed: $fail"
[ "$fail" -eq 0 ]
