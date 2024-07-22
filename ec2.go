package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func ec2List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "List and search EC2 instances",
	}

	var line, profile, out, preview string
	flags := cmd.Flags()
	flags.StringVarP(&profile, "profile", "p", "", "AWS profile to use")
	flags.StringVarP(&out, "out", "o", defaultOutTpl, "Output template (default is as per profile)")
	flags.StringVarP(&line, "line", "l", defaultLineTpl, "Line template (default is as per profile)")
	flags.StringVarP(&preview, "preview", "v", defaultPreviewTpl, "Preview template (default is as per profile)")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadDefaultConfig(cmd.Context(),
			config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			log.Printf("failed to load configs: %v\n", err)
			os.Exit(1)
		}

		cl := ec2.NewFromConfig(cfg)

		ins, err := listInstances(cmd.Context(), cl, args)
		if err != nil {
			log.Printf("failed to list instances: %v\n", err)
			os.Exit(1)
		}

		funcs := template.FuncMap{
			"tag": func(ins types.Instance, key string) string {
				for _, tag := range ins.Tags {
					if *tag.Key == key {
						return *tag.Value
					}
				}
				return ""
			},
			"tag_join": func(inst types.Instance, keys ...string) string {
				var tags []string
				for _, key := range keys {
					tags = append(tags, strings.TrimSpace(findTag(inst.Tags, key)))
				}
				return strings.Join(tags, " : ")
			},
			"name": func(ins types.Instance) string {
				return findName(ins.Tags)
			},
		}

		lineTpl := template.New("line").Funcs(funcs)
		outTpl := template.New("output").Funcs(funcs)
		previewTpl := template.New("preview").Funcs(funcs)

		lineTpl, err = lineTpl.Parse(line)
		if err != nil {
			log.Printf("failed to parse line template: %v\n", err)
			os.Exit(1)
		}

		outTpl, err = outTpl.Parse(out)
		if err != nil {
			log.Printf("failed to parse output template: %v\n", err)
			os.Exit(1)
		}

		previewTpl, err = previewTpl.Parse(preview)
		if err != nil {
			log.Printf("failed to parse preview template: %v\n", err)
			os.Exit(1)
		}

		idx, err := fuzzyfinder.Find(ins, func(i int) string {
			var buf strings.Builder
			_ = lineTpl.Execute(&buf, ins[i])
			return buf.String()
		}, fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			var buf strings.Builder
			_ = previewTpl.Execute(&buf, ins[i])
			return buf.String()
		}))
		if err != nil {
			log.Printf("failed to find instance: %v\n", err)
			os.Exit(1)
		}

		var out strings.Builder
		if err := outTpl.Execute(&out, ins[idx]); err != nil {
			log.Printf("failed to execute output template: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(out.String())
	}

	return cmd
}

func listInstances(ctx context.Context, cl *ec2.Client, filterSpecs []string) ([]types.Instance, error) {
	var filters []types.Filter

	for _, filter := range filterSpecs {
		split := strings.SplitN(filter, "=", 2)
		if len(split) < 2 {
			return nil, fmt.Errorf("invalid filter: %s", filter)
		}

		filters = append(filters, types.Filter{
			Name:   aws.String(split[0]),
			Values: []string{split[1]},
		})
	}

	res, err := cl.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	var instances []types.Instance
	for _, res := range res.Reservations {
		instances = append(instances, res.Instances...)
	}

	return instances, err
}

func findName(tags []types.Tag) string {
	return findTag(tags, "Name")
}

func findTag(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return "?"
}
