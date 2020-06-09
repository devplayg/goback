package goback

import (
	"net/http"
)

type SignIn struct {
	AccessKey string
	SecretKey string
}

func errorPageTpl() string {
	return `
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

{{define "css"}}
    <link rel="stylesheet" media="screen, print" href="/assets/css/custom.css">
{{end}}


{{define "content"}}
	<div class="h-alt-hf d-flex flex-column align-items-center justify-content-center text-center">
		<h1 class="page-error color-fusion-500">
			<a href="/login">
				ERROR <span class="text-gradient">{{.StatusCode}}</span>
				<small class="fw-500">
					{{.Status}}
				</small>
			</a>
		</h1>
	</div>
{{end}}
`
}

func isLogged(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, SignInSessionId)
	if err != nil {
		return false
	}

	if len(session.Values) < 1 {
		return false
	}

	return true
}
