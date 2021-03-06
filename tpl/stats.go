package tpl

import "os"

func Stats() string {
	tpl := "tpl/stats.html"
	if _, err := os.Stat(tpl); os.IsExist(err) {
		return displayWithLocalFile(tpl)
	}

	return `{{define "css"}}
    <link rel="stylesheet" media="screen, print" href="/assets/css/custom.css">
{{end}}

{{define "sidebar"}}
    <li>
        <a href="/backup/" title="Backup"><i class="fal fa-list-ul"></i><span class="nav-link-text">Backup logs</span></a>
    </li>
    <li class="active">
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
            <i class="fal fa-chart-bar mr-1"></i> Statistics
            <small>
                Monthly statistics
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
                <div id="toolbar-stats">
                    <form class="form-inline">
                        <select id="select-srcDirs" class="form-control"></select>
                        <select id="select-period" class="form-control ml-2 d-none">
                            <option value="weekly">Weekly</option>
                            <option value="monthly" selected>Monthly</option>
                        </select>
                    </form>
                    </select>
                </div>
                <table  id="table-stats"
                        class="table table-sm"
                        data-buttons-class="default"
                        data-url="/stats"
                        data-toggle="table"
                        data-cookie-id-table="stats"
                        data-cookie="true"
                        data-cookies-enabled="['bs.table.columns']"
                        data-toolbar="#toolbar-stats"
                        data-search="true"
                        data-side-pagination="client"
                        data-pagination="true"
                        data-show-export="true"
                        data-export-types="['csv', 'txt', 'excel']"
                        data-show-refresh="true"
                        data-show-columns="true"
                        data-sort-name="date"
                        data-page-size="20"
                        data-sort-order="desc">
                    <thead>
                    <tr>
                        <th data-field="date" data-sortable="true" data-formatter="yyyymmFormatter">Date</th>
                        <th data-field="report" data-sortable="true" data-formatter="reportFormatter">Report</th>
                        <th data-field="srcDir" data-sortable="true" data-formatter="shortDirFormatter">Directory</th>

                        <th data-field="countAdded" data-sortable="true" data-formatter="thousandCommaSep">Added</th>
                        <th data-field="sizeAdded" data-sortable="true" data-formatter="byteSizeFormatter">Added</th>

                        <th data-field="countModified" data-sortable="true" data-formatter="thousandCommaSep">Modified</th>
                        <th data-field="sizeModified" data-sortable="true" data-formatter="byteSizeFormatter">Modified</th>

                        <th data-field="countDeleted" data-sortable="true" data-formatter="thousandCommaSep">Deleted</th>
                        <th data-field="sizeDeleted" data-sortable="true" data-formatter="byteSizeFormatter">Deleted</th>

                        <th data-field="countFailed" data-sortable="true" data-formatter="thousandCommaSep" class="bg-gray-light">Failed</th>
                        <th data-field="sizeFailed" data-sortable="true" data-formatter="byteSizeFormatter" class="bg-gray-light">Failed</th>

                        <th data-field="countSuccess" data-sortable="true" data-formatter="thousandCommaSep" class="bg-gray-light">Backup Files</th>
                        <th data-field="sizeSuccess" data-sortable="true" data-formatter="byteSizeFormatter" class="bg-gray-light">Backup Size</th>
                    </tr>
                    </thead>
                </table>
            </div>
        </div>
    </div>
{{end}}

{{define "script"}}
    <script src="/assets/plugins/bootstrap-select/bootstrap-select.min.js"></script>
    <script src="/assets/js/custom.js"></script>
    <script>
        const Monthly = "monthly",
            Weekly = "weekly",
            SystemTz = "Asia/Seoul";
        moment.tz.setDefault(SystemTz);
        $( "#select-srcDirs" ).change(function () {
            if (this.value.length > 0) {
                $("#table-stats").bootstrapTable('filterBy', {
                    srcDir: [this.value]
                });
                return;
            }

            $("#table-stats").bootstrapTable('filterBy',{});
        });

        $(".table")
            .on('all.bs.table', function (e, data) {
                $('.has-tooltip').tooltip();
            })
            .on('load-success.bs.table', function(e, rows) {
                let dirMap = {};
                $.each(rows, function(i, s) {
                    dirMap[s.srcDir] = true;
                });
                Object.keys(dirMap).sort();

                $("#select-srcDirs").empty();
                $("<option/>", {
                    "value": "",
                }).text("Directories").appendTo($("#select-srcDirs"));
                $.each(dirMap, function(dir, r) {
                    $("<option/>", {
                        "value": dir
                    }).text(dir).appendTo($("#select-srcDirs"));
                });
            });

        function reportFormatter(val, row, idx) {
            return '<a class="btn btn-primary btn-xs" href="/report/' + moment(row.date).format("YYYYMM") + '/' + '" target="_blank"><i class="fal fa-file-alt"></i> Report</a>';
        }
    </script>
{{end}}
`
}
