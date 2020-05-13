# goback

json
struct
summary.db.lock - gob encoded
summary.db - zipped

powered by **BoltDB**

- s.runBackupJob(jobId)
	- backup.start()
		- b.startDirBackup(dir)
			- b.startBackup(srcDir, lastFileMap)
				- b.issueSummary(srcDir, Incremental)
				- b.getCurrentFileMaps(srcDir)
				- b.backupFiles()
				- b.writeResult(
					- b.writeChangesLog(lastFileMap)
					- b.writeFileMaps(currentFileMaps)
    - s.writeSummaries(summaries)
	- backupId := lastBackupId + 1
	
	
## Buckets of DB

`summary`

* key: Summary.Id
* value: Summary

`backup`

* key: Backup.Id
* value: nil

`config`

|Key|Value|
|---|---|
|storage|JSON|