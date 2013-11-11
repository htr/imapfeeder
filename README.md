imapfeeder
==========

Feeds your imap account with feeds.



Quickstart
----------


go get this repository and copy the sample configuration file to your $HOME:

```
go get github.com/htr/imapfeeder

cp $GOPATH/src/github.com/htr/imapfeeder/imapfeeder.json.sample ~/.imapfeeder.json
```

edit the configuration file:

```
{
    "cleanup": true,                // cleanup messages already seen and not flagged
    "folders": {                    // dictionary of folder names mapped to lists of feed urls
        "folder": [
            "feedurl1",
            "feedurl2"
        ],
        "folder2": [
            "feedurl3"
        ]
    },
    "imapserver": "imap.server.tld",  // imap server with tls support
    "from": "<email@server.tld",      // "From" header to be appended to each feeditem Author
    "folderprefix": "feeds/",         // a folder prefix. an empty string is accepted
    "to": "<email@server.tld>",       // "To" header.
    "labels": [],                     // additional labels/folders to where each new message should be appended
    "template": "<table>\n<tbody>\n<tr><td><a href=\"{{ .Link }}\">{{ .Title }}</a></td></tr>\n<hr />\n<tr><td>{{ .Author }}</td></tr>\n<tr><td>{{ .Content }}</td></tr>\n</tbody>\n</table>",     // the message body template
    "password": "password",           // imap username
    "username": "username"            // imap password
}

```


pull all feeds:

```
imapfeeder -pull
```

just do a cleanup:

```
imapfeeder -cleanup
```


test how a feed is processed:

```
imapfeeder -test-feed=feedurl
```


Details
-------


### imap and gmail

When connected to a gmail imap server, imapfeeder uses the `X-GM-EXT-1` features:

 * when *duplicating* messages to other labels/mailboxes: store X-GM-LABELS vs COPY
 * when deleting messages (cleanup): store X-GM-LABELS \Trash vs store \Deleted
 
 
TODO
----

* TODO
* ...








