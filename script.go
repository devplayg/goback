package goback

import "strings"

func customCss() string {
	css := `
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

        .text-soft {
			color: #aaaaaa	 !important;
		}

		.bred {
			border: 1px solid red;
		}

        @media print{@page {size: landscape}}
`
	return strings.TrimSpace(css)
}

func customScript() string {
	script := `
        let waitMeOptions = {
			effect : "bounce",
			text : 'Loading..',
			bg : "rgba(255,255,255,0.7)",
			color : "#616469"
		};

		function thousandCommaSep(n) {
            return n.toLocaleString();
        }

		function bytesToSize(bytes) {
			return humanizedSize(bytes, true, 1);
		}

		function thousandCommaSep(n) {
			return n.toLocaleString();
		}

		function getSrcDirs(summaries) {
			let dir = {};
			$.each(summaries, function(i, r) {
				dir[r.srcDir] = true;
			});
			return Object.keys(dir).sort();
		}

		function basename(path) {
            return path.replace(/^.*[\\\/]/, '');
        }

		function getRate(i, total) {
            let per =  Math.round((1 - (i / total)) * 100);
            per = (per - (per % 10) - 20) * 10;
            if (per < 100) {
                per = 50;
            }
            return per;
            ss
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


		// Formatters
		function dateFormatter(val, row, idx) {
			return moment(val).format();
		}
		
		function thCommaFormatter(val, row, idx) {
			if (val === 0) {
				return '<span class="text-muted">0</span>';
			}
			return thousandCommaSep(val);
		}
		
		function shortDirFormatter(val, row, idx) {
			if (val.length < 21) {
				return val;
			}
			let dir = basename(val);
			return '<span class="has-tooltip" title="' + val + '">.. ' + dir + '</span>';
		}
		
		function byteSizeFormatter(val, row, idx, field) {
			if (val === 0) {
				return '<span class="text-muted ">' + bytesToSize(val) + '</span>';
			}

			let textCss = (field.endsWith("Failed")) ? "text-danger" : "";
			//if (val < 1000) {
			//	return bytesToSize(val);
			//}
			return '<span class="has-tooltip ' + textCss + '" title="' + val.toLocaleString() + ' Bytes">' + bytesToSize(val) + '</span>';
		}
		
		function yyyymmFormatter(val, row, idx) {
			return moment(val).format("ll");
		}

		function extFormatter(val, row, idx) {
            if (val.length < 1) {
                return;
            }
            return val;
        }

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

        function backupResultFormatter(val, row, idx, field) {
            if (val === 0) {
                return '<span class="text-muted">' + val + '</span>';
            }
            let textCss = (field.endsWith("Failed")) ? "text-danger" : "";

            let link = '<a class="changed ' + textCss + '" href="javascript:void(0);" title="changed" data-field="' + field + '">' + val.toLocaleString() + '</a>';
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
			
			let stats = ";"
			let prefix = '<span class="text-danger" title="Started-&gt;Read-&gt;Compared-&gt;Copied-&gt;Logged">',
				suffix = '</span>';
			if (val === 4) {
				stats = 'Copied';
			} else if (val === 3) {
				stats = 'Compared';
			} else if (val === 2) {
				stats = 'Read';
			} else if (val === 1) {
				stats = 'Started';
			} else {
				stats = val + '<i class="fa fa-warn"></i>';
			}
		
            return prefix + stats + suffix;
        }

        function toFixedFormatter(val, row, idx) {
            return  val.toFixed(2);
        }

        function thCommaFormatter(val, row, idx) {
            if (val === 0) {
                return '<span class="text-muted">0</span>';
            }
            return thousandCommaSep(val);
        }
`

	return strings.TrimSpace(script)
}
