Color log

  $ cat "$ASSETDIR/log.txt" | paynt '0/^error$/red' '0/^warn$/yellow' '0/^info$/blue'
  \x1b[31merror this is a test\x1b[0m (esc)
  \x1b[34minfo  this is some info\x1b[0m (esc)
  \x1b[33mwarn  this is a warning\x1b[0m (esc)
  \x1b[34minfo  bla bla bla\x1b[0m (esc)
  \x1b[31merror something bad happened\x1b[0m (esc)
