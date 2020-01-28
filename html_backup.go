package goback

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func DisplayBackupTest() string {
	b, err := ioutil.ReadFile("static/backup.html")
	if err != nil {
		log.Error(err)
	}
	return string(b)
}

func DisplayBackup() string {
	return `
    {{define "content"}}
    <style>
        .pagination .page-link {
            border-width: 1px;
        }

        .badge-stats {
            margin-right: 3px;
        }
    </style>
    <div class="row">
        <div class="col">
            <div  class="panel ">
                <div class="panel-container show">
                    <div class="panel-content">
                        <div id="toolbar-backup">
                        </div>
                        <table  id="table-backup"
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
                                data-show-toggle="true"
                                data-page-size="20"
                                data-sort-order="desc">
                            <thead>
                            <tr>
                                <th data-field="id" data-sortable="true">ID</th>
                                <th data-field="date" data-sortable="true" data-formatter="dateFormatter">Date</th>
                                <th data-field="srcDir" data-sortable="true">Backup Dir.</th>
                                <th data-field="backupType" data-sortable="true" data-formatter="backupTypeFormatter">Type</th>
                                <th data-field="state" data-sortable="true" data-formatter="backupStateFormatter">State</th>
                                <th data-field="execTime" data-sortable="true" data-formatter="toFixedFormatter">Time(Sec)</th>

                                <th data-field="workerCount" data-sortable="true" data-visible="false">Workers</th>

                                <th data-field="totalCount" data-sortable="true" data-formatter="thCommaFormatter">Files</th>
                                <th data-field="totalSize" data-sortable="true" data-formatter="bytesToSize">Total Size</th>
                                <th data-field="totalSize" data-sortable="true" data-formatter="thCommaFormatter" data-visible="false">Total Size (B)</th>


                                <th data-field="countAdded" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Added</th>
                                <th data-field="sizeAdded" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents" data-visible="false">Added (B)</th>
                                <th data-field="_sizeAdded" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Added</th>

                                <th data-field="countModified" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Modified</th>
                                <th data-field="sizeModified" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents" data-visible="false">Modified (B)</th>
                                <th data-field="_sizeModified" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Modified</th>

                                <th data-field="countDeleted" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Deleted</th>
                                <th data-field="sizeDeleted" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents" data-visible="false">Deleted (B)</th>
                                <th data-field="_sizeDeleted" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Deleted</th>

                                <th data-field="countFailed" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents">Failed</th>
                                <th data-field="sizeFailed" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents" data-visible="false">Failed (B)</th>
                                <th data-field="_sizeFailed" data-sortable="true" data-formatter="backupResultFormatter" data-events="backupOperateEvents" data-visible="true">Failed</th>

                                <th data-field="message" data-sortable="true">Msg(Read/Compare/Copy/Write)</th>
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

        </div>
    </div>


    <div class="modal fade" id="modal-backup-changes" tabindex="-1" role="dialog" aria-hidden="true">
        <div class="modal-dialog  mw-100 w-75" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h2 class="modal-title">Backup logs</h2>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true"><i class="fal fa-times"></i></span>
                    </button>
                </div>
                <div class="modal-body">
                    <ul id="tabs-backup-changes" class="nav nav-tabs" role="tablist">
                        <li class="nav-item">
                            <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-added" role="tab"><i class="far fa-plus text-info"></i>
                                Added
                                <span class="badge border border-info rounded-pill bg-primary-500 stats-data stats-added"></span>
                            </a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-modified" role="tab"><i class="far fa-pen text-success"></i>
                                Modified
                                <span class="badge border border-info rounded-pill bg-success-500 stats-data stats-modified"></span>
                            </a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-deleted" role="tab"><i class="far fa-minus text-warning"></i>
                                Deleted
                                <span class="badge border border-warning rounded-pill bg-warning-500 stats-data stats-deleted"></span>
                            </a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link fs-lg" data-toggle="tab" href="#tab-backup-failed" role="tab"><i class="far fa-times text-danger"></i>
                                Failed
                                <span class="badge border border-danger rounded-pill bg-danger-500 stats-data stats-failed"></span>
                            </a>
                        </li>
                    </ul>
                    <div class="tab-content p-3">
                        <div class="tab-pane fade show active" id="tab-backup-added" role="tabpanel">
                            <div id="toolbar-backup-added">
                                <h3>
                                    <div class="stats-data stats-added-size"></div>
                                </h3>
                                <div class="stats-data stats-added-ext"></div>
                            </div>
                            <table  id="table-backup-added"
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
                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                    <th data-field="mtime" data-visible="false">ModTime</th>
                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                </tr>
                                </thead>
                            </table>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-modified" role="tabpanel">
                            <div id="toolbar-backup-modified">
                                <h3>
                                    <div class="stats-data stats-modified-size"></div>
                                </h3>
                                <div class="stats-data stats-modified-ext"></div>
                            </div>
                            <table  id="table-backup-modified"
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
                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                    <th data-field="mtime" data-visible="false">ModTime</th>
                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                </tr>
                                </thead>
                            </table>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-deleted" role="tabpanel">
                            <div id="toolbar-backup-deleted">
                                <h3>
                                    <div class="stats-data stats-deleted-size"></div>
                                </h3>
                                <div class="stats-data stats-deleted-ext"></div>
                            </div>
                            <table  id="table-backup-deleted"
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
                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                    <th data-field="mtime" data-visible="false">ModTime</th>
                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                </tr>
                                </thead>
                            </table>
                        </div>
                        <div class="tab-pane fade" id="tab-backup-failed" role="tabpanel">
                            <div id="toolbar-backup-failed">
                                <h3>
                                    <div class="stats-data stats-failed-size"></div>
                                </h3>
                                <div class="stats-data stats-failed-ext"></div>
                            </div>
                            <table  id="table-backup-failed"
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
                                    <th data-field="ext" data-sortable="true" data-formatter="extFormatter">Extension</th>
                                    <th data-field="name" data-sortable="true" data-formatter="backupChangesNameFormatter">Name</th>
                                    <th data-field="dir" data-sortable="true" data-visible="false">Directory</th>
                                    <th data-field="mtime" data-visible="false">ModTime</th>
                                    <th data-field="size" data-sortable="true" data-formatter="byteSizeFormatter">Size</th>
                                    <th data-field="sizeB" data-sortable="true" data-visible="false" data-formatter="sizeBFormatter">Size (B)</th>
                                    <th data-field="state" data-sortable="true" data-formatter="backupFileStateFormatter">Backup</th>
                                </tr>
                                </thead>
                            </table>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                </div>
            </div>
        </div>
    </div>
{{end}}

{{define "script"}}
    <script>

        $('#modal-backup-changes')
            .on('show.bs.modal', function (e) {
            }).on('hidden.bs.modal', function (e) {
                $(".table-data").bootstrapTable("removeAll");
                $(".stats-data").empty();
            });

        $(".table").on('all.bs.table', function (e, data) {
            $('.has-tooltip').tooltip();
        });

        window.backupOperateEvents = {
            'click .file': function (e, val, row, idx) {
                console.log(row);
                let $btn = $(e.currentTarget);
                // $("#modal-files .modal-title").text($btn.data("title"));

                // console.log($btn.data("field"));
                showChangedFiles(row.id, $btn.data("title"), $btn.data("field"));
                // getBackupFiles(row, $btn.data("field"));
                // $("#modal-files").modal("show");
            },
        };


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
            return humanizedSize(bytes, true);
        }

        function humanizedSize(bytes, si) {
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
            return bytes.toFixed(1)+' '+units[u];
        }

        // function bytesWithSize(bytes, prefix, suffix) {
        //     if (prefix === undefined) {
        //         prefix = "";
        //     }
        //     if (suffix === undefined) {
        //         suffix = "";
        //     }
        //     console.log("==================");
        //     console.log(bytes + " B");
        //     console.log(bytesToSize(bytes));
        //     // console.log((bytes + " B") === bytesToSize(bytes));
        //     if ((bytes + " B") == bytesToSize(bytes)) {
        //         return bytesToSize(bytes);
        //     }
        //     return bytesToSize(bytes) + prefix + bytes.toLocaleString() + " B" + suffix;
        // }

        // function basename(path) {
        //     return path.replace(/^.*[\\\/]/, '');
        // }
        //
        // function dirname(path) {
        //     // return path.substr(0, basename(path).lastIndexOf('.'));
        //
        //     return path.substring(0, path.lastIndexOf(basename(path)));
        //
        //     // return path.replacee
        //     // trimEnd(path, basename(path));
        // }

        function showChangedFiles(id, title, field) {
            let url = "/summaries/" + id + "/changes";
            $.ajax({
                url: url,
            }).done(function(report) {
                console.log(report);
                updateBackupStats("added", report.added, "primary");
                updateBackupStats("modified", report.modified, "success");
                updateBackupStats("deleted", report.deleted, "warning");
                updateBackupStats("failed", report.failed, "danger");

                $("#modal-backup-changes").modal("show");
                $('#tabs-backup-changes a[href="#tab-backup-' + getTab(field) + '"]').tab('show');
            });
        }

        function updateBackupStats(how, what, colorSuffix) {
            // console.log(what);
            $("#table-backup-" + how).bootstrapTable("load", what.files);
            $(".stats-" + how).text(what.files.length > 0 ? what.files.length : "");
            $(".stats-" + how + "-size").html(bytesToSize(what.size) + " <i>" + how + "</i>");

            let extTags = "",
                total = what.report.extRanking.length
            console.log("================");
            $.each(what.report.extRanking, function(i, r) {
                if (r.ext.length < 1) {
                    r.ext = "-";
                }
                extTags += '<span class="badge badge-stats bg-' + colorSuffix + '-' + getPer(i, total) + '">'+ r.ext + " / " + bytesToSize(r.size) + '</span>';
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
            let str = humanizedSize(val, true);
            if (val < 1000) {
                return '<span class="text-muted">' + str +'</span>';
            }

            if (val >= 10000000) {
                return '<span class="text-danger">' + str +'</span>';
            }

            return str;
            //
            // if (str.endsWith(" kB")) {
            //     return '<span class="">' + str +'</span>';
            // }
            //
            // if (str.endsWith(" MB")) {
            //     return '<span class="text-danger">' + str +'</span>';
            // }
        }

        function backupBasenameFormatter(val, row, idx) {
            return basename(val);
        }

        function backupDirFormatter(val, row, idx) {
            return '<span class="text-muted">' + dirname(row.path) + '</span>';
        }

        function backupResultFormatter(val, row, idx, field) {
            let humanized = false;
            let rowName = field;
            if (field.charAt(0) === "_") {
                rowName = field.substr(1);
                humanized = true;
            }

            let th = $('#table-backup').find("[data-field='" + rowName + "']");

            if (val === undefined) {
                val = row[rowName];
            }

            if (val === 0) {
                return '<span class="text-muted">' + val + '</span>';
            }

            let returnVal = val;
            if (humanized) {
                returnVal = humanizedSize(val, true);
            }

            let link = $("<a/>", {
                href: "#",
                class: "file",
                "data-title": th.text(),
                "data-field": rowName,
                "title": "",
            }).html(
                returnVal.toLocaleString()
            );
            return link[0].outerHTML;
        }

        function backupTypeFormatter(val, row, idx) {
            if (val === 1) {
                return 'Initial';
            }
            if (val === 2) {
                return 'Incremental';
            }
        }

        function backupExtFormatter(val, row, idx) {
            return '<a href="#" class="extension">Ext</a>';
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
            return val.toLocaleString();
        }

    </script>
{{end}}
`
}
