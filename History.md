# 0.11.0 / 2016-09-13

  * Drop &#34;awsenv-&#34; prefix in username

# 0.10.2 / 2016-05-02

  * Fix: Ensure directories for token-file are created


0.10.1 / 2016-04-17
==================

  * Migrated to go1.6 vendoring

0.10.0 / 2016-04-15
==================

  * Added auto-timeout for LockAgent

0.9.0 / 2016-04-14
==================

  * Added `awsenv passwd` command to change password
  * Updated README, executed `go fmt`

0.8.0 / 2016-02-15
==================

  * Added `run` command

0.7.1 / 2016-02-03
==================

  * Fix: Region variable was broken

0.7.0 / 2015-11-29
==================

  * Added support for AWS\_DEFAULT\_REGION

0.6.0 / 2015-10-01
==================

  * Do not show password in unlock command

0.5.3 / 2015-08-17
==================

  * Fix: Description for AWS region contained variable name

0.5.2 / 2015-08-17
==================

  * Fix: Region parameter had wrong long form

0.5.1 / 2015-06-23
==================

  * Fix: Switched to stable aws-sdk-go

0.5.0 / 2015-06-23
==================

  * Added command &#34;prompt&#34;
  * Added instructions for Homebrew

0.4.3 / 2015-05-29
==================

  * Fix: Shorthand &#34;-d&#34; was duplicated in console command

0.4.2 / 2015-05-28
==================

  * Fix: Persistent flags were not globally available

0.4.1 / 2015-05-28
==================

  * Fix: Do not require unlock for version command

0.4.0 / 2015-05-28
==================

  * Moved to spf13/cobra as CLI framework
  * Added support for direct jump to subconsole
  * Set AWS\_REGION in awsenv shell command

0.3.1 / 2015-05-23
==================

  * Added Godeps

0.3.0 / 2015-05-23
==================

  * Refactored code and followed linter advices
  * Added -u flag in install instructions

0.2.0 / 2015-05-21
==================

  * Added lockagent

0.1.0 / 2015-05-19
==================

  * Initial version