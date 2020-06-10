package tpl

func Settings() string {
	return `{{define "css"}}
    <link rel="stylesheet" media="screen, print" href="/assets/css/custom.css">
{{end}}

{{define "sidebar"}}
    <li>
        <a href="/backup/" title="Backup"><i class="fal fa-list-ul"></i><span class="nav-link-text">Backup logs</span></a>
    </li>
    <li>
        <a href="/stats/" title="Statistics"><i class="fal fa-chart-bar"></i><span class="nav-link-text">Statistics</span></a>
    </li>
    <li class="active">
        <a href="/settings/" title="Settings"><i class="fal fa-cog"></i><span class="nav-link-text">Settings</span></a>
    </li>
    <li>
        <a href="/logout" title="Sign out"><i class="fal fa-sign-out"></i><span class="nav-link-text">Sign out</span></a>
    </li>
{{end}}

{{define "content"}}
    <div class="subheader">
        <h1 class="subheader-title">
            <i class="subheader-icon fal fa-cog"></i> Settings
            <small>
                System configuration
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
    <div class="row">
        <div class="col-md-4">
            {{range $i, $storage := .Settings.Storages}}
                <form class="form" data-id="{{$storage.Id}}" data-type="storage">
                    <input type="hidden" name="id" value="{{$storage.Id}}"/>
                    <input type="hidden" name="checksum" value="{{$.Checksum}}" />

                    <div class="card border mb-4">
                        <div class="card-header d-flex pr-2 align-items-center flex-wrap">
                            <div class="card-title">
                                {{if eq $storage.Id 1}}
                                    <i class="fal fa-hdd"></i> Local Disk Storage
                                {{else if eq $storage.Id 2}}
                                    SFTP
                                {{end}}
                            </div>
                        </div>
                        <div class="card-body">
                            {{if eq $storage.Protocol 1}}
                                <div class="form-group">
                                    <label class="form-label">Local Directory Path</label>
                                    <input type="text" name="dir" class="form-control {{if DirExists $storage.Dir  }}is-valid{{else}}is-invalid{{end}}"  value="{{$storage.Dir}}" />
                                    {{if eq (DirExists $storage.Dir) false }}
                                        <div class="invalid-feedback">
                                            Directory does not exist
                                        </div>
                                    {{end}}
                                </div>

                            {{else if eq $storage.Protocol 4}}
                                <div class="row">
                                    <div class="col-lg-8">
                                        <div class="form-group">
                                            <label class="form-label">Hostname</label>
                                            <input type="text" name="host" class="form-control" value="{{$storage.Host}}">
                                        </div>
                                    </div>
                                    <div class="col-lg-4">
                                        <div class="form-group">
                                            <label class="form-label">Port</label>
                                            <input type="text" name="port" class="form-control" value="{{$storage.Port}}">
                                        </div>
                                    </div>
                                </div>

                                <div class="row mt-3">
                                    <div class="col-lg-6">
                                        <div class="form-group">
                                            <label class="form-label">Hostname</label>
                                            <input type="text" name="username" class="form-control" value="{{$storage.Username}}">
                                        </div>
                                    </div>
                                    <div class="col-lg-6">
                                        <div class="form-group">
                                            <label class="form-label">Password</label>
                                            <input type="text" name="password" class="form-control" value="{{$storage.Password}}">
                                        </div>
                                    </div>
                                </div>

                                <div class="form-group mt-3">
                                    <label class="form-label">Directory</label>
                                    <input type="text" name="dir" class="form-control" value="{{$storage.Dir}}">
                                </div>
                            {{end}}
                        </div>
                        <div class="card-body">
                            <button class="btn btn-primary">Apply</button>
                        </div>
                    </div>
                </form>
            {{end}}
        </div>
        <div class="col-md-8">
            {{range $idx, $job := .Settings.Jobs}}
                <form class="form" data-id="{{$job.Id}}" data-type="job">
                    <input type="hidden" name="id" value="{{$job.Id}}"/>
                    <input type="hidden" name="backupType" value="{{$job.BackupType}}"/>
                    <input type="hidden" name="checksum" value="{{$.Checksum}}" />

                    <div class="card border mb-4">
                        <div class="card-header d-flex pr-2 align-items-center flex-wrap">
                            <div class="card-title">#{{$job.Id}} </div>
                            <div class="custom-control d-flex custom-switch ml-auto">
                                <input id="job-checkbox-{{$job.Id}}" type="checkbox" name="enabled" class="custom-control-input" {{if $job.Enabled}}checked{{end}}>
                                <label class="custom-control-label fw-500" for="job-checkbox-{{$job.Id}}"></label>
                            </div>
                        </div>
                        <div class="card-body">
                            <div class="row">
                                <div class="col-5">
                                    <div class="form-group">
                                        <label class="form-label">Storage</label>
                                        <select class="form-control" name="storageId" data-style="btn-default" readonly>
                                            <option value="1" {{if eq $job.StorageId 1}}selected{{end}}>Local Disk</option>
                                            <!--{*<option value="2" {{if eq $job.StorageId 2}}selected{{end}}>SFTP</option>*}-->
                                        </select>

                                    </div>
                                </div>
                                <div class="col-7">
                                    <div class="form-group">
                                        <label class="form-label">Schedule</label>
                                        <input type="text" name="schedule" class="form-control"  value="{{$job.Schedule}}">
                                    </div>
                                </div>
                            </div>

                            <div class="directory">
                                <p class="form-label mt-3 mb-2">
                                    Directories to backup
                                    <a href="javascript:void(0);" class="btn btn-primary btn-xs btn-icon createDir">
                                        <i class="fal fa-plus"></i>
                                    </a>
                                </p>

                                <div class="dirs">
                                    {{range $j, $dirName := $job.SrcDirs }}
                                        <div class="form-group mb-2 input-dir">
                                            <div class="row">
                                                <div class="col-10">
                                                    <input type="text" name="srcDirs"
                                                           class="form-control {{if DirExists $dirName  }}is-valid{{else}}is-invalid{{end}}"
                                                           value="{{$dirName}}" />
                                                    {{if eq (DirExists $dirName) false }}
                                                        <div class="invalid-feedback">
                                                            Directory does not exist
                                                        </div>
                                                    {{end}}
                                                </div>
                                                <div class="col-2">
                                                    <a href="javascript:void(0);" class="removeDir"><i class="fal fa-times"></i></a>
                                                </div>
                                            </div>
                                        </div>
                                    {{end}}
                                    <div class="form-group mb-2 input-dir">
                                        <div class="row">
                                            <div class="col-10">
                                                <input type="text" name="srcDirs" class="form-control" />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div>
                                <button type="button" class="btn btn-danger runBackup d-none" data-id="{{$job.Id}}">Run backup</button>
                                <button class="btn btn-primary">Apply</button>
                            </div>
                        </div>
                    </div>
                </form>
            {{end}}

        </div>
    </div>123
{{end}}

{{define "script"}}
    <script src="/assets/js/custom.js"></script>
    <script>
        $(".form").each(function() {
            let c = $(this);
            $(this).validate({
                submitHandler: function (form, e) {
                    e.preventDefault();
                    let $form = $(this.currentForm),
                        url = "/settings/" + $form.data("type") + "/id/"+$form.data("id");
                    $.ajax({
                        url: url,
                        type: "PATCH",
                        data: $form.serialize(),
                    }).done(function (data) {
                        window.location.reload();
                    }).fail(function (jqXHR, textStatus, errorThrown) {
                        console.error(jqXHR);
                        Swal.fire("Error", jqXHR.responseJSON.error, "warning");
                    }).always(function() {
                    });
                },
            });
        });

        $(document).on("click", ".removeDir" , function(e) {
            e.preventDefault();
            $(this).closest(".input-dir").remove();
        });
        $(".createDir").click(function(e) {
            e.preventDefault();
            let elm = '<div class="form-group mb-2 input-dir"><div class="row"><div class="col-10"><input type="text" name="srcDirs" class="form-control"/></div><div class="col-2"><a href="javascript:void(0);" class="removeDir"><i class="fal fa-times"></i></a></div></div></div>';
            $(this).closest(".directory").find(".dirs").append($(elm));
        });

        $(".runBackup").click(function(e) {
            let id = $(this).data("id");

            $.ajax({
                url: "/backup/" + id + "/run",
                type: "GET",
            }).done(function(data) {
                // console.log(data);
            }).fail(function (jqXHR, textStatus, errorThrown) {
                console.error(jqXHR);
                Swal.fire("Error", jqXHR.responseJSON.error, "warning");
            }).always(function() {
            });
        });

        $(".custom-control-input").change(function() {
            let $card = $(this).closest(".card");
            $("input[type=text]", $card).prop("readonly", !this.checked);
            $(".custom-control-label", $card).text((this.checked) ? "Enabled" : "Disabled");
        });

        $(".custom-control-input").each(function(i, $c) {
            let $card = $c.closest(".card");
            $("input[type=text]", $card).prop("readonly", !this.checked);
            $(".custom-control-label", $card).text((this.checked) ? "Enabled" : "Disabled");
        });

    </script>
{{end}}
`
}
