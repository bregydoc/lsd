package main

import (
	"context"
	"fmt"

	proto "github.com/bregydoc/lsd/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	viper.SetDefault("GRPCPort", 10000)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", viper.GetInt("GRPCPort")), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	lsd := proto.NewLSDClient(conn)

	res, err := lsd.GenerateNewKeyPair(context.Background(), &proto.NewKeyPairPayload{
		UserID:               "bregydoc",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info("userID: ", res.UserID)
	log.Info("publicKey: ", res.PublicKey)
}

