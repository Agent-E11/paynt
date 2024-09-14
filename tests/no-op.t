  $ cat "$ASSETDIR/log.txt" | paynt > "$TMPDIR/log.txt"
  $ diff "$ASSETDIR/log.txt" "$TMPDIR/log.txt"

  $ cat "$ASSETDIR/log.txt" | paynt 'does not exist in file/red' > "$TMPDIR/log.txt"
  $ diff "$ASSETDIR/log.txt" "$TMPDIR/log.txt"
