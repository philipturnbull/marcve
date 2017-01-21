marcve
======

Generate random CVE decscriptions with Markov chains:

    $ curl https://cve.mitre.org/data/downloads/allitems.xml.gz | gunzip > allitems.xml
    $ marcve --filename allitems.xml &
    $ curl http://127.0.0.1:8080/cve/2019/1531
    Cross-site scripting vulnerability in the sourceAFRICA plugin 0.1.3 for
    WordPress allows remote attackers to bypass intended sandbox restrictions via a
    crafted app that leverages improper handling of O_DIRECT (direct IO) write
    requests.
