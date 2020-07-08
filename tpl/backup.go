package tpl

import "os"

func Backup() string {
	tpl := "tpl/backup.html"
	if _, err := os.Stat(tpl); os.IsExist(err) {
		return displayWithLocalFile(tpl)
	}

	return `{{define "css"}}
    <link rel="stylesheet" media="screen, print" href="/assets/css/custom.css">
{{end}}

{{define "sidebar"}}
    <li class="active">
        <a href="/backup/" title="Backup"><i class="fal fa-list-ul"></i><span class="nav-link-text">Backup logs</span></a>
    </li>
    <li>
        <a href="/stats/" title="Statistics"><i class="fal fa-chart-bar"></i><span class="nav-link-text">Statistics</span></a>
    </li>
    <li>
        <a href="/settings/" title="Settings"><i class="fal fa-cog"></i><span class="nav-link-text">Settings</span></a>
    </li>
    <li>
        <a href="/logout" title="Sign out"><i class="fal fa-sign-out"></i><span class="nav-link-text">Sign out</span></a>
    </li>
{{end}}

{{define "content"}}
    <div class="subheader">
        <h1 class="subheader-title">
            <i class="fal fa-list-ul mr-1"></i> Backup Logs
            <small>
                Backup history
                <div class="float-right">
                    <button class="btn btn-xs btn-secondary">
                        <span class="sysInfo-time" data-format="ll"></span>
                    </button>
                    <button class="btn btn-xs btn-secondary">
                        <span class="sysInfo-time" data-format="LTS"></span>
                    </button>
                </div>
            </small>
        </h1>
    </div>

    <div class="panel">
        <div class="panel-hdr d-none">
            <h2>
                <i class="fal fa-tag mr-1"></i> BACKUP
                <span class="fw-300">
                    <i>Logs</i>
                </span>
            </h2>
            <div class="panel-toolbar">
                <span class="badge badge-primary fw-300 ml-1">Aloha</span>
            </div>
        </div>
        <div class="panel-container show">
            <div class="panel-content">
                <div id="toolbar-backup">
                    <select class="form-control" id="select-srcDirs"></select>
                </div>
                <table  id="table-backup"
                        data-buttons-class="default"
                        class="table table-sm"
                        data-toggle="table"
                        data-cookie="true"
                        data-cookie-id-table="backup"
                        data-cookies-enabled="['bs.table.columns']"
                        data-toolbar="#toolbar-backup"
                        data-search="true"
                        data-url="/summaries"
                        data-pagination="true"
                        data-show-export="true"
                        data-export-types="['csv', 'txt', 'excel']"
                        data-side-pagination="client"
                        data-show-refresh="true"
                        data-show-columns="true"
                        data-sort-name="id"
                        data-page-size="20"
                        data-sort-order="desc">
                    <thead>
                    <tr>
                        <th data-field="id" data-sortable="true">ID</th>
                        <th data-field="backupId" data-sortable="true" data-visible="false">Backup ID</th>
                        <th data-field="date" data-sortable="true" data-formatter="dateFormatter">Date</th>
                        <th data-field="srcDir" data-sortable="true" data-formatter="shortDirFormatter">Directory</th>
                        <th data-field="keeper.protocol" data-sortable="true" data-formatter="backupKeeperFormatter" data-visible="false">Storage</th>
                        <th data-field="backupType" data-sortable="true" data-formatter="backupTypeFormatter">Type</th>
                        <th data-field="workerCount" data-sortable="true" data-visible="false">Workers</th>
                        <th data-field="state" data-sortable="true" data-formatter="backupStateFormatter">State</th>
                        <th data-field="execTime" data-sortable="true" data-formatter="toFixedFormatter">Time(Sec)</th>

                        <th data-field="totalCount" data-sortable="true" data-formatter="backupTotalCountFormatter" data-events="backupStatsEvents">Files</th>
                        <th data-field="totalSize" data-sortable="true" data-formatter="byteSizeFormatter">Total Size</th>

                        <th data-field="countAdded" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents">Added</th>
                        <th data-field="sizeAdded" data-sortable="true" data-formatter="byteSizeFormatter">Added</th>
                        {{/*                        <th data-field="sizeAdded" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Added (B)</th>*/}}

                        <th data-field="countModified" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents">Modified</th>
                        <th data-field="sizeModified" data-sortable="true" data-formatter="byteSizeFormatter">Modified</th>
                        {{/*                        <th data-field="sizeModified" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Modified (B)</th>*/}}

                        <th data-field="countDeleted" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents">Deleted</th>
                        <th data-field="sizeDeleted" data-sortable="true" data-formatter="byteSizeFormatter">Deleted</th>
                        {{/*                        <th data-field="sizeDeleted" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Deleted (B)</th>*/}}

                        <th data-field="countFailed" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents" class="bg-gray-light">Failed</th>
                        <th data-field="sizeFailed" data-sortable="true" data-formatter="byteSizeFormatter" class="bg-gray-light">Failed</th>
                        {{/*                        <th data-field="sizeFailed" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Failed (B)</th>*/}}

                        <th data-field="countSuccess" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents" class="bg-gray-light">Backup Files</th>
                        <th data-field="sizeSuccess" data-sortable="true" data-formatter="byteSizeFormatter" class="bg-gray-light">Backup Size</th>

                        <th data-field="message" data-sortable="false">Read / Compare / Copy / Write</th>

                        <th data-field="readingTime" data-sortable="true" data-formatter="dateFormatter" data-visible="false">1) Reading </th>
                        <th data-field="comparisonTime" data-sortable="true" data-formatter="dateFormatter" data-visible="false">2) Compare</th>
                        <th data-field="backupTime" data-sortable="true" data-formatter="dateFormatter" data-visible="false">3) Backup</th>
                        <th data-field="loggingTime" data-sortable="true" data-formatter="dateFormatter" data-visible="false">4) Logging </th>
                    </tr>
                    </thead>
                </table>
            </div>
        </div>
    </div>

    <div class="modal fade" id="modal-backup-changes" tabindex="-1" role="dialog" aria-hidden="true">
        <div class="modal-dialog mw-100 w-75" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title">
                        Changes #<span class="summary-id"></span>
                        <small class="mb-0">
                            <div class="summary-date"></div>
                            <div class="summary-srcDir"></div>
                        </small>
                    </h1>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true"><i class="fal fa-times"></i></span>
                    </button>
                </div>
                <div class="modal-body">
                    <div class="row">
                        <div class="col-lg-2">
                            <div class="alert alert-primary">
                                <div class="row">
                                    <div class="col-sm-6 col-md-6 col-lg-12">
                                        <div>
                                            <div class="badge border border-primary text-primary">Files</div>
                                            <div class="summary-totalCount h2"></div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Size</div>
                                            <div class="summary-totalSize h2"></div>
                                        </div>
                                        <div class="mt-2 mb-3">
                                            <div class="badge border border-primary text-primary">Execution time</div>
                                            <div class="summary-execTime h2 mb-0"></div>
                                            <div class="pl-2">- Reading: <span class="summary-readingTime 53"></span></div>
                                            <div class="pl-2">- Comparing: <span class="summary-comparingTime 53"></span></div>
                                            <div class="pl-2">- Coping: <span class="summary-backupTime 53"></span></div>
                                            <div class="pl-2">- Logging: <span class="summary-loggingTime 53"></span></div>
                                        </div>
                                    </div>

                                    <div class="col-sm-6 col-md-6 col-lg-12">
                                        <div>
                                            <div class="badge border border-primary text-primary">Added</div>
                                            <div>
                                                <span class="h2 summary-countAdded"></span>
                                                (<span class="summary-sizeAdded"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Modified</div>
                                            <div>
                                                <span class="h2 summary-countModified"></span>
                                                (<span class="summary-sizeModified"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Deleted</div>
                                            <div>
                                                <span class="h2 summary-countDeleted"></span>
                                                (<span class="summary-sizeDeleted"></span>)
                                            </div>
                                        </div>

                                        <div class="mt-2">
                                            <div class="badge border border-secondary text-secondary">Backup success</div>
                                            <div>
                                                <span class="h2 summary-countSuccess"></span>
                                                (<span class="summary-sizeSuccess"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-danger text-danger">Backup failure</div>
                                            <div>
                                                <span class="h2 summary-countFailed"></span>
                                                (<span class="summary-sizeFailed"></span>)
                                            </div>
                                        </div>

                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="col-lg-10">
                            <ul id="tabs-backup-changes" class="nav nav-tabs" role="tablist">
                                <li class="nav-item">
                                    <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-added" role="tab">
                                        Added
                                        <span class="badge border border-info rounded-pill bg-primary-500 stats-data stats-added"></span>
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-modified" role="tab">
                                        Modified
                                        <span class="badge border border-info rounded-pill bg-success-500 stats-data stats-modified"></span>
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-deleted" role="tab">
                                        Deleted
                                        <span class="badge border border-warning rounded-pill bg-warning-500 stats-data stats-deleted"></span>
                                    </a>
                                </li>
                                <li class="nav-item">
                                    <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-failed" role="tab">
                                        Failed
                                        <span class="badge border border-danger rounded-pill bg-danger-500 stats-data stats-failed"></span>
                                    </a>
                                </li>
                            </ul>
                            <div class="tab-content p-3">
                                <div class="tab-pane fade show active" id="tab-backup-added" role="tabpanel">
                                    <div class="row">
                                        <div class="col-lg-9">
                                            <div id="toolbar-backup-added">
                                                <h3 class="stats-data stats-added-size"></h3>
                                            </div>
                                            <table  id="table-backup-added"
                                                    data-buttons-class="default"
                                                    class="table table-data table-sm"
                                                    data-toolbar="#toolbar-backup-added"
                                                    data-toggle="table"
                                                    data-search="true"
                                                    data-pagination="true"
                                                    data-show-export="true"
                                                    data-pagination-v-align="bottom"
                                                    data-show-columns="true"
                                                    data-export-types="['csv', 'txt', 'excel']"
                                                    data-side-pagination="client"
                                                    data-sort-name="size"
                                                    data-page-size="15"
                                                    data-row-style="backupRowStyle"
                                                    data-sort-order="desc">
                                                <thead>
                                                <tr>
                                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                                    <th data-field="mtime" data-visible="false" data-formatter="dateFormatter">ModTime</th>
                                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                                </tr>
                                                </thead>
                                            </table>
                                        </div>
                                        <div class="col-lg-3">
                                            <div class="stats-data stats-added-ext"></div>
                                        </div>
                                    </div>
                                </div>
                                <div class="tab-pane fade" id="tab-backup-modified" role="tabpanel">
                                    <div class="row">
                                        <div class="col-lg-9">
                                            <div id="toolbar-backup-modified">
                                                <h3 class="stats-data stats-modified-size"></h3>
                                            </div>
                                            <table  id="table-backup-modified"
                                                    data-buttons-class="default"
                                                    class="table table-data table-sm"
                                                    data-toolbar="#toolbar-backup-modified"
                                                    data-toggle="table"
                                                    data-search="true"
                                                    data-pagination="true"
                                                    data-show-export="true"
                                                    data-pagination-v-align="bottom"
                                                    data-show-columns="true"
                                                    data-export-types="['csv', 'txt', 'excel']"
                                                    data-side-pagination="client"
                                                    data-sort-name="size"
                                                    data-page-size="15"
                                                    data-row-style="backupRowStyle"
                                                    data-sort-order="desc">
                                                <thead>
                                                <tr>
                                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                                    <th data-field="mtime" data-visible="false" data-formatter="dateFormatter">ModTime</th>
                                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                                </tr>
                                                </thead>
                                            </table>
                                        </div>
                                        <div class="col-lg-3">
                                            <div class="stats-data stats-modified-ext"></div>
                                        </div>
                                    </div>
                                </div>
                                <div class="tab-pane fade" id="tab-backup-deleted" role="tabpanel">
                                    <div class="row">
                                        <div class="col-lg-9">
                                            <div id="toolbar-backup-deleted">
                                                <h3 class="stats-data stats-deleted-size"></h3>
                                            </div>
                                            <table  id="table-backup-deleted"
                                                    data-buttons-class="default"
                                                    class="table table-data table-sm"
                                                    data-toolbar="#toolbar-backup-deleted"
                                                    data-toggle="table"
                                                    data-search="true"
                                                    data-pagination="true"
                                                    data-show-export="true"
                                                    data-pagination-v-align="bottom"
                                                    data-show-columns="true"
                                                    data-export-types="['csv', 'txt', 'excel']"
                                                    data-side-pagination="client"
                                                    data-sort-name="size"
                                                    data-page-size="15"
                                                    data-row-style="backupRowStyle"
                                                    data-sort-order="desc">
                                                <thead>
                                                <tr>
                                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                                    <th data-field="mtime" data-visible="false" data-formatter="dateFormatter">ModTime</th>
                                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                                </tr>
                                                </thead>
                                            </table>
                                        </div>
                                        <div class="col-lg-3">
                                            <div class="stats-data stats-deleted-ext"></div>
                                        </div>
                                    </div>
                                </div>
                                <div class="tab-pane fade" id="tab-backup-failed" role="tabpanel">
                                    <div class="row">
                                        <div class="col-lg-9">
                                            <div id="toolbar-backup-failed">
                                                <h3 class="stats-data stats-failed-size"></h3>
                                            </div>
                                            <table  id="table-backup-failed"
                                                    data-buttons-class="default"
                                                    class="table table-data table-sm"
                                                    data-toolbar="#toolbar-backup-failed"
                                                    data-toggle="table"
                                                    data-search="true"
                                                    data-pagination="true"
                                                    data-show-export="true"
                                                    data-pagination-v-align="bottom"
                                                    data-show-columns="true"
                                                    data-export-types="['csv', 'txt', 'excel']"
                                                    data-side-pagination="client"
                                                    data-sort-name="size"
                                                    data-page-size="15"
                                                    data-row-style="backupRowStyle"
                                                    data-sort-order="desc">
                                                <thead>
                                                <tr>
                                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                                    <th data-field="mtime" data-visible="false" data-formatter="dateFormatter">ModTime</th>
                                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                                </tr>
                                                </thead>
                                            </table>
                                        </div>
                                        <div class="col-lg-3">
                                            <div class="stats-data stats-failed-ext"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-primary" data-dismiss="modal">Close</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade" id="modal-backup-stats" tabindex="-1" role="dialog" aria-hidden="true">
        <div class="modal-dialog mw-100 w-75" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title">
                        Backup #<span class="summary-id"></span>
                        <small class="mb-0">
                            <div class="summary-date"></div>
                            <div class="summary-srcDir"></div>
                        </small>
                    </h1>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true"><i class="fal fa-times"></i></span>
                    </button>
                </div>
                <div class="modal-body">
                    <div class="row">
                        <div class="col-lg-2">
                            <div class="alert alert-primary">
                                <div class="row">
                                    <div class="col-sm-6 col-md-6 col-lg-12">
                                        <div>
                                            <div class="badge border border-primary text-primary">Files</div>
                                            <div class="summary-totalCount h2"></div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Size</div>
                                            <div class="summary-totalSize h2"></div>
                                        </div>
                                        <div class="mt-2 mb-3">
                                            <div class="badge border border-primary text-primary">Execution time</div>
                                            <div class="summary-execTime h2 mb-0"></div>
                                            <div class="pl-2">- Reading: <span class="summary-readingTime 53"></span></div>
                                            <div class="pl-2">- Comparing: <span class="summary-comparingTime 53"></span></div>
                                            <div class="pl-2">- Coping: <span class="summary-backupTime 53"></span></div>
                                            <div class="pl-2">- Logging: <span class="summary-loggingTime 53"></span></div>
                                        </div>
                                    </div>

                                    <div class="col-sm-6 col-md-6 col-lg-12">
                                        <div>
                                            <div class="badge border border-primary text-primary">Added</div>
                                            <div>
                                                <span class="h2 summary-countAdded"></span>
                                                (<span class="summary-sizeAdded"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Modified</div>
                                            <div>
                                                <span class="h2 summary-countModified"></span>
                                                (<span class="summary-sizeModified"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-primary text-primary">Deleted</div>
                                            <div>
                                                <span class="h2 summary-countDeleted"></span>
                                                (<span class="summary-sizeDeleted"></span>)
                                            </div>
                                        </div>

                                        <div class="mt-2">
                                            <div class="badge border border-secondary text-secondary">Backup success</div>
                                            <div>
                                                <span class="h2 summary-countSuccess"></span>
                                                (<span class="summary-sizeSuccess"></span>)
                                            </div>
                                        </div>
                                        <div class="mt-2">
                                            <div class="badge border border-danger text-danger">Backup failure</div>
                                            <div>
                                                <span class="h2 summary-countFailed"></span>
                                                (<span class="summary-sizeFailed"></span>)
                                            </div>
                                        </div>

                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="col-lg-5">
                            <!-- Extension Ranking -->
                            <div class="panel">
                                <div class="panel-container show">
                                    <div class="panel-content">
                                        <div id="toolbar-table-backup-stats-ext">
                                            <div class="h3 text-primary">By extension</div>
                                        </div>
                                        <table  id="table-backup-stats-ext"
                                                data-toolbar="#toolbar-table-backup-stats-ext"
                                                data-buttons-class="default"
                                                class="table table-data table-sm"
                                                data-toggle="table"
                                                data-search="true"
                                                data-pagination="true"
                                                data-show-export="true"
                                                data-pagination-v-align="bottom"
                                                data-export-types="['csv', 'txt', 'excel']"
                                                data-side-pagination="client"
                                                data-sort-name="size"
                                                data-page-size="15"
                                                data-sort-order="desc">
                                            <thead>
                                            <tr>
                                                <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                                <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Total Size</th>
                                                <th data-field="count" data-formatter="thCommaFormatter" data-sortable="true">Count</th>
                                            </tr>
                                            </thead>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="col-lg-5">
                            <!-- Size Distribution -->
                            <div class="panel">
                                <div class="panel-container show">
                                    <div class="panel-content">
                                        <div id="toolbar-table-backup-stats-size">
                                            <div class="h3 text-primary">By size</div>
                                        </div>
                                        <table  id="table-backup-stats-size"
                                                data-buttons-class="default"
                                                class="table table-data table-sm"
                                                data-toggle="table"
                                                data-toolbar="#toolbar-table-backup-stats-size"
                                                data-search="true"
                                                data-pagination="true"
                                                data-show-export="true"
                                                data-pagination-v-align="bottom"
                                                data-export-types="['csv', 'txt', 'excel']"
                                                data-side-pagination="client"
                                                data-sort-name="sizeDist"
                                                data-page-size="15"
                                                data-sort-order="desc">
                                            <thead>
                                            <tr>
                                                <th data-field="sizeDist" data-sortable="true" data-formatter="backupStatsSizeDistFormatter">Size</th>
                                                <th data-field="size" data-formatter="byteSizeFormatter" data-sortable="true">Total Size</th>
                                                <th data-field="count" data-formatter="thCommaFormatter" data-sortable="true">Count</th>
                                            </tr>
                                            </thead>
                                        </table>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-primary" data-dismiss="modal">Close</button>
                </div>
            </div>
        </div>
    </div>
{{end}}

{{define "script"}}
    <script src="/assets/js/custom.js"></script>
    <script>
        function backupRowStyle(row, index) {
            if (row.state === -1) {
                return {
                    classes: 'bg-warning-100'
                }
            }
            return {
                css: {
                    // color: 'blue'
                }
            }
        }

        function countWithRateFormatter(val, row, idx) {
            return val + '<span class="small float-right text-muted">' + row.rate + '%</span>'

        }

        function showChangedFiles(id, field) {
            let url = "/summaries/" + id + "/changes";
            $.ajax({
                url: url,
            }).done(function(report) {
                updateBackupStats("added", report.added, "primary");
                updateBackupStats("modified", report.modified, "success");
                updateBackupStats("deleted", report.deleted, "warning");
                updateBackupStats("failed", report.failed, "danger");

                $("#modal-backup-changes").modal("show");
                $('#tabs-backup-changes a[href="#tab-backup-' + getTab(field) + '"]').tab('show');
            }).fail(function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR);
                let msg = "";
                if (jqXHR.responseJSON !== undefined && jqXHR.responseJSON.error !== undefined) {
                    msg = jqXHR.responseJSON.error;
                } else {
                    msg = textStatus;
                }
                Swal.fire("Error", msg, "warning");
            });
        }

        function showStats(id) {
            let url = "/summaries/" + id;
            $.ajax({
                url: url,
            }).done(function(summary) {
                let total = 0;
                $.each(summary.stats.extRanking, function(i, r) {
                    total += r.count;
                });
                if (total > 0) {
                    $.each(summary.stats.extRanking, function(i, r) {
                        summary.stats.extRanking[i].rate = (r.count / total *100).toFixed(2);
                    });
                }
                $("#table-backup-stats-ext").bootstrapTable("load", summary.stats.extRanking);
                $("#table-backup-stats-size").bootstrapTable("load", summary.stats.sizeDist);
                $("#modal-backup-stats").modal("show");
            }).fail(function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR);
                let msg = "";
                if (jqXHR.responseJSON !== undefined && jqXHR.responseJSON.error !== undefined) {
                    msg = jqXHR.responseJSON.error;
                } else {
                    msg = textStatus;
                }
                Swal.fire("Error", msg, "warning");
            });
        }

        function updateBackupStats(how, what, colorSuffix) {
            $("#table-backup-" + how).bootstrapTable("load", what.files);
            $(".stats-" + how).text(what.files.length > 0 ? what.files.length : "");
            $(".stats-" + how + "-size").html(bytesToSize(what.size) + " <i>" + how + "</i>");

            let extTags = "",
                total = what.report.extRanking.length;

            $.each(what.report.extRanking, function(i, r) {
                extTags += '<a href="#" class="filterFiles" data-cond="' + r.ext + '" data-how="' + how + '"><span class="badge badge-stats bg-' + colorSuffix + '-' + getRate(i, total) + '">'+ r.ext + " / " + bytesToSize(r.size) + '</span></a>';
            });
            $(".stats-" + how + "-ext" ).html(extTags);
        }

        function getTab(field) {
            let tabs = ["Added", "Modified", "Deleted", "Failed"],
                selected = null;

            $.each(tabs, function(i, t) {
                if (field.endsWith(t)) {
                    selected = t;
                    return false;
                }
            });
            return (selected === null) ? tabs[0].toLowerCase() : selected.toLowerCase();
        }

        /*
        * Formatters
        */

        $("#modal-backup-changes, #modal-backup-stats")
            .on("hidden.bs.modal", function (e) {
                $(".table-data").bootstrapTable("load", []);
                $(".table-data").bootstrapTable("filterBy", {}); // reset filter
                $(".stats-data").empty();
            });

        $(".table")
            .on("refresh.bs.table", function (e) {
            })
            .on("all.bs.table", function (e, data) {
                $('.has-tooltip').tooltip();
            })
            .on("load-success.bs.table", function(e, rows) {
                let srcDir = {};
                $.each(rows, function(i, s) {
                    srcDir[s.srcDir] = true;
                });
                Object.keys(srcDir).sort();

                $("#select-srcDirs").empty();
                $("<option/>", {
                    "value": "",
                }).text("Directories").appendTo($("#select-srcDirs"));
                $.each(srcDir, function(srcDir, r) {
                    $("<option/>", {
                        "value": srcDir
                    }).text(srcDir).appendTo($("#select-srcDirs"));
                });
            });

        window.backupStatsEvents = {

            'click .stats': function (e, val, row, idx) {
                updateSummary(row);
                showStats(row.id);
            },
            'click .changed': function (e, val, row, idx) {
                let $btn = $(e.currentTarget);
                updateSummary(row);
                showChangedFiles(row.id, $btn.data("field"));
            },
        };

        function updateSummary(row) {
            $(".summary-id").text(row.id);
            $(".summary-totalCount").text(thousandCommaSep(row.totalCount));
            $(".summary-totalSize").text(bytesToSize(row.totalSize));
            $(".summary-srcDir").text(row.srcDir);
            $(".summary-date").html('<span class="text-soft">' + moment(row.date).format("lll") + '</span>');

            $(".summary-countAdded").text(row.countAdded);
            $(".summary-sizeAdded").html(bytesToSize(row.sizeAdded));
            $(".summary-countModified").text(row.countModified);
            $(".summary-sizeModified").html(bytesToSize(row.sizeModified));
            $(".summary-countDeleted").text(row.countDeleted);
            $(".summary-sizeDeleted").html(bytesToSize(row.sizeDeleted));
            $(".summary-countSuccess").text(row.countSuccess);
            $(".summary-sizeFailed").html(bytesToSize(row.sizeFailed));
            $(".summary-countFailed").text(row.countFailed);

            $(".summary-sizeSuccess").html(bytesToSize(row.sizeSuccess));

            $(".summary-execTime").text(row.execTime.toFixed(2) + " sec");
            $(".summary-readingTime").text( (moment.parseZone(row.readingTime).diff(moment.parseZone(row.date)) / 1000).toFixed(2) );
            $(".summary-comparingTime").text( (moment.parseZone(row.comparisonTime).diff(moment.parseZone(row.readingTime)) / 1000).toFixed(2) );
            $(".summary-backupTime").text( (moment.parseZone(row.backupTime).diff(moment.parseZone(row.comparisonTime)) / 1000).toFixed(2) );
            $(".summary-loggingTime").text( (moment.parseZone(row.loggingTime).diff(moment.parseZone(row.backupTime)) / 1000).toFixed(2) );
        }

        $( "#select-srcDirs" ).change(function () {
            if (this.value.length > 0) {
                $("#table-backup").bootstrapTable('filterBy', {
                    srcDir: [this.value]
                });
                return;
            }

            $("#table-backup").bootstrapTable('filterBy',{});
        });

        let filteredChecker = {};
        $('body').on('click', 'a.filterFiles', function() {
            let ext = $(this).data("cond"),
                how = $(this).data("how");

            if (filteredChecker[how] !== undefined) {
                if (filteredChecker[how] === ext) {
                    $("#table-backup-" + how).bootstrapTable('filterBy', {} );
                    filteredChecker[how] = null;
                    return;
                }
            }
            filteredChecker[how] = ext;
            $("#table-backup-" + how).bootstrapTable('filterBy', {
                ext: [ext]
            });
        });
    </script>
{{end}}
`
}
