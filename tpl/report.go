package tpl

import (
	"os"
)

func Report() string {
	tpl := "tpl/report.html"
	if _, err := os.Stat(tpl); os.IsExist(err) {
		return displayWithLocalFile(tpl)
	}

	return `{{define "css"}}
    <link rel="stylesheet" media="screen, print" href="/assets/css/custom.css">
    <link rel="stylesheet" media="screen, print" href="/assets/css/page-invoice.css">
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
    <li>
        <a href="/logout" title="Sign out"><i class="fal fa-sign-out"></i><span class="nav-link-text">Sign out</span></a>
    </li>
{{end}}

{{define "content"}}
    <div class="container">
        <div data-size="A4" class="bg-white">
            <h2 class="keep-print-font fw-500 mb-0 text-primary flex-1 position-relative">
                <img src="/img/logo.png"/> Monthly Backup Report
            </h2>

            <div class="row">
                <div class="col-sm-6">
                    <h3 class="display-4 fw-500 color-primary-600 keep-print-font pt-4 l-h-n m-0">
                        May, 2020
                    </h3>
                    <div class="text-dark fw-400 h1 mb-g keep-print-font">
                        Jan 1 ~ Jan 31, 2010
                    </div>
                </div>
                <div class="col-sm-6">

                    <div class="row mt-5 color-fusion-900">
                        <div class="col-sm-4 offset-sm-3"><h2 class="fw-500 color-primary-600">Summary</h2></div>
                    </div>
                    <div class="row mt-2 color-fusion-900">
                        <div class="col-sm-4 offset-sm-3">Added</div>
                        <div class="col-sm-5 text-right">
                            <span class="countAdded"></span>
                            (<span class="sizeAdded"></span>)
                        </div>
                    </div>
                    <div class="row mt-2">
                        <div class="col-sm-4 offset-sm-3">Modified</div>
                        <div class="col-sm-5 text-right">
                            <span class="countModified"></span>
                            (<span class="countDeleted"></span>)
                        </div>
                    </div>
                    <div class="row mt-2">
                        <div class="col-sm-4 offset-sm-3">Deleted</div>
                        <div class="col-sm-5 text-right">
                            <span class="countDeleted"></span>
                            (<span class="sizeDeleted"></span>)
                        </div>
                    </div>
                    <div class="row mt-2">
                        <div class="col-sm-4 offset-sm-3">Backup</div>
                        <div class="col-sm-5 text-right">
                            <span class="countSuccess"></span>
                            (<span class="sizeSuccess"></span>)
                        </div>
                    </div>
                </div>
            </div>
            <div class="row mt-5">
                <div class="col-sm-12">
                    <table  id="table-report"
                            data-buttons-class="default"
                            class="table table-sm"
                            data-toggle="table"
                            data-cookie="true"
                            data-cookie-id-table="report"
                            data-cookies-enabled="['bs.table.columns']"
                            data-toolbar="#toolbar-report"
                            data-export-types="['csv', 'txt', 'excel']"
                            data-side-pagination="client"
                            data-sort-name="id"
                            data-sort-order="desc">
                        <thead>
                        <tr>
                            <th data-field="date" data-sortable="true" data-formatter="yyyymmddFormatter">Date</th>
                            <th data-field="workerCount" data-sortable="true" data-visible="false">Workers</th>
                            <th data-field="countAdded" data-sortable="true" data-formatter="thousandCommaSep">Added</th>
                            <th data-field="countModified" data-sortable="true" data-formatter="thousandCommaSep">Modified</th>
                            <th data-field="countDeleted" data-sortable="true" data-formatter="thousandCommaSep">Deleted</th>
                            <th data-field="sizeSuccess" data-sortable="true" data-formatter="byteSizeFormatter">Backup Size</th>
                        </tr>
                        </thead>
                    </table>
                </div>
            </div>
        </div>
    </div>
{{end}}

{{define "script"}}
    <script src="/assets/js/custom.js"></script>
    <script>
        const url = new URL(window.location.href);
        $.ajax({
            url: url.pathname.substring(0, url.pathname.length-1),
        }).done(function(report) {
            $("#table-report").bootstrapTable("load", report);
            let countAdded = 0,
                countModified = 0,
                countDeleted = 0,
                countSuccess = 0,
                sizeAdded = 0,
                sizeModified = 0,
                sizeDeleted = 0,
                sizeSuccess = 0;

            $.each(report, function(i, r) {
                countAdded += r.countAdded;
                countModified += r.countModified;
                countDeleted += r.countDeleted;
                countSuccess += r.countSuccess;
                sizeAdded += r.sizeAdded;
                sizeModified += r.sizeModified;
                sizeDeleted += r.sizeDeleted;
                sizeSuccess += r.sizeSuccess;
            });

            $(".countAdded").text(thousandCommaSep(countAdded));
            $(".countModified").text(thousandCommaSep(countModified));
            $(".countDeleted").text(thousandCommaSep(countDeleted));
            $(".countSuccess").text(thousandCommaSep(countSuccess));
            $(".sizeAdded").text(bytesToSize(sizeAdded));
            $(".sizeModified").text(bytesToSize(sizeModified));
            $(".sizeDeleted").text(bytesToSize(sizeDeleted));
            $(".sizeSuccess").text(bytesToSize(sizeSuccess));
        });
    </script>
{{end}}
`
}
