package goback

func displayStats() string {
	return `
{{define "css"}}
    <style>
        .pagination .page-link {
            border-width: 1px;
        }

        .badge-stats {
            margin-right: 3px;
            font-weight: 400;
        }

        .bg-gray-light {
            background-color: #f7f7f7;
        }

        @media print{@page {size: landscape}}
    </style>
{{end}}

{{define "sidebar"}}
    <li>
        <a href="/backup/" title="Backup"><i class="fal fa-list-ul"></i><span class="nav-link-text">Backup logs</span></a>
    </li>
    <li>
        <a href="/stats/" title="Statistics"><i class="fal fa-chart-bar"></i><span class="nav-link-text">Statistics</span></a>
    </li>
    <li>
        <a href="/settings/" title="Settings"><i class="fal fa-cog"></i><span class="nav-link-text">Settings</span></a>
    </li>
{{end}}

{{define "content"}}
    <div class="subheader">
        <h1 class="subheader-title">
            <i class="fal fa-chart-bar mr-1"></i> Statistics
            <small>
                Weekly/Monthly statistics
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
                        data-buttons-class="default"
                        class="table table-sm"
                        data-toggle="table"
						data-cookie="true"
						data-cookie-id-table="stats"
                        data-toolbar="#toolbar-stats"
                        data-search="true"
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
                        <th data-field="date" data-sortable="true">Date</th>
                        <th data-field="srcDir" data-sortable="true" data-formatter="shortDirFormatter">Source</th>

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
              sysTz = "Asia/Seoul";
        // console.log(moment);
        // var moment = require('moment-timezone');
        // moment().tz("America/Los_Angeles").format();
        // import moment from 'moment';
        // import 'moment-timezone';
        // console.log(moment.tz);
        moment.tz.setDefault(sysTz);

        $.ajax({
            url: "/summaries",
        }).done(function(summaries) {
            // console.log(summaries);

            let stats = {};
            if ( $("#select-period").val() === Monthly) {
                $.each(summaries, function(i, r) {
                    let key = moment(r.date).format("YYYY-MM");
                    // console.log(key);
                    if (!(key in stats)) {
                        stats[key] = {};
                        // console.log("new! " + key);
                    }

                    if (!(r.srcDir in stats[key])) {
                        stats[key][r.srcDir] = {
                            date: key,
                            srcDir: r.srcDir,
                            countAdded: 0, sizeAdded: 0,
                            countModified: 0, sizeModified: 0,
                            countDeleted: 0, sizeDeleted: 0,
                            countFailed: 0, sizeFailed: 0,
                            countSuccess: 0, sizeSuccess: 0,
                        };

                    }

                    stats[key][r.srcDir].countAdded    += r.countAdded;
                    stats[key][r.srcDir].sizeAdded     += r.sizeAdded;
                    stats[key][r.srcDir].countModified += r.countModified;
                    stats[key][r.srcDir].sizeModified  += r.sizeModified;
                    stats[key][r.srcDir].countDeleted  += r.countDeleted;
                    stats[key][r.srcDir].sizeDeleted   += r.sizeDeleted;
                    stats[key][r.srcDir].countFailed   += r.countFailed;
                    stats[key][r.srcDir].sizeFailed    += r.sizeFailed;
                    stats[key][r.srcDir].countSuccess  += r.countSuccess;
                    stats[key][r.srcDir].sizeSuccess   += r.sizeSuccess;

                    console.log(r.date + " => " + key );
                });
            }

            let rows = [];
            $.each(stats, function(date, statsByDir) {
                // console.log(dir);
                $.each(statsByDir, function(dir, data) {
                    rows.push(data);
                });
            });
            console.log(rows);

            // console.log( Object.values(stats));

            // Initialize directory selector
            initDirSelector(getSrcDirs(summaries));

            // Update table
            $("#table-stats").bootstrapTable("load", rows);
        });

        console.log(moment.tz);

        $( "#select-srcDirs" ).change(function () {
            if (this.value.length > 0) {
                $("#table-stats").bootstrapTable('filterBy', {
                    srcDir: [this.value]
                });
                return;
            }

            $("#table-backup").bootstrapTable('filterBy',{});
        });

        function initDirSelector(dirs) {
            $("#select-srcDirs").empty();
            $("<option/>", {
                "value": "",
            }).text("Directories").appendTo($("#select-srcDirs"));
            $.each(dirs, function(i, dir) {
                console.log(dir);
                $("<option/>", {
                    "value": dir
                }).text(dir).appendTo($("#select-srcDirs"));
            });
        }

        $(".table")
            .on('all.bs.table', function (e, data) {
                $('.has-tooltip').tooltip();
            })

    </script>
{{end}}
`
}