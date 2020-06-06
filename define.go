package goback

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

const (
	FileModified = 1
	FileAdded    = 2
	FileDeleted  = 4

	FileBackupFailed    = -1
	FileBackupSucceeded = 1

	FilesDbName   = "files-%s.db"
	ChangesDbName = "changes-%d-%s.db"

	Initial     = 1
	Incremental = 2
	Full        = 4

	LocalDisk = 1
	Ftp       = 2
	Sftp      = 4

	GobEncoding  = 1
	JsonEncoding = 2

	GZIP = "GZIP"

	// Content Type
	ApplicationJson = "application/json"

	SignInSessionId = "gbSessionId"

	AccessKey = "GOBACK_ACCESS_KEY"
	SecretKey = "GOBACK_SECRET_KEY"

	// URI
	AssetUriPrefix = "/assets/"
	LoginUri       = "/login"
	HomeUri        = "/backup/"

	LogoImg = "H4sIAAAAAAAA/wBCBL37iVBORw0KGgoAAAANSUhEUgAAAB4AAAAeCAYAAAA7MK6iAAAACXBIWXMAAC4jAAAuIwF4pT92AAAD9ElEQVR42sWX709bVRjHT6uwV/oXGN8bX/hOZyIZY5vNwgphGjRZoiTGaYz6Yrhh293e/liCLiRtX1mD24hJMcB6SwddnU5cLFMGZJslkQED2vurP6hi9iPbKJzH57kUZdMtBFrX5JOn9+acz/eep+fe3DKXy8UeB/86IVIVRaOKLvcTW+Mf1yOD7xuI1ev4lHmO2phHsGPdIMZYmzF3vUt85IrXBjgFo37QOfBWS+/lEy29V75u6dsgvQYnaK5znYvc/xlstFakAU6T4PZuOyBNRq0X7kLj0O1NcAusF+7BgfBvEcHtqRbRSW7xYSt24u/icbSx976KvVv/411oiufuNJ0rFDcFziXH+52D75CT3A9ZMbZZFPHKRNYsTf2879zickM8X2z4tsBLwAZZHYtzydEsXRsmp+F2ifcHUwucostEm6O1I/h8fbywZI3lVvadzXPr2TxsAm7MRQe5Wju+eI7c2G6T+OCKqRVebMnbJ3/wWuKLUD+g36sfzMKWQIcl/ge0nDzv9j7Q7tUV045zOpngPlZlDV+fsgxk+d4z+vLegQxHYJMYDnI1hGcmccM+SRliaXevtdnsxXvvI1/Xjt2DebD0K8uWqMYR2CKcXOT82HeqhjIoSzRWjFeA95uZbvo3uhKdO8/k+B5JXtoTUaAsoIucb5766UvKoCwXbrK/Hxg27+dPv9o3m90lyXxXWF7ZjRWBLWK46rCiO2PzfvaU8UDBe5oJLk81teBg4JvmmogOdX1zxbrTKSgr6CT3wUD365SFmVW006rczqOssWsk+spphdf2zBZre+c4AmXCcJK7seuXCGUJLnc1oys43O5/trbn+q0dPbO8BJQZw4sZN4+0+57xCA7GjtsPmfb7I4df6E7By6Gp4vbQDFQCclPGa37p0HFHq4l9YnMwMXZ15PvpLE9Ma8XEtM4TMzqUFXQOo5syhNjVS0fa2hiLxuIf5jP6nRuFHCzms/zPhSxUAnLfWMjyhVx2qScs7WcZXdc1TQdF00DFWjk0kFV1hbImksmLTMGPqqocgXW17OiqzjGK/Hx8fHyYybKsIhyBSqHICkzJ1yCdThs5Y2NjwwwPVDqBQLmZT81DJp2FofnzsHP6RfhuLs5zch5GRkeMYG0tmFZezuBUOgVqWoXk/K/QMdnOr8xd5pqs8dGx0YssmUymSoHFSoDhRSWtFAvp34tKSlmSFRkSicQQCwaDTRh+U6NdXdoI9L0cqNqqi5xpJWUcT0xM3A4EAjXMbrczv9+/XZKkYDQaDfX394ewdpeTkjMUDoeDPp/vJcqktwGzIAjM4XD8L1CW2+02rb2KmPGA3ocqSinDTJnscf1p+wtMni0TTGHjQwAAAABJRU5ErkJgggEAAP//5/mgNEIEAAA="
)

var (
	SummaryBucket     = []byte("summary")
	BackupBucket      = []byte("backup")
	ConfigBucket      = []byte("config")
	KeyConfig         = []byte("config")
	KeyConfigChecksum = []byte("config_checksum")
)

const (
	Started = iota + 1
	Read
	Compared
	Copied
	Logged
)

const (
	kB = 1000
	MB = 1000000
	GB = 1000000000
	TB = 1000000000000
)

var fileSizeCategories = []int64{
	0,

	5 * kB,
	50 * kB,
	500 * kB,

	5 * MB,
	50 * MB,
	500 * MB,

	5 * GB,
	50 * GB,

	5 * TB,
}

var log *logrus.Logger

type dirInfo struct {
	checksum string
	dbPath   string
}

func newDirInfo(srcDir, dbDir string) *dirInfo {
	b := md5.Sum([]byte(srcDir))
	checksum := hex.EncodeToString(b[:])
	return &dirInfo{
		checksum: checksum,
		dbPath:   filepath.Join(dbDir, fmt.Sprintf(FilesDbName, checksum)),
	}
}
