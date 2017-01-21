marcve
======

Generate random CVE decscriptions with Markov chains:

    $ curl https://cve.mitre.org/data/downloads/allitems.xml.gz | gunzip > allitems.xml
    $ marcve --filename allitems.xml --id CVE-2020-1471 
    Unspecified vulnerability in Cisco Secure Elements Class Name field, which
    allows remote attackers to obtain sensitive environment variable and 2.2.x
    before 9.5.2 and 3.x, allows remote attackers to execute arbitrary web cache
    poisoning attacks via the network traffic after a URL in Linux SCTP data.
