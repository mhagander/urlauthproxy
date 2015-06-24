urlauthproxy
============

`urlauthproxy` is a trivial proxy that fakes authentication using secret URLs.
Any request to a secret URL will be proxied to a different (possibly also secret)
URL. There are a few usecases for this, such as:

* Getting access to a basic-authenticated resource from a service that does not
  support username/password given on the URL for access to such (looking at you,
  Google Calendar)
* Being able to expose a single URL that's protected behind a password that is also
  used for other things (such as the rest of the website, other calendars, etc)
  
`urlauthproxy` will only serve requests over encrypted https connections. It will
call out to clients regardless of https or not, so be careful with what's written
in the map.

urlmap
------
`urlauthproxy` uses a file called `urlmap` to map secret urls to the resources
it access. This is a very simple file of format

```
/secret/url/here=https://user:secret@some.where.com/whatever
```

Be careful not to put spaces around the `=` sign, as that will be interpreted
as part of the URL.

certificates
------------
`urlauthproxy` expects a certificate file and a key file, both in standard PEM
format. Commandline arguments control where to find these files.

commandline
-----------
Just use `urlauthproxy --help` to get a quick help.
