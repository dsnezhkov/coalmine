
![](docs/coalmine.png)
**Coalmine**: Canary Pest Control ;)

### Objective

Less detection, more fun.

### Usage
```
./bin/coalmine pdf -l ~/Downloads  -p 50  -c  -v 
```

```
./bin/coalmine doc  -l ./data -c  -v 
```

```
./bin/coalmine xls  -l ./data -c  -v 
```

```
Seek out canaries and honeytokens in commonly used file formats

Usage:
  coalmine [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  doc         Process doc
  help        Help about any command
  pdf         Process pdfs
  xls         Process xls

Flags:
  -h, --help      help for coalmine
  -v, --verbose   Verbose output

Use "coalmine [command] --help" for more information about a command.

```

PDF:
Canarytokens hide in potentially invalid streams and objects, get them all with brute force (passes):
```
Usage:
  coalmine pdf [flags]

Flags:
  -c, --candidate         Show unverified candidates (default true)
  -h, --help              help for pdf
  -l, --location string   File or Folder to fumigate
  -p, --pass int          Max passes for URL object extraction (default 10)

Global Flags:
  -v, --verbose   Verbose output
```

Gives you known canaries.
Gives you candidates if it cannot validate canaries.
```
✗ ./bin/coalmine pdf -l ~/Downloads  -p 50  -c  
/Users/dimas/Downloads/ssa89-nv.pdf:
        1: (http://www.ssa.gov/cbsv/docs/SampleUserAgreement.pdf)
        2: (http://www.socialsecurity.gov/foia/bluebook)
/Users/dimas/Downloads/j2d9n4auf7b5aeaph3jhlbtp3-1.pdf:
-->     1: (http://j2d9n4auf7b5aeaph3jhlbtp3.canarytokens.net/FADYIMGCFISKUHKSKTQGUOFQTFYHOQITFK)
/Users/dimas/Downloads/gdit4pt6d0kowc3ayq1sp330e.pdf:
-->     1: (http://gdit4pt6d0kowc3ayq1sp330e.canarytokens.net/MPDUBLBAUKKGLHHASBVACRIOPYMNQNYBPV)
/Users/dimas/Downloads/4146936_1_1CGL.pdf:
        1: (http://www.hiscox.com/manage-your-policy)
        2: (http://www.hiscox.com/manage-your-policy)
/Users/dimas/Downloads/detecting_malware_threats.pdf:
        1: (https://www.ibm.com/support/knowledgecenter/en/SSKLLW_9.5.0/com.ibm.bigfix.inventory.doc/Inventory/security/t_checksums_main.html)
/Users/dimas/Downloads/493e490f-394f-4497-9d6e-77c138b14ead.pdf:
        1: (http://www.ssa.gov/oact/COLA/Benefits.html)
        2: (http://www.medicare.gov)
        3: (http://www.ssa.gov/medicare)
        4: (http://www.ssa.gov/people/materials/pdfs/EN-05-10229.pdf)

```
XLS/WORD:

```
Usage:
  coalmine doc [flags]

Flags:
  -c, --candidate         Show unverified candidates (default true)
  -h, --help              help for doc
  -l, --location string   File or Folder to fumigate

Global Flags:
  -v, --verbose   Verbose output

```

```
 ./bin/coalmine doc  -l ./data -c  -v               
data/hhwci3lxddtv9a8bbw1vzp2u5.docx
data/hhwci3lxddtv9a8bbw1vzp2u5.doc
data/hhwci3lxddtv9a8bbw1vzp2u5.docm
 100% |████████████████████████████████████████████████| (20/20, 3833 it/s)
 100% |████████████████████████████████████████████████| (20/20, 3210 it/s)
 100% |████████████████████████████████████████████████| (20/20, 3176 it/s)
data/hhwci3lxddtv9a8bbw1vzp2u5.docx:
-->     1: http://canarytokens.com/feedback/articles/hhwci3lxddtv9a8bbw1vzp2u5/index.html
data/hhwci3lxddtv9a8bbw1vzp2u5.doc:
-->     1: http://canarytokens.com/feedback/articles/hhwci3lxddtv9a8bbw1vzp2u5/index.html
data/hhwci3lxddtv9a8bbw1vzp2u5.docm:
-->     1: http://canarytokens.com/feedback/articles/hhwci3lxddtv9a8bbw1vzp2u5/index.html


```

### Currently supported formats
- PDF
- DOC(X|M)
- XLS(X|M)

### TODO:
- Win Folder
- Redirects
- SQL
- Web bugs
- Cloned sites

As always, if they change we adjust.
