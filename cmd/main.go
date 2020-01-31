package main

import (
	"fmt"
	"net"

	"github.com/bregydoc/lsd"
	"github.com/bregydoc/lsd/api"
	proto "github.com/bregydoc/lsd/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	viper.SetDefault("DatabaseFilepath", "lsd.db")
	viper.SetDefault("Secure", false)
	viper.SetDefault("GRPCPort", 10000)
	viper.SetDefault("APIPort", 8080)

	viper.SetConfigName("lsd.config.yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/lsd")

	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
		panic(err)
	}

	lsdEngine, err := lsd.NewLSD(viper.GetString("DatabaseFilepath"), viper.GetBool("Secure"))
	if err != nil {
		log.Error(err)
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("GRPCPort")))
	if err != nil {
		log.Error("failed to listen: %v", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterLSDServer(grpcServer, lsdEngine)

	// TODO (bregydoc): determine whether to use TLS

	done := make(chan error, 1)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Error(err)
			done <- err
		}
	}()

	go func() {
		if err := api.New(lsdEngine, nil).Run(viper.GetInt("APIPort")); err != nil {
			log.Error(err)
			done <- err
		}
	}()

	if err := <-done; err != nil {
		log.Error(err)
		panic(err)
	}
}
