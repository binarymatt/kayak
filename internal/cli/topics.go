package cli

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/urfave/cli/v2"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func CreateTopic(ctx *cli.Context) error {
	client := buildClient(ctx.String("host"))
	_, err := client.CreateTopic(ctx.Context, connect.NewRequest(&kayakv1.CreateTopicRequest{
		Topic: &kayakv1.Topic{
			Name:       ctx.String("name"),
			Partitions: 1,
		},
	}))
	return err
}
func ListTopics(cCtx *cli.Context) error {
	client := buildClient(cCtx.String("host"))
	resp, err := client.ListTopics(cCtx.Context, connect.NewRequest(&kayakv1.ListTopicsRequest{}))
	if err != nil {
		return err
	}
	for _, topic := range resp.Msg.Topics {
		fmt.Println(topic.GetName())
	}
	return nil
}
func DeleteTopic(ctx *cli.Context) error {
	client := buildClient(ctx.String("host"))
	_, err := client.DeleteTopic(ctx.Context, connect.NewRequest(&kayakv1.DeleteTopicRequest{
		Topic: &kayakv1.Topic{
			Name:       ctx.String("name"),
			Archived:   ctx.Bool("archive"),
			Partitions: 1,
		},
	}))
	return err
}
