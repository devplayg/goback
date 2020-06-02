package goback

import "strings"

func customScript() string {
	script := `
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
		
		function thCommaFormatter(val, row, idx) {
			if (val === 0) {
				return '<span class="text-muted">0</span>';
			}
			return thousandCommaSep(val);
		}
		
		function thousandCommaSep(n) {
			return n.toLocaleString();
		}
		
		function shortDirFormatter(val, row, idx) {
			if (val.length < 21) {
				return val;
			}
			let dir = basename(val);
			return '<span class="has-tooltip" title="' + val + '">.. ' + dir + '</span>';
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
		
		function getSrcDirs(summaries) {
			let dir = {};
			$.each(summaries, function(i, r) {
				dir[r.srcDir] = true;
			});
			return Object.keys(dir).sort();
		}

`

	return strings.TrimSpace(script)
	// b, _ := compress.Compress([]byte(script), compress.GZIP)
	// return base64.StdEncoding.EncodeToString(b)
}
