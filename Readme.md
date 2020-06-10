# goback

[![Build Status](https://travis-ci.org/devplayg/goback.svg?branch=master)](https://travis-ci.org/devplayg/goback)

Cross-platform incremental backup service

<img src="https://github.com/devplayg/goback/raw/master/goback.png" width="180">

* Running on service
* Embedded web supported (Default :8000)
* Monthly statistics
* Crontab format scheduler provided

## Screenshots

Backup

<img src="https://github.com/devplayg/goback/raw/master/screenshots/backup.png" width="800">

- All files
    
    <img src="https://github.com/devplayg/goback/raw/master/screenshots/modal-backup.png" width="800">
    
- Changes

    <img src="https://github.com/devplayg/goback/raw/master/screenshots/modal-changes.png" width="800">
    
Stats

<img src="https://github.com/devplayg/goback/raw/master/screenshots/stats.png" width="800">

Settings

<img src="https://github.com/devplayg/goback/raw/master/screenshots/settings.png" width="800">


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


