

## Backup Process

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
	
	
	
## To do

- [ ] Monthly report
- [ ] Custom logo image
- [x] ID/PW settings
- [x] Loading scheduler on starting
- [ ] Displaying disk partition usage
- [ ] Graphing file growth
- [ ] Adding build note
- [ ] Changing badge to button on UI
