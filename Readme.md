# Backup

summary.db

files-[dir-checksum].db


### modal
	
- title, changes log
	
### display

- page
    - title change(rtsp -> goback)
    - backup icon

- main table
    - table col color, failed column,
    - table msg col, tooltip
    - select box
    - long name problem, backup dir
    - msg: first backup, 2.4s / -9223372036.9s / 0.0s / 9223372036.9s
    - main table toolbar summary (n backup)
    - filterBy, https://bootstrap-table.com/docs/api/methods/#filterby
    
- sub table
    - col seq => name, size, ext, msg***

- tx

    sftp
    
    
    $button.click(function () {
        $table.bootstrapTable('filterBy', {
            id: [1, 2, 3]
        })
    })
    
