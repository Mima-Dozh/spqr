package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"syscall"

	coordApp "github.com/pg-sharding/spqr/coordinator/app"
	"github.com/pg-sharding/spqr/coordinator/provider"
	"github.com/pg-sharding/spqr/pkg"
	"github.com/pg-sharding/spqr/pkg/config"
	"github.com/pg-sharding/spqr/pkg/datatransfers"
	"github.com/pg-sharding/spqr/pkg/spqrlog"
	"github.com/pg-sharding/spqr/qdb"
	"github.com/pg-sharding/spqr/router/app"
	"github.com/pg-sharding/spqr/router/instance"
	"github.com/pkg/errors"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

var (
	rcfgPath     string
	cpuProfile   bool
	memProfile   bool
	profileFile  string
	daemonize    bool
	console      bool
	logLevel     string
	gomaxprocs   int
	pgprotoDebug bool

	ccfgPath string
	qdbImpl  string

	persist bool

	routerPort   int
	routerROPort int
	adminPort    int
	grpcPort     int

	default_route_behaviour string

	rootCmd = &cobra.Command{
		Use:   "spqr-router run --config `path-to-config-folder`",
		Short: "spqr-router",
		Long:  "spqr-router",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Version:       pkg.SpqrVersionRevision,
		SilenceUsage:  false,
		SilenceErrors: false,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&rcfgPath, "config", "c", "/etc/spqr/router.yaml", "path to router config file")
	rootCmd.PersistentFlags().StringVarP(&profileFile, "profile-file", "p", "/etc/spqr/router.prof", "path to profile file")
	rootCmd.PersistentFlags().BoolVarP(&daemonize, "daemonize", "d", false, "daemonize router binary or not")
	rootCmd.PersistentFlags().BoolVarP(&console, "console", "", false, "console (not daemonize) router binary or not")
	rootCmd.PersistentFlags().BoolVar(&cpuProfile, "cpu-profile", false, "profile cpu or not")
	rootCmd.PersistentFlags().BoolVar(&memProfile, "mem-profile", false, "profile mem or not")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "log level")
	rootCmd.PersistentFlags().IntVarP(&gomaxprocs, "gomaxprocs", "", 0, "GOMAXPROCS value")

	rootCmd.PersistentFlags().StringVarP(&ccfgPath, "coordinator-config", "", "/etc/spqr/coordinator.yaml", "path to coordinator config file")
	rootCmd.PersistentFlags().StringVarP(&qdbImpl, "qdb-impl", "", "etcd", "which implementation of QDB to use.")
	rootCmd.PersistentFlags().BoolVarP(&persist, "persist", "", false, "tells router to persist its configuration in non-clustered setup")

	rootCmd.PersistentFlags().IntVarP(&routerPort, "router-port", "", 0, "router PostgreSQL port")
	rootCmd.PersistentFlags().IntVarP(&routerROPort, "router-ro-port", "", 0, "router read-only PostgreSQL port")
	rootCmd.PersistentFlags().IntVarP(&adminPort, "admin-port", "", 0, "router Metadata PostgreSQL interface admin port")
	rootCmd.PersistentFlags().IntVarP(&grpcPort, "grpc-port", "", 0, "router Metadata GRPC interface admin port")

	rootCmd.PersistentFlags().StringVarP(&default_route_behaviour, "default-route-behaviour", "", "", "router block or scatters-out failed to route statements")

	rootCmd.PersistentFlags().BoolVarP(&pgprotoDebug, "proto-debug", "", false, "reply router notice, warning, etc")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run router",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.LoadRouterCfg(rcfgPath)
		if err != nil {
			return err
		}
		rcfg := config.RouterConfig()

		spqrlog.ReloadLogger(rcfg.LogFileName)

		// Logger
		rlogLevel := rcfg.LogLevel
		if logLevel != "" {
			rlogLevel = logLevel
		}

		if err := spqrlog.UpdateZeroLogLevel(rlogLevel); err != nil {
			return err
		}

		if console && daemonize {
			return fmt.Errorf("simultaneous use of `console` and `daemonize`. Abort")
		}

		if persist && qdbImpl == "etcd" {
			return fmt.Errorf("Cannot persist metadata setup locally in clustered mode. Abort")
		}

		if !persist && rcfg.MemqdbPersistent {
			persist = true
		}

		if !console && (rcfg.Daemonize || daemonize) {
			cntxt := &daemon.Context{
				PidFileName: rcfg.PidFileName,
				PidFilePerm: 0644,
				WorkDir:     "./",
				Umask:       027,
				Args:        args,
			}

			d, err := cntxt.Reborn()
			if err != nil {
				log.Fatal("Unable to run: ", err)
			}
			if d != nil {
				return nil
			}

			defer func() {
				if err := cntxt.Release(); err != nil {
					spqrlog.Zero.Error().Msg("")
					spqrlog.Zero.Error().Err(err).Msg("")
				}
			}()

			spqrlog.Zero.Debug().Msg("daemon started")
		}

		if rcfg.UseCoordinatorInit && rcfg.UseInitSQL {
			return fmt.Errorf("cannot use initSQL and coordinator-based init simultaneously")
		}

		ctx, cancelCtx := context.WithCancel(context.Background())
		defer cancelCtx()

		var pprofCpuFile *os.File
		var pprofMemFile *os.File

		if cpuProfile {
			spqrlog.Zero.Info().Msg("starting cpu profile")
			pprofCpuFile, err = os.Create(path.Join(path.Dir(profileFile), "cpu"+path.Base(profileFile)))

			if err != nil {
				spqrlog.Zero.Info().
					Err(err).
					Msg("got an error while starting cpu profile")
				return err
			}

			if err := pprof.StartCPUProfile(pprofCpuFile); err != nil {
				spqrlog.Zero.Info().
					Err(err).
					Msg("got an error while starting cpu profile")
				return err
			}
		}

		if memProfile {
			spqrlog.Zero.Info().Msg("starting mem profile")
			pprofMemFile, err = os.Create(path.Join(path.Dir(profileFile), "mem"+path.Base(profileFile)))
			if err != nil {
				spqrlog.Zero.Info().
					Err(err).
					Msg("got an error while starting mem profile")
				return err
			}
		}

		if gomaxprocs > 0 {
			runtime.GOMAXPROCS(gomaxprocs)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

		/* will change on reload */
		rcfg.PgprotoDebug = rcfg.PgprotoDebug || pgprotoDebug
		rcfg.ShowNoticeMessages = rcfg.ShowNoticeMessages || pgprotoDebug

		if routerPort != 0 {
			rcfg.RouterPort = strconv.FormatInt(int64(routerPort), 10)
		}

		if routerROPort != 0 {
			rcfg.RouterROPort = strconv.FormatInt(int64(routerROPort), 10)
		}

		if adminPort != 0 {
			rcfg.AdminConsolePort = strconv.FormatInt(int64(adminPort), 10)
		}

		if grpcPort != 0 {
			rcfg.GrpcApiPort = strconv.FormatInt(int64(grpcPort), 10)
		}

		if default_route_behaviour != "" {
			rcfg.Qr.DefaultRouteBehaviour = default_route_behaviour
		}

		router, err := instance.NewRouter(ctx, rcfg, os.Getenv("NOTIFY_SOCKET"), persist)
		if err != nil {
			return errors.Wrap(err, "router failed to start")
		}

		app := app.NewApp(router)

		if rcfgPath != "" {
			if err := datatransfers.LoadConfig(rcfgPath); err != nil {
				return err
			}
		}
		if config.RouterConfig().WithCoordinator {
			go func() {
				if err := func() error {
					if err := config.LoadCoordinatorCfg(ccfgPath); err != nil {
						return err
					}

					db, err := qdb.NewXQDB(qdbImpl)
					if err != nil {
						return err
					}

					frTLS, err := config.CoordinatorConfig().FrontendTLS.Init(config.CoordinatorConfig().Host)
					if err != nil {
						return fmt.Errorf("init frontend TLS: %w", err)
					}

					coordinator, err := provider.NewCoordinator(frTLS, db)
					if err != nil {
						return err
					}

					app := coordApp.NewApp(coordinator)
					return app.Run(false)
				}(); err != nil {
					spqrlog.Zero.Error().Err(err).Msg("")
				}
			}()
		}
		go func() {
			defer cancelCtx()
			for {
				s := <-sigs
				spqrlog.Zero.Info().Str("signal", s.String()).Msg("received signal")

				switch s {
				case syscall.SIGUSR1:
					spqrlog.ReloadLogger(rcfg.LogFileName)
				case syscall.SIGUSR2:
					if cpuProfile {
						// write profile
						pprof.StopCPUProfile()
						spqrlog.Zero.Info().Str("fname", pprofCpuFile.Name()).Msg("writing cpu prof")

						if err := pprofCpuFile.Close(); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
					}
					if memProfile {
						// write profile
						spqrlog.Zero.Info().Str("fname", pprofMemFile.Name()).Msg("writing mem prof")

						if err := pprof.WriteHeapProfile(pprofMemFile); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
						if err := pprofMemFile.Close(); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
					}
					return
				case syscall.SIGHUP:
					// reread config file
					err := router.RuleRouter.Reload(rcfgPath)
					if err != nil {
						spqrlog.Zero.Error().Err(err).Msg("")
					}
					spqrlog.ReloadLogger(rcfg.LogFileName)
				case syscall.SIGINT, syscall.SIGTERM:
					if cpuProfile {
						// write profile
						pprof.StopCPUProfile()

						spqrlog.Zero.Info().Msg("writing cpu prof")
						if err := pprofCpuFile.Close(); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
					}

					if memProfile {
						// write profile
						spqrlog.Zero.Info().Msg("writing mem prof")

						if err := pprof.WriteHeapProfile(pprofMemFile); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
						if err := pprofMemFile.Close(); err != nil {
							spqrlog.Zero.Error().Err(err).Msg("")
						}
					}
					return
				default:
					return
				}
			}
		}()

		/* initialize metadata */
		if rcfg.UseInitSQL {
			i := instance.NewInitSQLMetadataBootstraper(rcfg.InitSQL)
			if err := i.InitializeMetadata(ctx, router); err != nil {
				return err
			}
		} else if rcfg.UseCoordinatorInit {
			/* load config if not yet */
			if err := config.LoadCoordinatorCfg(ccfgPath); err != nil {
				return err
			}
			e := instance.NewEtcdMetadataBootstraper(config.CoordinatorConfig().QdbAddr)
			if err := e.InitializeMetadata(ctx, router); err != nil {
				return err
			}
		} else {
			/* TODO: maybe error-out? */
			router.Initialize()
		}

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := app.ServeRouter(ctx)
			if err != nil {
				spqrlog.Zero.Error().Err(err).Msg("")
			}
			wg.Done()
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := app.ServeGrpcApi(ctx)
			if err != nil {
				spqrlog.Zero.Error().Err(err).Msg("")
			}
			wg.Done()
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := app.ServeAdminConsole(ctx)
			if err != nil {
				spqrlog.Zero.Error().Err(err).Msg("")
			}
			wg.Done()
		}(wg)

		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			err := app.ServceUnixSocket(ctx)
			if err != nil {
				spqrlog.Zero.Error().Err(err).Msg("")
			}
			wg.Done()
		}(wg)

		wg.Wait()

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		spqrlog.Zero.Fatal().Err(err).Msg("")
	}
}
