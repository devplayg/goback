package goback

import (
	"context"
	"encoding/json"
	"github.com/devplayg/golibs/compress"
	"github.com/devplayg/himma/v2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const LogoImg = "H4sIAAAAAAAA/wBCBL37iVBORw0KGgoAAAANSUhEUgAAAB4AAAAeCAYAAAA7MK6iAAAACXBIWXMAAC4jAAAuIwF4pT92AAAD9ElEQVR42sWX709bVRjHT6uwV/oXGN8bX/hOZyIZY5vNwgphGjRZoiTGaYz6Yrhh293e/liCLiRtX1mD24hJMcB6SwddnU5cLFMGZJslkQED2vurP6hi9iPbKJzH57kUZdMtBFrX5JOn9+acz/eep+fe3DKXy8UeB/86IVIVRaOKLvcTW+Mf1yOD7xuI1ev4lHmO2phHsGPdIMZYmzF3vUt85IrXBjgFo37QOfBWS+/lEy29V75u6dsgvQYnaK5znYvc/xlstFakAU6T4PZuOyBNRq0X7kLj0O1NcAusF+7BgfBvEcHtqRbRSW7xYSt24u/icbSx976KvVv/411oiufuNJ0rFDcFziXH+52D75CT3A9ZMbZZFPHKRNYsTf2879zickM8X2z4tsBLwAZZHYtzydEsXRsmp+F2ifcHUwucostEm6O1I/h8fbywZI3lVvadzXPr2TxsAm7MRQe5Wju+eI7c2G6T+OCKqRVebMnbJ3/wWuKLUD+g36sfzMKWQIcl/ge0nDzv9j7Q7tUV045zOpngPlZlDV+fsgxk+d4z+vLegQxHYJMYDnI1hGcmccM+SRliaXevtdnsxXvvI1/Xjt2DebD0K8uWqMYR2CKcXOT82HeqhjIoSzRWjFeA95uZbvo3uhKdO8/k+B5JXtoTUaAsoIucb5766UvKoCwXbrK/Hxg27+dPv9o3m90lyXxXWF7ZjRWBLWK46rCiO2PzfvaU8UDBe5oJLk81teBg4JvmmogOdX1zxbrTKSgr6CT3wUD365SFmVW006rczqOssWsk+spphdf2zBZre+c4AmXCcJK7seuXCGUJLnc1oys43O5/trbn+q0dPbO8BJQZw4sZN4+0+57xCA7GjtsPmfb7I4df6E7By6Gp4vbQDFQCclPGa37p0HFHq4l9YnMwMXZ15PvpLE9Ma8XEtM4TMzqUFXQOo5syhNjVS0fa2hiLxuIf5jP6nRuFHCzms/zPhSxUAnLfWMjyhVx2qScs7WcZXdc1TQdF00DFWjk0kFV1hbImksmLTMGPqqocgXW17OiqzjGK/Hx8fHyYybKsIhyBSqHICkzJ1yCdThs5Y2NjwwwPVDqBQLmZT81DJp2FofnzsHP6RfhuLs5zch5GRkeMYG0tmFZezuBUOgVqWoXk/K/QMdnOr8xd5pqs8dGx0YssmUymSoHFSoDhRSWtFAvp34tKSlmSFRkSicQQCwaDTRh+U6NdXdoI9L0cqNqqi5xpJWUcT0xM3A4EAjXMbrczv9+/XZKkYDQaDfX394ewdpeTkjMUDoeDPp/vJcqktwGzIAjM4XD8L1CW2+02rb2KmPGA3ocqSinDTJnscf1p+wtMni0TTGHjQwAAAABJRU5ErkJgggEAAP//5/mgNEIEAAA="

var (
	WebAssetMap himma.AssetMap
)

type Controller struct {
	router    *mux.Router
	summaryDb *os.File
	addr      string
	dbDir     string
	version   string
	app       *himma.Config
	summaries []*Summary
	server    *Server
}

func NewController(server *Server, app *himma.Config) *Controller {
	return &Controller{
		server:    server,
		addr:      server.appConfig.Address,
		dbDir:     server.dbDir,
		summaries: make([]*Summary, 0),
		version:   app.Version,
		app:       app,
	}
}

func (c *Controller) init() error {
	if err := c.initRouter(); err != nil {
		return err
	}
	uiAssetMap, err := himma.GetAssetMap(
		himma.Bootstrap4,
		himma.BootstrapDatepicker_1_9_0,
		himma.BootstrapSelect_1_13_9,
		himma.BootstrapTable_1_16_0,
		himma.Holder_2_9,
		himma.JqueryMask_1_14_16,
		himma.JqueryValidation_1_19_1,
		himma.JsCookie_2_2_1,
		himma.Moment_2_24_0,
		himma.ProgressbarJs_1_0_1,
		himma.WaitMe_31_10_17,
	)
	if err != nil {
		return err
	}
	// f, _ := ioutil.ReadFile("logo-goback.png")
	// uiAssetMap["/assets/img/logo.png"] = LogoImg
	// d, _ := compress.Compress(f, compress.GZIP)
	// fmt.Println(base64.StdEncoding.EncodeToString(d))
	// himma.Add(uiAssetMap,"/assets/img/logo.png", LogoImg)
	uiAssetMap.AddZippedAndBase64Encoded("/assets/img/logo.png", LogoImg)
	uiAssetMap.AddRaw("/assets/js/custom.js", customScript())
	WebAssetMap = uiAssetMap

	//if err := c.loadSummaryDb(); err != nil {
	//	return err
	//}
	return nil
}

func (c *Controller) initRouter() error {
	c.router = mux.NewRouter()
	c.router.HandleFunc("/", c.DisplayDefault)

	// Assets
	c.router.HandleFunc("/assets/{u1}/{u2}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}/{u2}/{u3}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})
	c.router.HandleFunc("/assets/{u1}/{u2}/{u3}/{u4}", func(w http.ResponseWriter, r *http.Request) {
		GetAsset(w, r)
	})

	// Backup
	c.router.HandleFunc("/backup/", c.DisplayBackup)
	c.router.HandleFunc("/summaries", c.GetSummaries)
	c.router.HandleFunc("/stats", c.GetStats)
	c.router.HandleFunc("/summaries/{id:[0-9]+}/changes", c.GetChangesLog)

	// Statistics
	c.router.HandleFunc("/stats/", c.DisplayStats)

	// Settings
	c.router.HandleFunc("/settings/", c.DisplaySettings)
	c.router.HandleFunc("/settings/job/id/{id:[0-9]+}", c.UpdateJob).Methods(http.MethodPatch)
	c.router.HandleFunc("/settings/storage/id/{id:[0-9]+}", c.UpdateStorage).Methods(http.MethodPatch)
	c.router.HandleFunc("/backup/{id:[0-9]+}/run", c.RunBackupJob).Methods(http.MethodGet)

	http.Handle("/", c.router)

	return nil
}

func (c *Controller) Start() error {
	if err := c.init(); err != nil {
		return err
	}
	defer c.Stop()

	srv := &http.Server{
		Handler:      c.router,
		Addr:         c.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	ch := make(chan struct{})
	go func() {
		<-c.server.Ctx.Done()
		c.server.Log.Info("2) got stop signal from engine")
		if err := srv.Shutdown(context.Background()); err != nil {
			c.server.Log.Error(err)
		}
		close(ch)
	}()

	c.server.Log.Debug("1) HTTP server has been started")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		c.server.Log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-ch
	c.server.Log.Debug("3) HTTP server has been stopped")
	return nil
}

func (c *Controller) Stop() error {
	return nil
}

func GetAsset(w http.ResponseWriter, r *http.Request) {
	if content, hasAsset := WebAssetMap[r.RequestURI]; hasAsset {
		w.Header().Set("Content-Type", DetectContentType(filepath.Ext(r.RequestURI)))
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(content)), 10))
		w.Header().Set("Content-Encoding", compress.GZIP)
		w.Write(content)
	}
}

func DetectContentType(ext string) string {
	ctype := mime.TypeByExtension(filepath.Ext(ext))
	if ctype == "" {
		return "application/octet-stream"
	}
	return ctype
}

func Response(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	if statusCode != http.StatusOK {
		log.WithFields(logrus.Fields{
			"ip":     r.RemoteAddr,
			"uri":    r.RequestURI,
			"method": r.Method,
			"length": r.ContentLength,
		}).Error(err)
	}
	w.Header().Add("Content-Type", "application/json")
	b, _ := json.Marshal(map[string]interface{}{
		"error": err.Error(),
	})
	w.WriteHeader(statusCode)
	w.Write(b)
}
