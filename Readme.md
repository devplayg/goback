goback
------

The `goback` makes incremental backups of directories.

    ./goback -s /dir/to/backup -d /dir/to/be/saved
    
    
Usage and examples
-------------------    
    
Run a backup for multiple directories.

    ./goback -s /dir/to/backup1 \
             -s /dir/to/backup2 \
             -d /dir/to/be/saved

The `goback` supports a Web UI for viewing backup results.

    ./goback -w /dir/to/be/saved

* Url: http://127.0.0.1:8000    


Screenshots
------------

Logs

<img src="img/goback-log.jpg">

![텍스트](img/goback-log.jpg)

Changes

Statistics


