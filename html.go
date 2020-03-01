package goback

import (
	"fmt"
	"io/ioutil"
)

func DisplayWithLocalFile(name string) string {
	b, err := ioutil.ReadFile(fmt.Sprintf("static/%s.html", name))
	if err != nil {
		log.Error(err)
	}
	return string(b)
}

func DisplayBackup() string {
	return `{{define "css"}}
    <style>
        .pagination .page-link {
            border-width: 1px;
        }

        .badge-stats {
            margin-right: 3px;
            font-weight: 400;
        }

        @media print{@page {size: landscape}}
    </style>
{{end}}

{{define "sidebar"}}
    <li>
        <a href="/backup/" title="Backup"><i class="fal fa-database"></i><span class="nav-link-text">Backup</span></a>
    </li>
    <li>
        <a href="/settings/" title="Settings"><i class="fal fa-cog"></i><span class="nav-link-text">Settings</span></a>
    </li>
{{end}}

{{define "content"}}
    <div class="panel">
        <div class="panel-hdr">
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
                        <th data-field="date" data-sortable="true" data-formatter="dateFormatter">Date</th>
                        <th data-field="srcDir" data-sortable="true" data-formatter="shortDirFormatter">Source</th>
                        <th data-field="keeper.protocol" data-sortable="true" data-formatter="backupKeeperFormatter">Storage</th>
                        <th data-field="backupType" data-sortable="true" data-formatter="backupTypeFormatter">Type</th>
                        <th data-field="state" data-sortable="true" data-formatter="backupStateFormatter">State</th>
                        <th data-field="execTime" data-sortable="true" data-formatter="toFixedFormatter">Time(Sec)</th>

                        <th data-field="workerCount" data-sortable="true" data-visible="false">Workers</th>

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

                        <th data-field="countFailed" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupStatsEvents">Failed</th>
                        <th data-field="sizeFailed" data-sortable="true" data-formatter="byteSizeFormatter">Failed</th>
                        {{/*                        <th data-field="sizeFailed" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Failed (B)</th>*/}}


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
                    <h2 class="modal-title">
                        <i class="fal fa-file-alt"></i> Changes
                        <i class="summary fw-300 small ml-2"></i>
                    </h2>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true"><i class="fal fa-times"></i></span>
                    </button>
                </div>
                <div class="modal-body">
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
                                <div class="col-lg-8">
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
                                <div class="col-lg-4">
                                    <div class="stats-data stats-added-ext"></div>
                                </div>
                            </div>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-modified" role="tabpanel">
                            <div class="row">
                                <div class="col-lg-8">
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
                                <div class="col-lg-4">
                                    <div class="stats-data stats-modified-ext"></div>
                                </div>
                            </div>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-deleted" role="tabpanel">
                            <div class="row">
                                <div class="col-lg-8">
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
                                <div class="col-lg-4">
                                    <div class="stats-data stats-deleted-ext"></div>
                                </div>
                            </div>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-failed" role="tabpanel">
                            <div class="row">
                                <div class="col-lg-8">
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
                                <div class="col-lg-4">
                                    <div class="stats-data stats-failed-ext"></div>
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
                    <h2 class="modal-title">
                        <i class="fal fa-sort-amount-down"></i> Statistics
                        <i class="summary fw-300 small ml-2"></i>
                    </h2>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true"><i class="fal fa-times"></i></span>
                    </button>
                </div>
                <div class="modal-body">
                    <div class="row">
                        <div class="col">
                            <!-- Extension Ranking -->
                            <div class="panel">
                                <div class="panel-hdr">
                                    <h2>
                                        Rankings by <i>extension</i>
                                    </h2>
                                    <div class="panel-toolbar">
                                        <span class="badge badge-primary fw-300 ml-1">Statistics</span>
                                    </div>
                                </div>
                                <div class="panel-container show">
                                    <div class="panel-content">
                                        <table  id="table-backup-stats-ext"
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
                        <div class="col">
                            <!-- Size Distribution -->
                            <div class="panel">
                                <div class="panel-hdr">
                                    <h2>
                                        Rankings by <i>fileSize distribution</i>
                                    </h2>
                                    <div class="panel-toolbar">
                                        <span class="badge badge-primary fw-300 ml-1">Statistics</span>
                                    </div>
                                </div>
                                <div class="panel-container show">
                                    <div class="panel-content">
                                        <table  id="table-backup-stats-size"
                                                data-buttons-class="default"
                                                class="table table-data table-sm"
                                                data-toggle="table"
                                                data-toolbar="#toolbar-backup-stats-size"
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
    <script src="/assets/plugins/bootstrap-select/bootstrap-select.min.js"></script>
    <script>
        function backupRowStyle(row, index) {
            let classes = [
                'bg-blue',
                'bg-green',
                'bg-orange',
                'bg-yellow',
                'bg-red'
            ]

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

        function extFormatter(val, row, idx) {
            if (val.length < 1) {
                return;
            }
            return val;
        }

        function dateFormatter(val, row, idx) {
            return moment(val).format("YYYY-MM-DD HH:mm:ss");
        }

        function bytesToSize(bytes) {
            return humanizedSize(bytes, true, 1);
        }

        function humanizedSize(bytes, si, toFixed) {
            let thresh = si ? 1000 : 1024;
            if(Math.abs(bytes) < thresh) {
                return bytes + ' B';
            }
            let units = si
                ? ['kB','MB','GB','TB','PB','EB','ZB','YB']
                : ['KiB','MiB','GiB','TiB','PiB','EiB','ZiB','YiB'];
            let u = -1;
            do {
                bytes /= thresh;
                ++u;
            } while(Math.abs(bytes) >= thresh && u < units.length - 1);
            return bytes.toFixed(toFixed)+' '+units[u];
        }

        function basename(path) {
            return path.replace(/^.*[\\\/]/, '');
        }
        // function dirname(path) {
        //     // return path.substr(0, basename(path).lastIndexOf('.'));
        //
        //     return path.substring(0, path.lastIndexOf(basename(path)));
        //
        //     // return path.replacee
        //     // trimEnd(path, basename(path));
        // }

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
            });
        }

        function showStats(backup) {
            $("#table-backup-stats-ext").bootstrapTable("load", backup.stats.extRanking);
            $("#table-backup-stats-size").bootstrapTable("load", backup.stats.sizeDist);
            $("#modal-backup-stats").modal("show");
        }

        function updateBackupStats(how, what, colorSuffix) {
            $("#table-backup-" + how).bootstrapTable("load", what.files);
            $(".stats-" + how).text(what.files.length > 0 ? what.files.length : "");
            $(".stats-" + how + "-size").html(bytesToSize(what.size) + " <i>" + how + "</i>");

            let extTags = "",
                total = what.report.extRanking.length;

            $.each(what.report.extRanking, function(i, r) {
                extTags += '<a href="#" class="filterFiles" data-cond="' + r.ext + '" data-how="' + how + '"><span class="badge badge-stats bg-' + colorSuffix + '-' + getPer(i, total) + '">'+ r.ext + " / " + bytesToSize(r.size) + '</span></a>';
            });
            $(".stats-" + how + "-ext" ).html(extTags);
        }

        function getPer(i, total) {
            let per =  Math.round((1 - (i / total)) * 100);
            per = (per - (per % 10) - 20) * 10;
            if (per < 100) {
                per = 50;
            }
            return per;
            ss
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

        function backupKeeperFormatter(val, row, idx) {
            if (val === 1) {
                return "Local disk";
            }
            if (val === 2) {
                return "Remote (SFTP)";
            }
            if (val === 4) {
                return "Remote (SFTP)";
            }
        }

        function backupStatsSizeDistFormatter(val, row, idx) {
            if (val >= 5000000000000) {
                return "Big file";
            }
            if (val === 0) {
                return val;
            }
            return '<span class="has-tooltip" title="' + bytesToSize(val / 10) + ' ~ ' +  bytesToSize(val) + '">' + humanizedSize(val / 10, true, 0) + " ~ " +  humanizedSize(val, true, 0) + '</span>';
        }

        function backupTotalCountFormatter(val, row, idx) {
            return '<a href="javascript:void(0);" class="stats">' + thCommaFormatter(val, row, idx) + '</a>';
        }

        function shortDirFormatter(val, row, idx) {
            if (val.length < 21) {
                return val;
            }
            let dir = basename(val);
            return '<span class="has-tooltip" title="' + val + '">.. ' + dir + '</span>';
        }

        function sizeBFormatter(val, row, idx) {
            return row.size.toLocaleString();
        }

        function backupFileStateFormatter(val, row, idx) {
            if (val === -1) {
                return '<i><span class="text-danger">failed</span><i>';
            }
            if (val === 1) {
                return "<i>done</i>";
            }
        }
        function backupChangesNameFormatter(val, row, idx) {
            return '<span class="has-tooltip" title="' + row.dir + '">' + val + '</span>';
        }
        function byteSizeFormatter(val, row, idx) {
            if (val === 0) {
                return '<span class="text-muted ">' + bytesToSize(val) + '</span>';
            }
            if (val < 1000) {
                return bytesToSize(val);
            }
            return '<span class="has-tooltip" title="' + val.toLocaleString() + ' Bytes">' + bytesToSize(val) + '</span>';
        }

        function backupResultFormatter(val, row, idx, field) {
            if (val === 0) {
                return '<span class="text-muted">' + val + '</span>';
            }

            // return '<a href="javascript:void(0);" class="file">' + val.toLocaleString() + '</a>';
            let link = '<a class="changed" href="javascript:void(0);" title="changed" data-field="' + field + '">' + val.toLocaleString() + '</a>';
            return link;


            // let th = $('#table-backup').find("[data-field='" + field + "']");
            // // console.log(th.text());
            // let link = $("<a/>", {
            //     href: "javascript:void(0);",
            //     class: "file",
            //     "data-title": th.text(),
            //     "data-field": field,
            //     "title": "",
            // }).html(
            //     val.toLocaleString()
            // );
            // return link[0].outerHTML;
        }

        function backupTypeFormatter(val, row, idx) {
            if (val === 1) {
                return '<span class="badge badge-primary">Initial</span>';
            }
            if (val === 2) {
                return 'Incremental';
            }
        }

        function backupStateFormatter(val, row, idx) {
            if (val === 5) {
                return 'Completed';
            }
            return val;
        }

        function toFixedFormatter(val, row, idx) {
            return  val.toFixed(2);
        }

        function thCommaFormatter(val, row, idx) {
            if (val === 0) {
                return '<span class="text-muted">0</span>';
            }
            return val.toLocaleString();
        }

        $("#modal-backup-changes, #modal-backup-stats")
            .on('show.bs.modal', function (e) {

            })
            .on('hidden.bs.modal', function (e) {
                $(".table-data").bootstrapTable("load", []);
                $(".table-data").bootstrapTable("filterBy", {}); // reset filter
                $(".stats-data").empty();
            });

        $(".table")
            .on('all.bs.table', function (e, data) {
                $('.has-tooltip').tooltip();
            })
            .on('load-success.bs.table', function(e, rows) {
                let srcDir = {};
                $.each(rows, function(i, s) {
                    srcDir[s.srcDir] = true;
                });
                Object.keys(srcDir).sort();

                $("#select-srcDirs").empty();
                $("<option/>", {
                    "value": "",
                }).text("All").appendTo($("#select-srcDirs"));
                $.each(srcDir, function(srcDir, r) {
                    $("<option/>", {
                        "value": srcDir
                    }).text(srcDir).appendTo($("#select-srcDirs"));
                });
            });

        window.backupStatsEvents = {
            'click .stats': function (e, val, row, idx) {
                // let $btn = $(e.currentTarget);
                let tags = '<button type="button" class="btn btn-sm btn-outline-default mr-1">Backup ID: ' + row.id + '</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default mr-1">' + row.totalCount + ' files</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default mr-1">' + bytesToSize(row.totalSize) + '</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default">' + row.srcDir + '</button>';

                $("#modal-backup-stats .modal-title .summary").html(tags);
                showStats(row);
            },
            'click .changed': function (e, val, row, idx) {
                let $btn = $(e.currentTarget);
                console.log(   $btn.data("field") );

                let tags = '<button type="button" class="btn btn-sm btn-outline-default mr-1">Backup ID: ' + row.id + '</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default mr-1">' + row.totalCount + ' files</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default mr-1">' + bytesToSize(row.totalSize) + '</button>';
                tags += '<button type="button" class="btn btn-sm btn-outline-default">' + row.srcDir + '</button>';

                $("#modal-backup-changes .modal-title .summary").html(tags);

                showChangedFiles(row.id, $btn.data("field"));
            },
        };

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
