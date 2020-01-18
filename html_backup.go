package goback

func DisplayBackup() string {
	return `{{define "content"}}
				<div class="row">
                    <div class="col">
                        <div  class="panel ">
                            <div class="panel-container show">
                                <div class="panel-content">
                                    <div id="toolbar-backup">
                                    </div>
                                    <table  id="table-streams"
                                            data-toggle="table"
                                            data-toolbar="#toolbar-backup"
                                            data-search="true"
                                            data-url="/summaries"
                                            data-pagination="true"
                                            data-side-pagination="client"
                                            data-show-refresh="true"
                                            data-show-columns="true"
                                            data-sort-name="id"
											data-page-size="20"
                                            data-sort-order="desc">
                                        <thead>
                                        <tr>
                                            <th data-field="id" data-sortable="true">id</th> 
											<th data-field="date" data-sortable="true">date</th> 
											<th data-field="srcDirs" data-sortable="true">srcDirs </th> 
											<th data-field="dstDir" data-sortable="true">dstDir</th> 
											<th data-field="workerCount" data-sortable="true">workerCount </th> 
											<th data-field="totalSize" data-sortable="true">totalSize </th> 
											<th data-field="totalCount" data-sortable="true">totalCount</th> 
											<th data-field="countAdded" data-sortable="true">countAdded</th> 
											<th data-field="countModified" data-sortable="true">countModified </th> 
											<th data-field="countDeleted" data-sortable="true">countDeleted</th> 
											<th data-field="sizeAdded" data-sortable="true">sizeAdded </th> 
											<th data-field="sizeModified" data-sortable="true">sizeModified</th> 
											<th data-field="sizeDeleted" data-sortable="true">sizeDeleted </th> 
											<th data-field="failedCount" data-sortable="true">failedCount </th> 
											<th data-field="failedSize" data-sortable="true">failedSize</th> 
											<th data-field="successCount" data-sortable="true">successCount</th> 
											<th data-field="successSize" data-sortable="true">successSize </th> 
											<th data-field="extensions" data-sortable="true">extensions</th> 
											<th data-field="sizeDistribution" data-sortable="true">sizeDistribution</th> 
											<th data-field="message" data-sortable="true">message </th> 
											<th data-field="version" data-sortable="true">version </th> 
											<th data-field="readingTime" data-sortable="true">readingTime </th> 
											<th data-field="comparisonTime " data-sortable="true">comparisonTime</th> 
											<th data-field="backupTime " data-sortable="true">backupTime</th> 
											<th data-field="loggingTime" data-sortable="true">loggingTime </th> 
											<th data-field="execTime" data-sortable="true">execTime</th> 
                                        </tr>
                                        </thead>
                                    </table>
                                </div>
                            </div>
                        </div>

                    </div>
                </div>
{{end}}

{{define "script"}}
{{end}}
`
}
