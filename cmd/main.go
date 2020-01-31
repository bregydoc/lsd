package main

import (
	"fmt"
	"net"
	"os"

	"github.com/bregydoc/lsd"
	"github.com/bregydoc/lsd/api"
	proto "github.com/bregydoc/lsd/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := os.Setenv("GIN_MODE", "release"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	viper.SetDefault("DatabaseFilepath", "lsd.db")
	viper.SetDefault("Secure", false)
	viper.SetDefault("GRPCPort", 10000)
	viper.SetDefault("APIPort", 8080)
	viper.SetDefault("WSPort", 3300)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/lsd")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	lsdEngine, err := lsd.NewLSD(viper.GetString("DatabaseFilepath"), viper.GetBool("Secure"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("GRPCPort")))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterLSDServer(grpcServer, lsdEngine)

	// TODO (bregydoc): determine whether to use TLS

	done := make(chan error, 1)

	go func() {
		log.Infof("serving grpc server on :%d", viper.GetInt("GRPCPort"))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
			done <- err
		}
	}()

	go func() {
		log.Infof("serving http api server on :%d", viper.GetInt("APIPort"))
		if err := api.New(lsdEngine, nil).Run(viper.GetInt("APIPort")); err != nil {
			log.Fatal(err)
			done <- err
		}
	}()

	go func() {
		log.Infof("serving ws server on :%d", viper.GetInt("WSPort"))
		if err := lsdEngine.RunWSService(viper.GetInt("WSPort")); err != nil {
			log.Fatal(err)
			done <- err
		}
	}()

	if err := <-done; err != nil {
		log.Fatal(err)
	}


}
