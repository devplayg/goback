# goback

Cross-platform incremental backup service

<img src="goback.png" width="200"">

* Running on service
* Embedded web supported (Default :8000)
* Monthly statistics
* Crontab format scheduler provided


### Database

Powered by BoltDB

1. `summary`

* key: Summary.Id
* value: Summary

2) `backup`

* key: Backup.Id
* value: nil

3. `config`

|Key|Value|
|---|---|
|backup|Config|
|backup_checksum|sha256(Config)|


