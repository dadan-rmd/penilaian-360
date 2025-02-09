package cmd

import (
	"errors"
	"flag"
	"os"

	"penilaian-360/internal/app/commons"
	"penilaian-360/internal/app/commons/applicationConstants"
	"penilaian-360/internal/app/migration"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate Up DB KBR PRIME API",
	Long:  `Please you know what are you doing by using this command`,
	Run: func(cmd *cobra.Command, args []string) {
		envName, _ := cmd.Flags().GetString("env")
		loadEnv(envName)

		opt, err := initCommonOptions()
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		runMigration(&opt, applicationConstants.MigrateUp)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "migratedown",
	Short: "Migrate Up DB Loan API",
	Long:  `Please you know what are you doing by using this command`,
	Run: func(cmd *cobra.Command, args []string) {
		envName, _ := cmd.Flags().GetString("env")
		loadEnv(envName)

		opt, err := initCommonOptions()
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		runMigration(&opt, applicationConstants.MigrateDown)
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(migrateDownCmd)

	migrateUpCmd.PersistentFlags().StringP("env", "e", "prod", "environment type")
	migrateDownCmd.PersistentFlags().StringP("env", "e", "prod", "environment type")
}

func runMigration(opt *commons.Options, direction int) {
	pathMigration := os.Getenv("APP_MIGRATION_PATH")
	migrationDir := flag.String("migration-dir", pathMigration, "migration directory")
	log.Info().Msg("path migration : " + pathMigration)

	migrationConf, errMigrationConf := migration.NewMigrationConfig(*migrationDir,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE_NAME"),
		"mysql")
	if errMigrationConf != nil {
		log.Error().Msg(errMigrationConf.Error())
		return
	}

	var errMigration error
	switch direction {
	case applicationConstants.MigrateUp:
		errMigration = migration.MigrateUp(migrationConf)
	case applicationConstants.MigrateDown:
		errMigration = migration.MigrateDown(migrationConf)
	default:
		errMigration = errors.New("Unknown migration direction")
	}
	if errMigration != nil {
		if errMigration.Error() != "no change" {
			log.Error().Msg(errMigration.Error())
			return
		}
		log.Info().Msg("Migration success : no change table . . .")
	}
}
